package main

import (
	"fmt"
	"machine"
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
	fix *gps.Fix
}

var loraConfig = sx127x.Config{
	Frequency:       868000000,
	SpreadingFactor: 7,
	Bandwidth:       125000,
	CodingRate:      6,
	TxPower:         17,
}

var (
	st                   status
	lightLed             volatile.Register8
	uartConsole, uartGps *machine.UART
	send_data            = string("")
	send_delay           = int(0)
	packet               [255]byte
)

// processCmd parses commands and execute actions
func processCmd(cmd string) error {
	ss := strings.Split(cmd, " ")
	switch ss[0] {
	case "help":
		println("reset: reset rfm69 device")
		println("send xxxxxxx: send string over the air ")
		println("get: temp|mode|freq|regs")
		println("set: freq <433900000> set transceiver frequency (in Hz)")
		println("mode: <rx,tx,standby,sleep>")

	case "reset":
		//d.Reset()
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
				return errors.New("Unknown command mode")
			}
		}
	default:
		return errors.New("Unknown command")
	}

	return nil
}

// serial() function is a gorouting for handling USART rx data
func serial(serial *machine.UART) string {
	println("A0")
	input := make([]byte, 100) // serial port buffer
	println("A1")

	i := 0

	for {

		if i == 100 {
			println("Serial Buffer overrun")
			i = 0
		}

		if serial.Buffered() > 0 {

			data, _ := serial.ReadByte() // read a character

			switch data {
			case 10: // pressed return key
				println("A3")
				cmd := string(input[:i])
				println(cmd)
				i = 0
			default: // pressed any other key
				uartConsole.WriteByte(data)
				serial.WriteByte('A')
				input[i] = data
				i++
			}
		}

		time.Sleep(10 * time.Millisecond)
	}

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
	dio0Pin := machine.PC13
	dio0Pin.Configure(machine.PinConfig{Mode: machine.PinOutput})
	loraRadio := sx127x.New(machine.SPI0, csPin, rstPin)
	var err = loraRadio.Configure(loraConfig)
	if err != nil {
		fmt.Println(err)
		return
	}

	//-----------------------------------------------------
	// Start
	println("*************")
	println("Zaza Receiver")
	println("*************")

	println("Receiving LoRa packets...")

	for {
		packetSize := loraRadio.ParsePacket(0)
		if packetSize > 0 {
			println("Got packet, RSSI=", loraRadio.LastPacketRSSI())
			size := loraRadio.ReadPacket(packet[:])
			println(string(packet[:size]))
		}
	}

}
