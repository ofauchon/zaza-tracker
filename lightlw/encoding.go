package lightlw

import (
	"github.com/jacobsa/crypto/cmac"
)

func (r *LightLW) CalculateUplinkJoinMIC(micBytes []uint8, key [16]uint8) [4]uint8 {
	var mic [4]uint8

	hash, _ := cmac.New(key[:])
	hash.Write(micBytes)
	hb := hash.Sum([]byte{})

	copy(mic[:], hb[0:4])
	return mic
}
