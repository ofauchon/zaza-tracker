package core

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

/*
lora setfreq 848300000
lora tx coucousssssssssssssssssss
*/
var (
	send_data  = string("")
	send_delay = int(0)
)

// processCmd() processes console commands
func processCmd(cmd string) error {
	ss := strings.Split(cmd, " ")
	switch ss[0] {
	case "help":
		println("mode [none|tracker|monitor]: Switch between tracker and monitor mode ")
		println("lora rx <duration millisec>: wait for rx packets for x millisec ")
		println("lora tx <message>: send lora packet ")
		println("lora setfreq <freq>: change frequency (in Hz) ")
		println("lorawan join")
		println("get: temp|mode|freq|regs")
		println("set: freq <868000000> set transceiver frequency (in Hz)")
		println("radiomode: <rx,tx,standby,sleep>")
		println("debug: <gps,none> enable debug or none")

	case "mode":
		if len(ss) == 2 {
			switch ss[1] {
			case "console":
				println("Start Console mode")
				TrackerLoopDisable()
				MonitorLoopDisable()
			case "monitor":
				println("Start Monitor mode")
				TrackerLoopDisable()
				MonitorLoopEnable()
			case "tracker":
				println("Start Monitor mode")
				MonitorLoopDisable()
				TrackerLoopEnable()
			}
		}

	case "get":
		if len(ss) == 2 {
			switch ss[1] {
			case "freq":
				println("Freq:")
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
func ConsoleTaskLoop() string {
	println("ConsoleTask Start")
	input := make([]byte, 300) // serial port buffer

	i := 0

	for {

		if i == 300 {
			println("Serial Buffer overrun")
			i = 0
		}

		if uartConsole.Buffered() > 0 {

			data, _ := uartConsole.ReadByte() // read a character

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
