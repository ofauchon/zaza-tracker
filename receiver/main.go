package main

import (
	"device/stm32"
	"fmt"
	"machine"
	"runtime/interrupt"
	"runtime/volatile"
	"strconv"
	"strings"
	"time"

	"tinygo.org/x/drivers/gps"
	"tinygo.org/x/drivers/lora/sx127x"
)

const (
	led    = machine.LED_RED
	button = machine.PB14
)

type status struct {
	fix *gps.Fix
}

var loraConfig = sx127x.Config{
	Frequency:       868000000,
	SpreadingFactor: 12,
	Bandwidth:       125000,
	CodingRate:      8,
	TxPower:         5,
	PaBoost:         true,
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

// processCmd parses commands and execute actions
func processCmd(cmd string) {
	ss := strings.Split(cmd, " ")
	switch ss[0] {
	case "help":
		println("reset: reset sx1276 device")
		println("loratx xxxxxxx: send 1 Lora packet every second until keypressed")
		println("lorarx : listen to lora packets until keypressed")
		println("get: sx1276config|regs")
		println("set: freq <433900000> set transceiver frequency (in Hz)")
		//println("mode: <rx,tx,standby,sleep>")

	case "reset":
		loraRadio.Reset()
		println("Reset done !")

	// Send 1 packet every 10 Seconds
	case "loratx":
		if len(ss) == 2 {
			keypressed = false
			go func() {
				cnt := int(0)
				for !keypressed {
					cnt++
					machine.LED_BLUE.Set(true)
					time.Sleep(250 * time.Millisecond)
					machine.LED_BLUE.Set(false)
					println("LoraTX Send: ", strconv.Itoa(cnt))
					loraRadio.SendPacket([]byte(strconv.Itoa(cnt)))
					time.Sleep(10 * time.Second)
				}
				println("LoraTX: Stopped by user")
			}()

		}
	// Listen for Lora packets for 20 Seconds
	case "lorarx":
		keypressed = false
		go func() {
			//Continuous mode leaves radio always on
			loraRadio.ReceiveContinuous()

			for !keypressed {

				packetSize := loraRadio.ParsePacket(0)
				if packetSize > 0 {
					//println("Got packet, RSSI=", loraRadio.LastPacketRSSI())
					size := loraRadio.ReadPacket(packet[:])
					println("RX: ", string(packet[:size]), " packetsize", packetSize)
				}

				time.Sleep(500 * time.Millisecond)
			}

			println("LoraRX: Stopped by user")
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

func Handler_EXTI4_15(inter interrupt.Interrupt) {
	println("IRQ: ", stm32.EXTI.PR.Get())
	loraRadio.ParsePacket(0)
}

//----------------------------------------------------------------------------------------------//

// main is where the program begins :-)
func main() {

	// Led
	machine.LED_RED.Configure(machine.PinConfig{Mode: machine.PinOutput})

	//Button
	machine.BUTTON.Configure(machine.PinConfig{Mode: machine.PinInputPulldown})

	// UART0 (Console)
	uartConsole = &machine.UART0
	uartConsole.Configure(machine.UARTConfig{TX: machine.UART_TX_PIN, RX: machine.UART_TX_PIN, BaudRate: 9600})
	// SPI and lx1276
	machine.SPI0.Configure(machine.SPIConfig{})
	csPin := machine.PA15
	csPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	rstPin := machine.PB0
	rstPin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	loraRadio = sx127x.New(machine.SPI0, csPin, rstPin)
	var err = loraRadio.Configure(loraConfig)
	if err != nil {
		fmt.Println(err)
		return
	}

	// RFM95 DIO0 Interrupt
	machine.RFM95_DIO0_PIN.Configure(machine.PinConfig{Mode: machine.PinInputFloating})
	//stm32.RCC.APB2ENR.SetBits(stm32.RCC_APB2ENR_AFIOEN)                                // Enable AFIO
	stm32.SYSCFG_COMP.EXTICR4.ReplaceBits(0b001, 0xf, stm32.SYSCFG_EXTICR4_EXTI13_Pos) // Enable PORTB EXTI only
	stm32.EXTI.RTSR.SetBits(stm32.EXTI_RTSR_RT13)                                      // Detect Rising Edge of EXTI0 Line
	stm32.EXTI.IMR.SetBits(stm32.EXTI_IMR_IM13)                                        // Enable EXTI0 line

	intr := interrupt.New(stm32.IRQ_EXTI4_15, Handler_EXTI4_15)
	intr.SetPriority(0xc0)
	intr.Enable()

	go serial(uartConsole)

	//-----------------------------------------------------
	// Start
	println("*************")
	println("Zaza Receiver")
	println("*************")

	println("Press help for commands.")

	// 3 Green blinks at start
	for i := 0; i < 3; i++ {
		machine.LED_GREEN.Set(true)
		time.Sleep(250 * time.Millisecond)
		machine.LED_GREEN.Set(false)
		time.Sleep(1 * time.Second)
	}

	processCmd("lorarx")

	for {
		time.Sleep(1 * time.Second)
	}

}
