package main

import (
	"device/stm32"
	"machine"
	"runtime/interrupt"
	"strconv"
	"strings"
	"time"

	"errors"

	"github.com/ofauchon/zaza-tracker/drivers"
	"github.com/ofauchon/zaza-tracker/libs"
	"tinygo.org/x/drivers/gps"
	"tinygo.org/x/drivers/lora/sx127x"
)

const (
	DBG_GPS = 1
)

type status struct {
	fix   *gps.Fix
	debug uint8
}

// Lorawan configuration
var (
	uartConsole, uartGps *machine.UART
	loraRadio            sx127x.Device
	st                   status
	send_data            = string("")
	send_delay           = int(0)
	cycle                uint32
	btnCount             uint32
	enableLoraInts       bool

	loraConfig = sx127x.Config{
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
)

// processCmd parses commands and execute actions
func processCmd(cmd string) error {
	ss := strings.Split(cmd, " ")
	switch ss[0] {
	case "help":
		println("reset: reset lora device")
		println("get: temp|mode|freq|regs")
		println("set: freq <8680000000> set transceiver frequency (in Hz)")
		println("mode: <rx,tx,standby,sleep>")
		println("debug: <gps,none> enable debug or none")

	case "reset":
		loraRadio.Reset()
		println("Reset done !")

	case "send":
		if len(ss) == 2 {
			println("Scheduled data to send :", ss[1])
			send_data = ss[1]
			/*
				err := d.Send([]byte(send_data))
				if err != nil {
					println("Send error", err)
				}
			*/
		}
	case "get":
		if len(ss) == 2 {
			switch ss[1] {
			case "freq":
				//println("Freq:")
			case "temp":
				//temp, _ := d.ReadTemperature(0)
				println("Temperature:")
			case "mode":
				//mode := d.GetMode()
				println(" Mode:")
			case "regs":
				println(" Regs:")
				/*
					for i := uint8(0); i < 0x60; i++ {
						val, _ := d.ReadReg(i)
						println(" Reg: ", strconv.FormatInt(int64(i), 16), " -> ", strconv.FormatInt(int64(val), 16))
					}
				*/
			default:
				return errors.New("Unknown command get")
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
			println("invalid use of set command")
		}

	case "mode":
		if len(ss) == 2 {
			switch ss[1] {
			case "standby":
				//loraRadio.Standby()
				println("TODO: Mode changed to Standby !")
			case "sleep":
				//loraRadio.Sleep()
				println("TODO: Mode changed to Sleep !")
			default:
				return errors.New("Unknown command mode")
			}
		}
	case "debug":
		if len(ss) == 2 {
			switch ss[1] {
			case "gps":
				st.debug |= DBG_GPS
			case "none":
				st.debug = 0
			default:
				return errors.New("Unknown command gps")
			}
		}
	default:
		return errors.New("Unknown command")
	}

	return nil
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

// gpsTask process incoming GPS sentences from the driver
func gpsTask(pGps gps.Device, pParser gps.Parser) {
	var fix gps.Fix
	println("Start gpsTask")
	for {
		s, err := pGps.NextSentence()
		if err != nil {
			//println(err)
			continue
		}

		if (st.debug & DBG_GPS) > 0 {
			println("DGB:", s)
		}

		fix, err = pParser.Parse(s)
		if err != nil {
			//println(err)
			continue
		}

		if fix.Valid {
			st.fix = &fix
		}

		time.Sleep(500 * time.Millisecond)

	}
}

func GpsEnable() {
	// GPS Standby_L (PB3)
	machine.GPS_STANDBY_PIN.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.GPS_STANDBY_PIN.Set(false)
	// GPS Reset OFF (PB4)
	machine.GPS_RESET_PIN.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.GPS_RESET_PIN.Set(false)
	// GPS Power ON (PB5)
	machine.GPS_POWER_PIN.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.GPS_POWER_PIN.Set(true)
}

func GpsDisable() {
	// GPS Standby_L (PB3)
	machine.GPS_STANDBY_PIN.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.GPS_STANDBY_PIN.Set(false)
	// GPS Reset OFF (PB4)
	machine.GPS_RESET_PIN.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.GPS_RESET_PIN.Set(false)
	// GPS Power ON (PB5)
	machine.GPS_POWER_PIN.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.GPS_POWER_PIN.Set(false)
}

// Interrupt handler from RFM95_DIO0 (RxDone Event) on PB13
func gpios_int(inter interrupt.Interrupt) {

	irqStatus := stm32.EXTI.PR.Get()

	if (irqStatus & stm32.EXTI_PR_PIF10) > 0 { // PC10 : DIO1 : RX_TMOUT
		loraRadio.DioIntHandler()
		stm32.EXTI.PR.Set(irqStatus)
	} else if (irqStatus & stm32.EXTI_PR_PIF13) > 0 { // PC13 : DIO0 : RX_DONE/TXDONE
		loraRadio.DioIntHandler()
		stm32.EXTI.PR.Set(irqStatus)
	} else if (irqStatus & stm32.EXTI_PR_PIF14) > 0 { // PB14 : Button
		btnCount++
		stm32.EXTI.PR.Set(irqStatus)
	}

}

func hw_init() {
	// SYSCFGEN is NEEDED FOR IRQ HANDLERS (button + Dio) .. Do not remove
	stm32.RCC.APB2ENR.SetBits(stm32.RCC_APB2ENR_SYSCFGEN)
	// BUTTON PB14 INTERRUPT
	machine.BUTTON.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})
	stm32.SYSCFG.EXTICR4.ReplaceBits(stm32.SYSCFG_EXTICR4_EXTI14_PB14, 0xf, stm32.SYSCFG_EXTICR4_EXTI14_Pos) // Enable PORTB On line 14
	stm32.EXTI.RTSR.SetBits(stm32.EXTI_RTSR_RT14)                                                            // Detect Rising Edge of EXTI14 Line
	//stm32.EXTI.FTSR.SetBits(stm32.EXTI_FTSR_FT14)                                 // Detect Falling Edge of EXTI14 Line
	stm32.EXTI.IMR.SetBits(stm32.EXTI_IMR_IM14) // Enable EXTI14 line
	// GPIOS: Leds
	machine.LED_RED.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.LED_GREEN.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.LED_BLUE.Configure(machine.PinConfig{Mode: machine.PinOutput})
	// UART0 (Console)
	uartConsole = &machine.UART0
	uartConsole.Configure(machine.UARTConfig{TX: machine.UART_TX_PIN, RX: machine.UART_TX_PIN, BaudRate: 9600})
	// UART1 (GPS)
	uartGps = &machine.UART1
	uartGps.Configure(machine.UARTConfig{TX: machine.UART1_TX_PIN, RX: machine.UART1_TX_PIN, BaudRate: 9600})
	// GPS driver
	gps1 := gps.NewUART(uartGps)
	parser1 := gps.NewParser()
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
	loraRadio.ConfigureLoraModem()

	// Configure interrupt for DIO0 (PC13) ... We watch after rising edge
	machine.RFM95_DIO0_PIN.Configure(machine.PinConfig{Mode: machine.PinInputFloating})
	stm32.SYSCFG.EXTICR4.ReplaceBits(stm32.SYSCFG_EXTICR4_EXTI13_PC13, 0xf, stm32.SYSCFG_EXTICR4_EXTI13_Pos) // Enable PORTC On line 13
	stm32.EXTI.FTSR.SetBits(stm32.EXTI_RTSR_RT13)                                                            // Detect Rising Edge of EXTI13 Line
	stm32.EXTI.IMR.SetBits(stm32.EXTI_IMR_IM13)                                                              // Enable EXTI13 line

	// Configure interrupt for DIO1 (PB10)... We watch after rising edge
	machine.RFM95_DIO1_PIN.Configure(machine.PinConfig{Mode: machine.PinInputFloating})
	stm32.SYSCFG.EXTICR3.ReplaceBits(stm32.SYSCFG_EXTICR3_EXTI10_PB10, 0xf, stm32.SYSCFG_EXTICR3_EXTI10_Pos) // Enable PORTC On line 10
	stm32.EXTI.FTSR.SetBits(stm32.EXTI_RTSR_RT10)                                                            // Detect Rising Edge of EXTI10 Line
	stm32.EXTI.IMR.SetBits(stm32.EXTI_IMR_IM10)                                                              // Enable EXTI10 line

	// Enable interrupts
	intr := interrupt.New(stm32.IRQ_EXTI4_15, gpios_int)
	intr.SetPriority(0x0)
	intr.Enable()

	// Start Serial console
	go serial(uartConsole)

	// Start GPS Device
	//GpsEnable()
	go gpsTask(gps1, parser1)

}

//-------- LORA ---------------

// Get Random 32bit
func getRand() uint32 {
	// Enable Random Number Generator Clock
	stm32.RCC.AHBENR.SetBits(stm32.RCC_AHBENR_RNGEN)

	// Enable PRNG periperal
	if !stm32.RNG.CR.HasBits(stm32.RNG_CR_RNGEN) {
		stm32.RNG.CR.SetBits(stm32.RNG_CR_RNGEN)
	}

	for !stm32.RNG.SR.HasBits(stm32.RNG_SR_DRDY) {
	}

	return stm32.RNG.DR.Get()
}

// loraJoin connects Lorawan network
/*
 *  Debug : ssh dragino
 * Logread -f

JOIN RQ
lora_pkt_fwd[1854]: PKT_FWD~ JOIN_REQ+ {"AppEUI":, "02010000000000A0", "DevEUI":, "01A00A000000ADDE"}
lora_pkt_fwd[1854]: RXTX~ {"rxpk":[{"tmst":1775631811,"time":"2021-01-15T13:35:12.236455Z","chan":1,"rfch":1,"freq":868.300000,"stat":1,"modu":"LORA","datr":"SF7BW125","codr":"4/5","lsnr":9.2,"rssi":-51,"size":23,"data":"AAIBAAAAAACgAaAKAAAArd44/YpABP8="}]}

JOIN ACCEPT
lora_pkt_fwd[1854]: PKT_FWD~ JOIN_ACCEPT+ {"NetID": "E43833", "DevAddr": "A014F9AB"}
lora_pkt_fwd[1854]: RXTX~ {"txpk":{"imme":false,"tmst":1780631811,"freq":868.3,"rfch":0,"powe":14,"modu":"LORA","datr":"SF7BW125","codr":"4/5","ipol":true,"size":33,"ncrc":true,"data":"IBkunuQ4M6v5FKDuKo2PUPk7ZSQNrpfEzjXarhHAJn/S"}}


*/
func loraJoin(config *drivers.ATConfig) {
	l := &libs.LightLW{}

	//l.Otaa.AppEUI = config.GetCurrentValue("APPEUI").([8]uint8)
	//l.Otaa.DevEUI = config.GetCurrentValue("DEVEUI").([8]uint8)
	//l.Otaa.AppKey = config.GetCurrentValue("APPKEY").([16]uint8)

	l.Otaa.AppEUI = [8]uint8{0xA0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x02}
	l.Otaa.DevEUI = [8]uint8{0xDE, 0xAD, 0x00, 0x00, 0x00, 0x0A, 0xA0, 0x01}
	l.Otaa.AppKey = [16]uint8{0x2C, 0x44, 0xFC, 0xF8, 0x6C, 0x7B, 0x76, 0x7B, 0x8F, 0xD3, 0x12, 0x4F, 0xCE, 0x7A, 0x32, 0x16}
	l.Otaa.DevNonce = uint16(getRand())

	// Send join packet
	j := l.GenerateJoinRequest()
	loraRadio.LoraSleep()
	loraRadio.TxLora([]byte(j))
	time.Sleep(1 * time.Second) // 1s should be enough for TX

	// Go Sleep and Receive
	//loraRadio.LoraSleep()
	loraRadio.GoReceive()

	go func() {
		ch := loraRadio.GetRxPktChannel()
		for {
			// Wait for incominl packet
			pk := <-ch
			println("Size:", len(pk))
			for i := 0; i < len(pk); i++ {
				print(drivers.ByteToHex(pk[i]), "")
			}
			println("")

			// Decode packet
			dec := l.DecodeJoinAccept(pk)
			for i := 0; i < len(dec); i++ {
				print(drivers.ByteToHex(dec[i]), "")
			}
			println("")

			time.Sleep(100 * time.Millisecond)
		}
	}() //Routine

}

//----------------------------------------------------------------------------------------------//

// main is where the program begins :-)
func main() {

	// Initialize all hardware
	hw_init()

	// 3 Blinks at poweron
	for i := uint8(0); i < 6; i++ {
		machine.LED_GREEN.Set((i % 2) == 0)
		time.Sleep(250 * time.Millisecond)

	}

	println("Zaza Tracker")

	config := drivers.NewATConfig()

	cycle = 1
	for {

		machine.LED_GREEN.Set(false)
		machine.LED_RED.Set(false)
		machine.LED_BLUE.Set(false)

		if (cycle % 10) == 0 {
			machine.LED_RED.Set(true)
		}
		if ((cycle + 1) % 10) == 0 {
			if st.fix != nil && st.fix.Valid {
				machine.LED_BLUE.Set(true)
			}
		}

		if (cycle % 30) == 0 {
			if st.fix != nil && st.fix.Valid {
				pkt := strconv.FormatFloat(float64(st.fix.Latitude), 'f', -1, 32) + ";"
				pkt += strconv.FormatFloat(float64(st.fix.Longitude), 'f', -1, 32) + ";"
				pkt += strconv.FormatFloat(float64(st.fix.Altitude), 'f', -1, 32) + ";"
				pkt += strconv.FormatFloat(float64(st.fix.Heading), 'f', -1, 32) + ";"
				println("Send Lora: ", pkt)

			}
		}

		if btnCount > 0 {
			btnCount = 0
			println("Button pressed, join lora now")
			loraJoin(config)
			//loraRadio.TxLora([]byte(pkt))

		}

		time.Sleep(1 * time.Second)
		cycle++
	}

}
