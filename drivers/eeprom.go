package drivers

import (
	"device/stm32"
	"unsafe"
)

const (
	EEPROM_ADDR_START uint32 = 0x08080000
	EEPROM_ADDR_END   uint32 = 0x08080FFF
)

type Eeprom struct {
}

// Unlock removes eeprom write protection
func (e *Eeprom) Unlock() {
	// Wait Flash not busy
	for stm32.Flash.SR.HasBits(stm32.Flash_SR_BSY) {
	}
	if stm32.Flash.PECR.HasBits(stm32.Flash_PECR_PELOCK) {
		stm32.Flash.PEKEYR.Set(0x89ABCDEF)
		stm32.Flash.PEKEYR.Set(0x02030405)
	}
}

// Lock enables eeprom write protection
func (e *Eeprom) Lock() {
}

func (e *Eeprom) ReadUint8(pos uint32) uint8 {
	return *(*uint8)(unsafe.Pointer(uintptr(0x08080000 + pos)))
}

func (e *Eeprom) WriteUint8(val uint8, pos uint32) {
	ptr := unsafe.Pointer(uintptr(EEPROM_ADDR_START + pos))
	*(*uint8)(ptr) = val
}
