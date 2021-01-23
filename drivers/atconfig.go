package drivers

const (
	AT_TYPE_CMD = 0
	AT_TYPE_PAR = 1
)

type ATConfigEntry struct {
	Class        uint8
	Short        string
	Desc         string
	CurrentValue interface{}
	FactoryValue interface{}
}

type ATConfig struct {
	Version         uint8
	ATConfigEntries []*ATConfigEntry
}

func NewATConfig() *ATConfig {
	r := ATConfig{}
	r.Version = 10 // v1.0
	r.ATConfigEntries = []*ATConfigEntry{
		// Configuration : Type, Description, Current value, Factory value
		&ATConfigEntry{AT_TYPE_PAR, "U8A", "Test: Get or Set array of 4 uint8", nil, [4]uint8{0x00, 0x00, 0x00, 0x00}},
		&ATConfigEntry{AT_TYPE_PAR, "U8", "Test:Get or Set uint8", nil, uint8(0)},
		&ATConfigEntry{AT_TYPE_PAR, "U32", "Test:Get or Set uint32", nil, uint32(0)},
		&ATConfigEntry{AT_TYPE_PAR, "APPEUI", "Get or Set the Application EUI", nil, [8]uint8{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		&ATConfigEntry{AT_TYPE_PAR, "APPKEY", "Get or Set the Application Key", nil, [8]uint8{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		&ATConfigEntry{AT_TYPE_PAR, "APPSKEY", "Get or Set the Application Session Key", nil, [16]uint8{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		&ATConfigEntry{AT_TYPE_PAR, "DADDR", "Get or Set the Device Address", nil, [4]uint8{0x00, 0x00, 0x00, 0x00}},
		&ATConfigEntry{AT_TYPE_PAR, "DEUI", "Get or Set the Device EUI", nil, [8]uint8{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		&ATConfigEntry{AT_TYPE_PAR, "NWKID", "Get or Set the Network ID", nil, [4]uint8{0x00, 0x00, 0x00, 0x00}},
		&ATConfigEntry{AT_TYPE_PAR, "NWKSKEY", "Get or Set the Network Session Key", nil, [16]uint8{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}},
		&ATConfigEntry{AT_TYPE_PAR, "CFM", "Get or Set the confirmation mode (0-1)", nil, uint8(0)},
		&ATConfigEntry{AT_TYPE_PAR, "TDC", "Get or set the application data transmission interval in ms", nil, uint32(30000)},
		&ATConfigEntry{AT_TYPE_PAR, "NJM", "Get or Set the Network Join Mode. (0: ABP, 1:OTAA)", nil, uint8(1)},
		&ATConfigEntry{AT_TYPE_PAR, "ADR", "Get or Set the Adaptive Data Rate setting. (0: off,1: on)", nil, uint8(0)},
		&ATConfigEntry{AT_TYPE_PAR, "CLASS", "Get or Set the ETSI Duty Cycle setting - 0=disable,1=enable - Only for testing", nil, uint8(0)},
		&ATConfigEntry{AT_TYPE_PAR, "DR", "Get or Set the Data Rate. (0-7 corresponding to DR_X)", nil, uint8(0)},
		&ATConfigEntry{AT_TYPE_PAR, "FCD", "Get or Set the Frame Counter Downlink", nil, uint8(0)},
		&ATConfigEntry{AT_TYPE_PAR, "FCU", "Get or Set the Frame Counter Uplink", nil, uint8(0)},
		&ATConfigEntry{AT_TYPE_PAR, "JN1DL", "Get or Set the Join Accept Delay between the end of the Tx and the Join Rx Window 1 in ms", nil, uint8(0)},
		&ATConfigEntry{AT_TYPE_PAR, "JN2DL", "Get or Set the Join Accept Delay between the end of the Tx and the Join Rx Window 2 in ms", nil, uint8(0)},
		&ATConfigEntry{AT_TYPE_PAR, "PNM", "Get or Set the public network mode. (0: off, 1: on)", nil, uint8(0)},
		&ATConfigEntry{AT_TYPE_PAR, "RX1DL", "Get or Set the delay between the end of the Tx and the Rx Window 1 in ms", nil, uint32(0)},
		&ATConfigEntry{AT_TYPE_PAR, "RX2DL", "Get or Set the delay between the end of the Tx and the Rx Window 2 in ms", nil, uint32(0)},
		&ATConfigEntry{AT_TYPE_PAR, "RX2DR", "Get or Set the Rx2 window data rate (0-7 corresponding to DR_X)", nil, uint8(0)},
		&ATConfigEntry{AT_TYPE_PAR, "RX2FQ", "Get or Set the Rx2 window frequency", nil, uint8(0)},
		&ATConfigEntry{AT_TYPE_PAR, "TXP", "Get or Set the Transmit Power (0-5, MAX:0, MIN:5,according to LoRaWAN Spec)", nil, uint8(0)},
		&ATConfigEntry{AT_TYPE_PAR, "CHS", "Get or Set Frequency (Unit: Hz) for Single Channel Mode", nil, uint32(0)},
		&ATConfigEntry{AT_TYPE_PAR, "ACE", "Get or set the Alarm data transmission interval in ms", nil, uint32(60000)},
		&ATConfigEntry{AT_TYPE_PAR, "KAT", "Get or set the keep alive time interval in ms", nil, uint32(21600000)},
		&ATConfigEntry{AT_TYPE_PAR, "LON", "Get or set the LED flashing of position, downlink and uplink", nil, uint8(1)},
		&ATConfigEntry{AT_TYPE_PAR, "MLON", "Get or set the LED of movement detection (Disable(0), Enable (1))", nil, uint8(0)},
		&ATConfigEntry{AT_TYPE_PAR, "MD", "Get or set the mode of motion detection (0:Disable,1:Move,2:Collide,3:Customized)", nil, uint8(1)},
		&ATConfigEntry{AT_TYPE_PAR, "PDOP", "Get or set the PDOP value", nil, uint8(0)},
		&ATConfigEntry{AT_TYPE_PAR, "FTIME", "Set max GPS positioning time until GPS poweroff", nil, uint8(0)},
		&ATConfigEntry{AT_TYPE_PAR, "NMEA353", "Get or set the search mode of GPS (For L76-L only)", nil, uint8(0)},
		// Command : Type, Description, Current value, Factory value
		&ATConfigEntry{AT_TYPE_CMD, "Z", "Trigger a reset of the MCU", nil, nil},
		&ATConfigEntry{AT_TYPE_CMD, "VER", "Get current image version and Frequency Band", nil, nil},
		&ATConfigEntry{AT_TYPE_CMD, "RSSI", "Get the RSSI of the last received packet", nil, nil},
		&ATConfigEntry{AT_TYPE_CMD, "SNR", "Get the SNR of the last received packet", nil, nil},
		&ATConfigEntry{AT_TYPE_CMD, "SEND", "Send text data along with the application port", nil, nil},
		&ATConfigEntry{AT_TYPE_CMD, "SENDB", "Send hexadecimal data along with the application port", nil, nil},
		&ATConfigEntry{AT_TYPE_CMD, "HWVER", "Get the LGT92 of hardware version and gps of version.", nil, nil},
		&ATConfigEntry{AT_TYPE_CMD, "FDR", "Reset Parameters to Factory Default, Keys Reserve", nil, nil},
		&ATConfigEntry{AT_TYPE_CMD, "JOIN", "Join network", nil, nil},
		&ATConfigEntry{AT_TYPE_CMD, "NJS", "Get the join status", nil, nil},
		&ATConfigEntry{AT_TYPE_CMD, "RECV", "Print last received data in raw format", nil, nil},
		&ATConfigEntry{AT_TYPE_CMD, "RECVB", "Print last received data in binary format (with hexadecimal values)", nil, nil}}
	return &r
}

//PackConfig returns a []byte representation of the configuration
func (a *ATConfig) PackConfig() []byte {
	var buf []uint8

	for _, v := range a.ATConfigEntries {
		println("Packing param ", v.Short)
		if v.Class == AT_TYPE_PAR {
			if v.CurrentValue != nil {
				switch v.CurrentValue.(type) {
				case [4]uint8:
					val := v.CurrentValue.([4]uint8)
					buf = append(buf, val[:]...)
				case [8]uint8:
					val := v.CurrentValue.([8]uint8)
					buf = append(buf, val[:]...)
				case [16]uint8:
					val := v.CurrentValue.([16]uint8)
					buf = append(buf, val[:]...)

				case uint8:
					val := v.CurrentValue.(uint8)
					buf = append(buf, val)

				case uint32:
					bb := v.CurrentValue.(uint32)
					buf = append(buf, uint8((bb>>24)&0xF))
					buf = append(buf, uint8((bb>>16)&0xF))
					buf = append(buf, uint8((bb>>8)&0xF))
					buf = append(buf, uint8((bb>>0)&0xF))
				}
			}
		}
	}
	// Preprend some headers
	// [MAGIC_0x0F][VERSION_0x01][LEN_HI][LEN_LOW][.......PAYLOAD........]
	sz := len(buf)
	buf = append([]uint8{uint8((sz & 0xFF) >> 8), uint8(sz & 0xFF)}, buf...) // Prepend 16 bit size
	buf = append([]uint8{0xF1, 0x1F, 0x01}, buf...)                          // Prepend magic (0x0F) and version (0x01)

	return buf
}

//UnPackConfig returns a []byte representation of the configuration
func (a *ATConfig) UnPackConfig(data []uint8) {
	pos := int(0)

	if data[pos] == 0xF1 && data[pos+1] == 0x1F {
		if data[pos+2] == 1 {

		} else {
			println("Invalid Config: version should be 1")
		}
	} else {
		println("Invalid Config: bad magic")
	}

	pos = 3
	for _, v := range a.ATConfigEntries {
		if v.Class == AT_TYPE_PAR {
			switch v.FactoryValue.(type) {
			case [4]uint8:
				v.CurrentValue = data[pos : pos+4]
				pos += 4
			case [8]uint8:
				v.CurrentValue = data[pos : pos+8]
				pos += 8
			case [16]uint8:
				v.CurrentValue = data[pos : pos+16]
				pos += 16

			case uint8:
				v.CurrentValue = data[pos]
				pos += 1

			case uint32:
				a := uint32(data[pos])
				a = (a << 8) | uint32(data[pos+1])
				a = (a << 8) | uint32(data[pos+2])
				a = (a << 8) | uint32(data[pos+3])
				a = (a << 8) | uint32(data[pos+4])
				v.CurrentValue = a
			}
		}
	}
}

// DumpConfig returns string with current config
func (a *ATConfig) DumpConfig() string {

	var st string
	for _, v := range a.ATConfigEntries {
		if v.Class == AT_TYPE_PAR {
			if v.CurrentValue != nil {
				st += v.Short + ": "
				switch v.CurrentValue.(type) {
				case []uint8:
					bb := v.FactoryValue.([]uint8)
					for i := 0; i < len(bb); i++ {
						st += ByteToHex(bb[i])
					}
					st += "\n"
				case uint8:
					bb := v.CurrentValue.(uint8)
					st += ByteToHex(bb) + "\n"
				case uint32:
					bb := v.CurrentValue.(uint32)
					st += "0x" + ByteToHex(uint8((bb>>24)&0xF))
					st += ByteToHex(uint8((bb >> 16) & 0xF))
					st += ByteToHex(uint8((bb >> 8) & 0xF))
					st += ByteToHex(uint8(bb&0xF)) + "\n"
				}
			}
		}
	}
	return st
}

func (a *ATConfig) GetConfigEntry(v string) *ATConfigEntry {
	for _, a := range a.ATConfigEntries {
		if a.Class == AT_TYPE_PAR && a.Short == v {
			return a
		}
	}
	return nil
}

// Returns current value of given configuration parameter
func (a *ATConfig) GetCurrentValue(v string) interface{} {
	r := a.GetConfigEntry(v)
	if r != nil {
		return r.CurrentValue
	}
	return nil
}

func (a *ATConfig) GetCommandDesc(s string) string {
	for _, v := range a.ATConfigEntries {
		if v.Short == s {
			return v.Desc
		}
	}
	return "Unknown command"
}

func (a *ATConfig) FactoryReset() {
	for _, v := range a.ATConfigEntries {
		v.CurrentValue = v.FactoryValue
	}
}
