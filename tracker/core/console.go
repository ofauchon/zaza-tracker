package core

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/ofauchon/zaza-tracker/libs"
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
		println("lora rx <duration millisec>: wait for rx packets for x millisec ")
		println("lora tx <message>: send lora packet ")
		println("lora setfreq <freq>: change frequency (in Hz) ")
		println("lorawan join")
		println("get: temp|mode|freq|regs")
		println("set: freq <868000000> set transceiver frequency (in Hz)")
		println("mode: <rx,tx,standby,sleep>")
		println("debug: <gps,none> enable debug or none")

	// LORA COMMANDS
	case "lora":

		if len(ss) < 2 {
			break
		}

		switch ss[1] {
		case "tx":
			if len(ss) == 3 {
				txdata := []byte(ss[2])
				println("TX Data :", libs.BytesToHexString(txdata))
				err := loraTx(radio, txdata)
				if err != nil {
					println("tx error", err)
				} else {
					println("OK")
				}
			}

		case "rx":
			if len(ss) == 3 {
				timeout, _ := strconv.Atoi(ss[2])
				println("Lora RX for ", timeout, "ms ")
				if uartConsole.Buffered() > 0 {
					println("Stopped by user")
					break
				}
				println("Waiting for RX Packet with timeout", timeout, " ms.")
				data, err := loraRx(radio, timeout)
				if err != nil {
					println("RX: ", libs.BytesToHexString(data))
					println("OK")
				} else {
					println(err)
				}

			}

		case "setfreq":
			if len(ss) == 3 {
				freq, err := strconv.Atoi(ss[2])
				if err != nil {
					println("Bad frequence", err)
				} else {
					println("> Switch to standby and change freq")
					radio.SetStandby()
					radio.SetRfFrequency(uint32(freq))
					println("OK")
				}
			}

		//radio.SetModulationParams(LORA_SF, sx126x.SX126X_LORA_BW_125_0, sx126x.SX126X_LORA_CR_4_7, sx126x.SX126X_LORA_LOW_DATA_RATE_OPTIMIZE_OFF)
		case "setmod":
			if len(ss) == 6 {
				sf, _ := strconv.Atoi(ss[2])
				bw, _ := strconv.Atoi(ss[3])
				cr, _ := strconv.Atoi(ss[4])
				ldr, _ := strconv.Atoi(ss[5])
				println("> Switch to standby")
				radio.SetStandby()
				println("> SetModulationParam:", uint8(sf), uint8(bw), uint8(cr), uint8(ldr))
				radio.SetModulationParams(uint8(sf), uint8(bw), uint8(cr), uint8(ldr))
				println("OK")
			}
		}

	case "lorawan":
		if len(ss) == 2 {
			switch ss[1] {
			case "join":
				println("Start Join")
				// Send join packet
				payload, _ := loraStack.GenerateJoinRequest()
				loraTx(radio, payload)
				println("Wait Join Response")
				resp, err := loraRx(radio, 10000)
				err = loraStack.DecodeJoinAccept(resp)
				if err == nil {
					println("Join Accept Response OK")
				}
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
func ConsoleTask() string {
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
