package drivers

import (
	"device/stm32"
	"unsafe"
)

const (
	EEPROM_ADDR_START = 0x08080000
	EEPROM_ADDR_END   = 0x08080FFF
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

// ReadU8 Reads uint8 value
func (e *Eeprom) ReadU8(pos uint32) uint8 {
	return *(*uint8)(unsafe.Pointer(uintptr(0x08080000 + pos)))
}

// ReadU8 Reads []uint8 array
func (e *Eeprom) ReadU8Array(pos, size int) []uint8 {
	var r []uint8
	for i := 0; i < size; i++ {
		v := *(*uint8)(unsafe.Pointer(uintptr(0x08080000 + pos + i)))
		r = append(r, v)
	}
	return r
}

// ReadU8 Write uint8
func (e *Eeprom) WriteU8(val uint8, pos uint32) {
	ptr := unsafe.Pointer(uintptr(EEPROM_ADDR_START + pos))
	*(*uint8)(ptr) = val
}

// WriteU8Array writes []uint8 in eeprom
func (e *Eeprom) WriteU8Array(val []uint8, pos int) {
	for i := 0; i < len(val); i++ {
		ptr := unsafe.Pointer(uintptr(EEPROM_ADDR_START + pos + i))
		*(*uint8)(ptr) = val[i]
	}
}
