package drivers

import "encoding/hex"

const (
	DEVEUI_POS = 10
	APPEUI_POS = 18
	APPKEY_POS = 26
)

type Config struct {
	DevEUI [8]uint8
	AppEUI [8]uint8
	AppKey [16]uint8
}

// Read() Loads configuration from eeprom
func (c *Config) Read(e Eeprom) {
	copy(c.DevEUI[:], e.ReadU8Array(DEVEUI_POS, 8))
	copy(c.AppEUI[:], e.ReadU8Array(APPEUI_POS, 8))
	copy(c.AppKey[:], e.ReadU8Array(APPKEY_POS, 16))

}

// Write() Saves configuration in eeprom
func (c *Config) Write(e Eeprom) {
	e.Unlock()
	e.WriteU8(0x0F, 0x00)
	e.WriteU8Array(c.DevEUI[:], DEVEUI_POS)
	e.WriteU8Array(c.AppEUI[:], APPEUI_POS)
	e.WriteU8Array(c.AppKey[:], APPKEY_POS)

}

// Test() Initialise sample configuration
func (c *Config) Default() {
	c.DevEUI = [8]uint8{0xA1, 0xA2, 0xA3, 0xA4, 0xA5, 0xA6, 0xA7, 0xA8}
	c.AppEUI = [8]uint8{0xB1, 0xB2, 0xB3, 0xB4, 0xB5, 0xB6, 0xB7, 0xB8}
	c.AppKey = [16]uint8{0xC1, 0xC2, 0xC3, 0xC4, 0xC5, 0xC6, 0xC7, 0xC8, 0xD1, 0xD2, 0xD3, 0xD4, 0xD5, 0xD6, 0xD7, 0xD8}
}

// Dump() displays configuration
func (c *Config) Dump() {
	println("Configuration Dump: ")
	println("DevEUI: ", hex.EncodeToString(c.DevEUI[:]))
	println("AppEUI: ", hex.EncodeToString(c.AppEUI[:]))
	println("AppKey: ", hex.EncodeToString(c.AppKey[:]))
}
