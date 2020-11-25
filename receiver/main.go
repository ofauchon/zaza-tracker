package main

import (
	"device/stm32"
	"encoding/hex"
	"machine"
	"runtime/interrupt"
	"runtime/volatile"
	"strconv"
	"strings"
	"time"

	"github.com/ofauchon/zaza-tracker/drivers"
	"tinygo.org/x/drivers/gps"
	"tinygo.org/x/drivers/lora/sx127x"
)

// Lorawan tests (The thing network)
/*
    device = "track01"
	myAppKey  = "2C44FCF86C7B767B8FD3124FCE7A3216"
	myDevEUI  = "D0000000000AA001"
	myJoinEUI/AppEUI = "A000000000000102"
*/

const (
	led    = machine.LED_RED
	button = machine.PB14
)

type status struct {
	fix *gps.Fix
}

/*
var loraConfig = sx127x.Config{
	Frequency:       868000000,
	SpreadingFactor: 12,
	Bandwidth:       125000,
	CodingRate:      8,
	TxPower:         5,
	PaBoost:         true,
}
*/

// Lorawan configuration
var loraConfig = sx127x.Config{
	Frequency:            868300000,
	SpreadingFactor:      7,
	Bandwidth:            125000,
	CodingRate:           4,
	TxPower:              20,
	PaBoost:              true,
	SyncWord:             0x34, // Lorawan sync
	ImplicitHeaderModeOn: false,
	CrcOn:                true,
}

var (
	st                   status
	lightLed             volatile.Register8
	uartConsole, uartGps *machine.UART
	send_data            = string("")
	send_delay           = int(0)
	packet               [255]byte
	loraRadio            sx127x.Device
	keypressed           bool
)

var e2p drivers.Eeprom // Eeprom driver

var rxPktChan chan []byte // Channel for packet RX

var config *drivers.ATConfig

// loraJoin connects Lorawan network
func loraJoin() {

	l := &drivers.LightLW{}
	l.Otaa.AppEUI = config.GetCurrentValue("APPEUI").([8]uint8)
	l.Otaa.DevEUI = config.GetCurrentValue("DEVEUI").([8]uint8)
	l.Otaa.AppKey = config.GetCurrentValue("APPKEY").([16]uint8)
	l.Otaa.DevNonce = uint16(0xCC85)

	j := l.GenerateJoinRequest()

	println(hex.EncodeToString(j))

}

// processCmd() parses commands and execute actions
func processCmd(cmd string) {

	println("Processing [" + cmd + "]")
	if len(cmd) > 6 && strings.HasPrefix(cmd, "AT+") && strings.Contains(cmd, "=") {
		// Processing AT+XXX=YYY Set configuration
		r := strings.Split(cmd, "+")
		if len(r) == 2 && len(r[1]) > 3 {
			s := strings.Split(r[1], "=")
			if len(s) == 2 && len(s[0]) > 0 && len(s[1]) > 0 {
				cf := config.GetConfigEntry(s[0])
				if cf != nil {
					switch cf.FactoryValue.(type) {
					case []uint8:
						bb, e := hex.DecodeString(s[1])
						if e == nil && len(bb) == len(cf.FactoryValue.([]uint8)) {
							cf.CurrentValue = bb
							println("Set " + s[0] + " OK")
						} else {
							println("Wrong parameter, expected [" + string(len(bb)) + "]uint8")
						}
					case uint32:
						i, err := strconv.ParseInt(s[1], 10, 64)
						if err == nil && i > 0x0 && i < 0xFFFFFFFF {
							cf.CurrentValue = uint32(i)
						} else {
							println("Wrong parameter, expected uint32")
						}
					case uint8:
						i, err := strconv.Atoi(s[1])
						if err == nil && i > 0x0 && i < 0xFF {
							cf.CurrentValue = uint8(i)
						} else {
							println("Wrong parameter, expected uint8")
						}
					}
				} else {
					println("Unknown " + s[0] + " config parameter.")
				}
			}
		}
	}

	ss := strings.Split(cmd, " ")
	switch ss[0] {

	case "AT+CFG":
		println("Configuration Dump:")
		println(config.DumpConfig())
	case "AT+FDR":
		println("Resetting configuration to factory defaults")
		config.FactoryReset()
	case "ATZ":
		println("System reboot Now")
		stm32.SCB.AIRCR.Set(uint32(0x05fa0004))
	case "AT+VER":
		println("1.3 EU868 ZAZA")
	case "AT+HWVER":
		println("L70-RL")

	case "help":
		//TODO
	case "conf":
		if len(ss) > 1 {
			switch ss[1] {
			case "read":
				//config.Read(e2p)
			case "write":
				//config.Write(e2p)
			case "default":
				//config.Default()
			case "dump":
				//config.Dump()
			}
		}

	case "lorawan":
		if len(ss) > 1 {
			switch ss[1] {
			case "join":
				loraJoin()
			}
		}

	case "eepw":
		if len(ss) == 3 {
			t1, err := strconv.ParseUint(ss[1], 16, 64)
			p := uint32(t1)
			if err == nil {
				t2, err := strconv.ParseUint(ss[2], 16, 64)
				b := uint8(t2)
				if err == nil {
					e2p.Unlock()
					println("Write eeprom : offset:", p, " byte:", b)
					e2p.WriteU8(b, p)
				}
			} else {
				println("Wrong byte value")
			}
		} else {
			println("Wrong pos value")

		}

	case "eepr":
		if len(ss) == 2 {
			p, err := strconv.ParseUint(ss[1], 16, 64)
			if err == nil {
				v := e2p.ReadU8(uint32(p))
				println("Read eeprom pos:", p, " value:", v)

			} else {
				println("Wrong pos value")
			}
		}

	// Send Lora packets
	case "loratx":
		if len(ss) == 2 {
			tmp, err := hex.DecodeString(ss[1])
			if err != nil {
				println("Invalid packet payload, can't send")
			} else {
				loraRadio.TxLora(tmp)
				println("LoraTX ", len(tmp), "bytes sent")

			}
		}

	// Listen for Lora packets
	case "lorarx":
		keypressed = false
		go func() {
			println("lorarx: Start RXContinuous")
			loraRadio.SetOpMode(sx127x.OPMODE_RX)

			rxchan := loraRadio.GetRxPktChannel()

			for !keypressed {
				//println("RX packet: Waiting for new packet")
				packet := <-rxchan
				println("RX packet: '", string(packet), "'")

				/*
					packetSize := loraRadio.ParsePacket(0)
					if packetSize > 0 {
						//println("Got packet, RSSI=", loraRadio.LastPacketRSSI())
						size := loraRadio.ReadPacket(packet[:])
						println("RX: ", string(packet[:size]), " packetsize", packetSize)
					}
				*/
			}
			println("lorarx: Stop RXContinuous")
		}()

	case "get":
		if len(ss) == 2 {
			switch ss[1] {
			case "sx1276config":
				println("Frequency:", loraRadio.GetFrequency())
				println("SpreadingFactor:", loraRadio.GetSpreadingFactor())
				println("Bandwidth:", loraRadio.GetBandwidth())
			case "temp":
				//temp, _ := d.ReadTemperature(0)
				println("Temperature:")
			case "regs":
				println(" Regs:")
				loraRadio.PrintRegisters(true)
			default:
				println("Invalid use of 'get'")
			}
		}

	case "set":
		if len(ss) == 3 {
			switch ss[1] {
			case "freq":
				val, _ := strconv.ParseUint(ss[2], 10, 32)
				//	d.SetFrequency(uint32(val))
				println("Freq set to ", val)
			case "power":
				val, _ := strconv.ParseUint(ss[2], 10, 32)
				//	d.SetTxPower(uint8(val))
				println("TxPower set to ", val)
			}
		} else {
			println("invalid use of 'set' command")
		}

	case "mode":
		if len(ss) == 2 {
			switch ss[1] {
			case "standby":
				//d.SetMode(rfm69.RFM69_MODE_STANDBY)
				//d.WaitForMode()
				println("Mode changed !")
			case "sleep":
				//d.SetMode(rfm69.RFM69_MODE_SLEEP)
				//d.WaitForMode()
				println("Mode changed !")
			case "tx":
				//d.SetMode(rfm69.RFM69_MODE_TX)
				//d.WaitForMode()
				println("Mode changed !")
			case "rx":
				//d.SetMode(rfm69.RFM69_MODE_RX)
				//d.WaitForMode()
				println("Mode changed !")
			default:
				println("Invalid use of 'mode'")
			}
		}
	default:
		println("Command Error")
	}

}

// serial() function is a gorouting for handling USART rx data
func serial(serial *machine.UART) string {
	input := make([]byte, 100) // serial port buffer

	i := 0

	for {

		if i == 100 {
			println("Serial Buffer overrun")
			i = 0
		}

		if serial.Buffered() > 0 {
			keypressed = true
			data, _ := serial.ReadByte() // read a character
			switch data {
			case 13: // pressed return key
				uartConsole.Write([]byte("\r\n"))
				cmd := string(input[:i])
				processCmd(cmd)
				i = 0
			default: // pressed any other key
				uartConsole.WriteByte(data)
				input[i] = data
				i++
			}
		}

		time.Sleep(10 * time.Millisecond)
	}

}

// Interrupt handler from RFM95_DIO0 (RxDone Event) on PB13
func gpios_int(inter interrupt.Interrupt) {
	irqStatus := stm32.EXTI.PR.Get()
	stm32.EXTI.PR.Set(irqStatus)

	if (irqStatus & 0x2000) > 0 { // PC13 : DIO

		/*
		   println("pgios_int: packetSize:", packetSize)
		   		if packetSize > 0 {
		   			size := loraRadio.ReadPacket(packet[:])
		   			rxPktChan <- packet[:size]
		   		}
		*/
		loraRadio.DioIntHandler()

	}
	if (irqStatus & 0x4000) > 0 { // PB14 : Button
		println("Button: ", machine.BUTTON.Get())
	}

}

// Get Random 32bit
func getRng() uint32 {
	// Enable PRNG
	stm32.RNG.CR.SetBits(stm32.RNG_CR_RNGEN)
	if stm32.RNG.SR.HasBits(1) {
		return stm32.RNG.DR.Get()
	}
	return 0
}

// hw_init() is responsible for all hardware init
// Refs: https://stackoverflow.com/questions/63746239/enable-external-interrupts-arm-cortex-m0-stm32g070
func hw_init() {

	// SYSCFGEN is NEEDED FOR IRQ HANDLERS (button + Dio) .. Do not remove
	stm32.RCC.APB2ENR.SetBits(stm32.RCC_APB2ENR_SYSCFGEN)

	// BUTTON PB14
	machine.BUTTON.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	stm32.SYSCFG_COMP.EXTICR4.ReplaceBits(0b001, 0xf, stm32.SYSCFG_EXTICR4_EXTI14_Pos) // Enable PORTB On line 14
	stm32.EXTI.RTSR.SetBits(stm32.EXTI_RTSR_RT14)                                      // Detect Rising Edge of EXTI14 Line
	stm32.EXTI.FTSR.SetBits(stm32.EXTI_FTSR_FT14)                                      // Detect Falling Edge of EXTI14 Line
	stm32.EXTI.IMR.SetBits(stm32.EXTI_IMR_IM14)                                        // Enable EXTI14 line

	// GPIOS: Leds
	machine.LED_RED.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.LED_GREEN.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.LED_BLUE.Configure(machine.PinConfig{Mode: machine.PinOutput})
	// UART0 (Console)
	uartConsole = &machine.UART0
	uartConsole.Configure(machine.UARTConfig{TX: machine.UART_TX_PIN, RX: machine.UART_TX_PIN, BaudRate: 9600})
	go serial(uartConsole)
	/*
		// UART1 (GPS)
		uartGps = &machine.UART1
		uartGps.Configure(machine.UARTConfig{TX: machine.UART1_TX_PIN, RX: machine.UART1_TX_PIN, BaudRate: 9600})
		// GPS driver
		gps1 := gps.NewUART(uartGps)
		parser1 := gps.NewParser()
	*/
	// Spi bus configuration
	machine.SPI0.Configure(machine.SPIConfig{})
	// Lora sx1276/rfm95 configuration
	csPin := machine.PA15
	csPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	rstPin := machine.PB0
	rstPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	loraRadio = sx127x.New(machine.SPI0, csPin, rstPin)
	// Prepare Lora chil
	loraRadio.Init(loraConfig)
	// Configure PB13 (connected to DIO0.rfm95)
	machine.RFM95_DIO0_PIN.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	stm32.SYSCFG_COMP.EXTICR4.ReplaceBits(0b010, 0xf, stm32.SYSCFG_EXTICR4_EXTI13_Pos) // Enable PORTC On line 13
	stm32.EXTI.RTSR.SetBits(stm32.EXTI_RTSR_RT13)                                      // Detect Rising Edge of EXTI13 Line
	stm32.EXTI.FTSR.SetBits(stm32.EXTI_FTSR_FT13)                                      // Detect Falling Edge of EXTI13 Line
	stm32.EXTI.IMR.SetBits(stm32.EXTI_IMR_IM13)                                        // Enable EXTI13 line
	// Enable interrupts
	intr := interrupt.New(stm32.IRQ_EXTI4_15, gpios_int)
	intr.SetPriority(0x0)
	intr.Enable()

}

//----------------------------------------------------------------------------------------------//
//----------------------------------------------------------------------------------------------//
// main is where the program begins :-)
//----------------------------------------------------------------------------------------------//
//----------------------------------------------------------------------------------------------//
func main() {

	config = drivers.NewATConfig()

	// Packets will be sent to rxPktChan Channel
	rxPktChan = make(chan []byte)

	hw_init()

	println("*** Zaza Receiver ***")

	if e2p.ReadU8(0) == 0x0F {
		println("*** Eeprom contains 0x0F at pos 0\n***Loading configuration")
		//	config.Read(e2p)
	}

	println("Press ? for commands.")

	// 3 Green blinks at start
	for i := 0; i < 6; i++ {
		machine.LED_GREEN.Set((i % 2) == 0)
		time.Sleep(250 * time.Millisecond)
	}

	// Force RX
	//	processCmd("loratx 1de965a196b3")

	// Wait forever
	for {
		//		println("Button: ", machine.BUTTON.Get())
		time.Sleep(1 * time.Second)
	}

}
