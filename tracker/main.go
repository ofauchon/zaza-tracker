package main

import (
	"device/stm32"
	"machine"
	"runtime/interrupt"
	"runtime/volatile"
	"strconv"
	"strings"
	"time"

	"errors"

	"tinygo.org/x/drivers/gps"
	"tinygo.org/x/drivers/lora/sx127x"
)

const (
	led    = machine.LED_RED
	button = machine.PB14
)

type status struct {
	fix   *gps.Fix
	debug uint8
}

const (
	DBG_GPS = 1
)

var loraRadio sx127x.Device

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
	stm32.EXTI.PR.Set(irqStatus)

	if (irqStatus & 0x2000) > 0 { // PC13 : DIO

		loraRadio.DioIntHandler()
	}
	if (irqStatus & 0x4000) > 0 { // PB14 : Button
		println("Button: ", machine.BUTTON.Get())
	}

}

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

	// Start Serial console
	go serial(uartConsole)

	// Start GPS Device
	GpsEnable()
	go gpsTask(gps1, parser1)

}

//----------------------------------------------------------------------------------------------//

// main is where the program begins :-)
func main() {

	hw_init()

	println("************")
	println("Zaza Tracker")
	println("************")

	var cycle uint32
	var btnCount uint32

	cycle = 1
	for {
		machine.LED_GREEN.Set(false)
		machine.LED_RED.Set(false)
		machine.LED_BLUE.Set(false)

		if (cycle % 5) == 0 {
			machine.LED_GREEN.Set(true)
		}

		if (cycle % 60) == 0 {
			if st.fix != nil && st.fix.Valid {
				machine.LED_BLUE.Set(true)

				pkt := strconv.FormatFloat(float64(st.fix.Latitude), 'f', -1, 32) + ";"
				pkt += strconv.FormatFloat(float64(st.fix.Longitude), 'f', -1, 32) + ";"
				pkt += strconv.FormatFloat(float64(st.fix.Altitude), 'f', -1, 32) + ";"
				pkt += strconv.FormatFloat(float64(st.fix.Heading), 'f', -1, 32) + ";"
				println("Send Lora: ", pkt)

				loraRadio.TxLora([]byte(pkt))

			}
		}

		if machine.BUTTON.Get() {
			btnCount++
		}

		if btnCount > 5 {
			machine.LED_RED.Set(!machine.LED_RED.Get())
			btnCount = 0
		}

		time.Sleep(1 * time.Second)
		cycle++
	}

}
