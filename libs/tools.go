package libs

// byteToHex return string hex representation of byte
func ByteToHex(b byte) string {
	bb := (b >> 4) & 0x0F
	ret := ""
	if bb < 10 {
		ret += string(rune('0' + bb))
	} else {
		ret += string(rune('A' + (bb - 10)))
	}

	bb = (b) & 0xF
	if bb < 10 {
		ret += string(rune('0' + bb))
	} else {
		ret += string(rune('A' + (bb - 10)))
	}
	return ret
}

// BytesToHexString converts byte slice to hex string representation
func BytesToHexString(data []byte) string {
	s := ""
	for i := 0; i < len(data); i++ {
		s += ByteToHex(data[i])
	}
	return s
}

//ByteToString returns string with binary representation of val
func IntToBinString(val int, bitsize int) string {
	s := ""
	for j := (bitsize - 1); j >= 0; j-- {
		if val&(1<<j) != 0 {
			s += "1"
		} else {
			s += "0"
		}
	}
	return s
}

//BytesToBinString returns binary string representation of byte array
func BytesToBinString(data []byte) string {
	s := ""
	for i := 0; i < len(data); i++ {
		s += IntToBinString(int(data[i]), 8)
		s += " "
	}
	return s
}

// revert inverts de order of a given byte slice
func Revert(s []byte) []byte {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}
