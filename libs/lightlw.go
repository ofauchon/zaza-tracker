package libs

import (
	"bytes"
	"crypto/aes"
	"encoding/binary"
	"errors"
	"math"

	"github.com/jacobsa/crypto/cmac"
	"github.com/ofauchon/zaza-tracker/drivers"
)

/*
This code was inspired from various projects:
https://github.com/BeelanMX/Beelan-LoRaWAN/blob/master/src/arduino-rfm/LoRaMAC.cpp
https://github.com/brocaar/lorawan
&https://github.com/trnlink/ls/blob/b5e69ea94d6650c9217997a5db4aa8a1cea68598/pkg/crypto/data_messages.go
http://www.techplayon.com/lora-device-activation-call-flow-join-procedure-using-otaa-and-abp/
*/

const (
	maxUploadPayloadSize = 220
)

type LoraEvent struct {
	eventType int
	eventData []byte
}

//LoraSession is used to store session data of a LoRaWAN session
type LoraSession struct {
	NwkSKey    [16]uint8
	AppSKey    [16]uint8
	DevAddr    [4]uint8
	FCntDown   uint32
	FCntUp     uint32
	CFList     [16]uint8
	RXDelay    uint8
	DLSettings uint8
}

//LoraOtaa is used to store session data of a LoRaWAN session
type LoraOtaa struct {
	DevEUI   [8]uint8
	AppEUI   [8]uint8
	AppKey   [16]uint8
	DevNonce [2]uint8
	AppNonce [3]uint8
	NetID    [3]uint8
}

//LoraMsg is used to store information of a LoRaWAN message to transmit or received
type LoraMsg struct {
	MacHeader    uint8
	DevAddr      [4]uint8
	FrameControl uint8
	FrameCounter uint16
	FramePort    uint8
	FrameOptions [15]uint8
	MIC          [4]uint8
	Direction    uint8
}

//LoraSettings  used for storing settings of the mote
type LoraSettings struct {
	Confirm        uint8 //0x00 Unconfirmed, 0x01 Confirmed
	Mport          uint8 //Port 1-223
	MoteClass      uint8 //0x00 Class A, 0x01 Class C
	DatarateTx     uint8 //See RFM file
	DatarateRx     uint8 //See RFM file
	ChannelTx      uint8 //See RFM file
	ChannelRx      uint8 //See RFM filed
	ChannelHopping uint8 //0x00 No hopping, 0x01 Hopping
	TransmitPower  uint8 //0x00 to 0x0F
}

const (
	ChanAll0     = uint8(0)
	ChanAll1     = uint8(1)
	ChanAll2     = uint8(2)
	ChanAll3     = uint8(3)
	ChanAll4     = uint8(4)
	ChanAll5     = uint8(5)
	ChanAll6     = uint8(6)
	ChanAll7     = uint8(7)
	ChanAll8     = uint8(8)
	ChanEu868RX2 = uint8(8)
	ChanMulti    = uint8(20)
)

type LightLW struct {
	Session  LoraSession
	Settings LoraSettings
	Otaa     LoraOtaa
}

// GenerateJoinRequest Generates a Lora Join request
func (r *LightLW) GenerateJoinRequest() ([]uint8, error) {
	// TODO: Add checks
	var buf []uint8
	buf = append(buf, 0x00)
	buf = append(buf, drivers.Revert(r.Otaa.AppEUI[:])...)
	buf = append(buf, drivers.Revert(r.Otaa.DevEUI[:])...)
	buf = append(buf, r.Otaa.DevNonce[:]...)
	mic := r.genPayloadMIC(buf, r.Otaa.AppKey)
	buf = append(buf, mic[:]...)
	return buf, nil
}

// DecodeJoinAccept Decodes a Lora Join Accept packet
func (r *LightLW) DecodeJoinAccept(phyPload []uint8) error {
	data := phyPload[1:] // Remove trailing 0x20
	// Prepare AES Cipher
	block, err := aes.NewCipher(r.Otaa.AppKey[:])
	if err != nil {
		return errors.New("Lora Cipher error 1")
	}
	buf := make([]byte, len(data))
	for k := 0; k < len(data)/aes.BlockSize; k++ {
		block.Encrypt(buf[k*aes.BlockSize:], data[k*aes.BlockSize:])
	}
	copy(r.Otaa.AppNonce[:], buf[0:3])
	copy(r.Otaa.NetID[:], buf[3:6])
	copy(r.Session.DevAddr[:], buf[6:10])
	r.Session.DLSettings = buf[10]
	r.Session.RXDelay = buf[11]

	if len(buf) > 16 {
		copy(r.Session.CFList[:], buf[12:28])
	}
	rxMic := buf[len(buf)-4:]

	dataMic := []byte{}
	dataMic = append(dataMic, phyPload[0])
	dataMic = append(dataMic, r.Otaa.AppNonce[:]...)
	dataMic = append(dataMic, r.Otaa.NetID[:]...)
	dataMic = append(dataMic, r.Session.DevAddr[:]...)
	dataMic = append(dataMic, r.Session.DLSettings)
	dataMic = append(dataMic, r.Session.RXDelay)
	dataMic = append(dataMic, r.Session.CFList[:]...)
	computedMic := r.genPayloadMIC(dataMic[:], r.Otaa.AppKey)
	if !bytes.Equal(computedMic[:], rxMic[:]) {
		return errors.New("Wrong Mic")
	}

	// Generate NwkSKey
	// NwkSKey = aes128_encrypt(AppKey, 0x01|AppNonce|NetID|DevNonce|pad16)
	sKey := []byte{}
	sKey = append(sKey, 0x01)
	sKey = append(sKey, r.Otaa.AppNonce[:]...)
	sKey = append(sKey, r.Otaa.NetID[:]...)
	sKey = append(sKey, drivers.Revert(r.Otaa.DevNonce[:])...)
	for i := 0; i < 7; i++ {
		sKey = append(sKey, 0x00) // PAD to 16
	}
	block.Encrypt(buf, sKey)
	copy(r.Session.NwkSKey[:], buf[0:16])

	// Generate AppSKey
	// AppSKey = aes128_encrypt(AppKey, 0x02|AppNonce|NetID|DevNonce|pad16)
	sKey[0] = 0x02
	block.Encrypt(buf, sKey)
	copy(r.Session.AppSKey[:], buf[0:16])

	// Reset counters
	r.Session.FCntDown = 0
	r.Session.FCntUp = 0

	return nil
}

// GenMessage Forge an uplink message
func (r *LightLW) GenMessage(dir uint8, payload []uint8) ([]uint8, error) {
	var buf []uint8
	buf = append(buf, 0b01000000) // FHDR Unconfirmed up

	buf = append(buf, drivers.Revert(r.Session.DevAddr[:])...)
	buf = append(buf, 0x00)                                                            // FCtl : No ADR, No RFU, No ACK, No FPending, No FOpt
	buf = append(buf, uint8((r.Session.FCntUp>>8)&0xFF), uint8(r.Session.FCntUp&0xFF)) // FCnt Up

	buf = append(buf, 0x01) // FPort=1

	fCnt := uint32(0)
	if dir == 0 {
		fCnt = r.Session.FCntUp
	} else {
		fCnt = r.Session.FCntDown
	}
	data, err := r.genFRMPayload(dir, drivers.Revert(r.Session.DevAddr[:]), fCnt, payload, false)
	if err != nil {
		return nil, err
	}
	buf = append(buf, data[:]...)

	// Mic
	//func (r *LightLW) calcMessageMIC(payload []uint8, key [16]uint8, dir uint8, addr []byte, fCnt uint32, lenMessage uint8) [4]uint8 {

	//println("[]byte used for mic: ", drivers.BytesToHexString(buf))
	mic := r.calcMessageMIC(buf, r.Session.NwkSKey, dir, drivers.Revert(r.Session.DevAddr[:]), fCnt, uint8(len(buf)))
	buf = append(buf, mic[:]...)

	return buf, nil
}

// encryptMessage encrypts Frame Header (Sec 4.3.3 lorawan 1.0.3 specification)
// dir : 0(uplink) 1(downlink)
// addr : devAddr
// fCnt : Frame counter (up or down)
func (r *LightLW) genFRMPayload(dir uint8, addr []uint8, fCnt uint32, payload []byte, isFOpts bool) ([]byte, error) {
	k := len(payload) / aes.BlockSize
	if len(payload)%aes.BlockSize != 0 {
		k++
	}
	if k > math.MaxUint8 {
		return nil, errors.New("Payload too big !")
	}
	encrypted := make([]byte, 0, k*16)
	cipher, err := aes.NewCipher(r.Session.AppSKey[:])
	if err != nil {
		panic(err) // types.AES128Key
	}

	var a [aes.BlockSize]byte
	a[0] = 0x01
	a[5] = dir
	copy(a[6:10], drivers.Revert(addr))
	binary.LittleEndian.PutUint32(a[10:14], fCnt)
	var s [aes.BlockSize]byte
	var b [aes.BlockSize]byte
	for i := uint8(0); i < uint8(k); i++ {
		copy(b[:], payload[i*aes.BlockSize:])
		if !isFOpts {
			a[15] = i + 1
		}
		cipher.Encrypt(s[:], a[:])
		for j := 0; j < aes.BlockSize; j++ {
			b[j] = b[j] ^ s[j]
		}
		encrypted = append(encrypted, b[:]...)
	}
	return encrypted[:len(payload)], nil
}

// getPayloadMIC computes MIC given the payload and the key
func (r *LightLW) genPayloadMIC(payload []uint8, key [16]uint8) [4]uint8 {
	var mic [4]uint8
	hash, _ := cmac.New(key[:])
	hash.Write(payload)
	hb := hash.Sum([]byte{})
	copy(mic[:], hb[0:4])
	return mic
}

// getPayloadMIC computes MIC given the payload and the key
func (r *LightLW) calcMessageMIC(payload []uint8, key [16]uint8, dir uint8, addr []byte, fCnt uint32, lenMessage uint8) [4]uint8 {

	var b0 [aes.BlockSize]byte

	copy(b0[0:5], []byte{0x49, 0x00, 0x00, 0x00, 0x00})
	b0[5] = dir
	copy(b0[6:10], drivers.Revert(addr)) // Test
	binary.LittleEndian.PutUint32(b0[10:14], fCnt)
	b0[14] = 0x00
	b0[15] = lenMessage

	var full []byte

	full = append(full, b0[:]...)
	full = append(full, payload...)

	//	println("calcMessageMic : B0+payload ", drivers.BytesToHexString(full))

	var mic [4]uint8
	hash, _ := cmac.New(key[:])
	hash.Write(full)
	hb := hash.Sum([]byte{})
	copy(mic[:], hb[0:4])
	return mic
}
