package core

import (
	"device/stm32"
	"unsafe"
)

// getRand() returns Random 32bit
func GetRand32() uint32 {
	println("RNG: start Rand request")

	// Enable PRNG clock and peripheral
	stm32.RCC.AHB3ENR.SetBits(stm32.RCC_AHB3ENR_RNGEN)

	if !stm32.RNG.CR.HasBits(stm32.RNG_CR_RNGEN) {
		println("RNG: Device Disabled, re-enabling")
		if stm32.RNG.CR.HasBits(stm32.RNG_CR_NISTC) {
			// Disable NISTC
			println("RNG: Disable NIST Checks")
			stm32.RNG.CR.Set(stm32.RNG_CR_NISTC | stm32.RNG_CR_CONDRST)
		}
		// Enable RNG
		println("RNG: Enable RNG")
		stm32.RNG.CR.SetBits(stm32.RNG_CR_RNGEN)
		println("RNG: CR register = ", stm32.RNG.CR.Get())
	}
	// Wait for data ready
	for !stm32.RNG.SR.HasBits(stm32.RNG_SR_DRDY) {
		println("RNG: SR register =", stm32.RNG.SR.Get())
	}
	ret := stm32.RNG.DR.Get()
	println("RNG: Got random:", ret)
	return ret
}

func GetCpuId() [3]uint32 {
	var ret [3]uint32
	ret[0] = *(*uint32)(unsafe.Pointer(uintptr(0x1FFF7590)))
	ret[1] = *(*uint32)(unsafe.Pointer(uintptr(0x1FFF7594)))
	ret[2] = *(*uint32)(unsafe.Pointer(uintptr(0x1FFF7598)))
	return ret
}
