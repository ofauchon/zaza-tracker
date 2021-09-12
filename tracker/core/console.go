package core

import (
	"errors"
	"machine"
	"strconv"
	"strings"
	"time"
)

var (
	send_data  = string("")
	send_delay = int(0)
)

// processCmd() processes console commands
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
		//		loraRadio.Reset()
		//		println("Reset done !")

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

		/*
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
		*/
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
	case "show":
		if len(ss) == 2 {
			switch ss[1] {
			case "fix":
				if currentState.lastValidFix.Valid && currentState.lastValidFix.Valid {
					f := currentState.lastValidFix
					pkt := strconv.FormatFloat(float64(f.Latitude), 'f', -1, 32) + ";"
					pkt += strconv.FormatFloat(float64(f.Longitude), 'f', -1, 32) + ";"
					pkt += strconv.FormatFloat(float64(f.Altitude), 'f', -1, 32) + ";"
					pkt += strconv.FormatFloat(float64(f.Heading), 'f', -1, 32) + ";"
					println("fix:", pkt)
				} else {
					println("No fix yet")

				}
			default:
				return errors.New("Unknown command mode")

			}
		}

	case "debug":
		if len(ss) == 2 {
			switch ss[1] {
			case "gps":
				currentState.debug |= DBG_GPS
				println("gps debug =", currentState.debug)
			case "none":
				currentState.debug = 0
				println("gps debug =", currentState.debug)
			default:
				return errors.New("Unknown command gps")
			}
		}
	default:
		return errors.New("Unknown command")
	}

	return nil
}

// consoleTask receive and processes commands
func ConsoleTask(serialPort *machine.UART) string {
	println("ConsoleTask Start")
	input := make([]byte, 300) // serial port buffer

	i := 0

	for {

		if i == 300 {
			println("Serial Buffer overrun")
			i = 0
		}

		if serialPort.Buffered() > 0 {

			data, _ := serialPort.ReadByte() // read a character

			switch data {
			case 13: // pressed return key
				serialPort.Write([]byte("\r\n"))
				cmd := string(input[:i])
				processCmd(cmd)
				i = 0
			default: // pressed any other key
				serialPort.WriteByte(data)
				input[i] = data
				i++
			}
		}

		time.Sleep(10 * time.Millisecond)
	}

}
