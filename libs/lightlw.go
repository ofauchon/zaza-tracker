package libs

import (
	"crypto/aes"

	"github.com/jacobsa/crypto/cmac"
)

/*
This code was inspired from various projects:
https://github.com/BeelanMX/Beelan-LoRaWAN/blob/master/src/arduino-rfm/LoRaMAC.cpp
https://github.com/brocaar/lorawan
https://runkit.com/avbentem/deciphering-a-lorawan-otaa-join-accept
*/

const (
	maxUploadPayloadSize = 220
)

//Struct used to store session data of a LoRaWAN session
type LoraSession struct {
	NwkSKey      [8]uint8
	AppSKey      [8]uint8
	DevAddr      [4]uint8
	FrameCounter uint16
	CFList       [16]uint8
	RXDelay      uint8
	DLSettings   uint8
}

type LoraOtaa struct {
	DevEUI   [8]uint8
	AppEUI   [8]uint8
	AppKey   [16]uint8
	DevNonce uint16
	AppNonce [3]uint8
	NetID    [3]uint8
}

//Struct to store information of a LoRaWAN message to transmit or received
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

//Struct used for storing settings of the mote
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

// GenerateJoinRequest Generates a Join request
/*
 *     Message Type = Join Request
 *           AppEUI = 70B3D57ED00000DC
 *           DevEUI = 00AFEE7CF5ED6F1E
 *         DevNonce = CC85
 *              MIC = 587FE913
 *
 * =>  00DC0000D07ED5B3701E6FEDF57CEEAF0085CC587FE913
 */
func (r *LightLW) GenerateJoinRequest() []uint8 {

	var rfmBuffer []uint8
	// Initialise message
	lmsg := &LoraMsg{}
	lmsg.MacHeader = 0x00
	lmsg.Direction = 0x00

	rfmBuffer = append(rfmBuffer, lmsg.MacHeader)

	// Load AppEUI
	for i := 0; i < 8; i++ {
		rfmBuffer = append(rfmBuffer, r.Otaa.AppEUI[7-i])
	}

	// Load DevEUI
	for i := 0; i < 8; i++ {
		rfmBuffer = append(rfmBuffer, r.Otaa.DevEUI[7-i])
	}

	// Load DevNounce
	if r.Otaa.DevNonce == 0 {
		println("Warning: LoraLW DevNonce=0")
	}
	rfmBuffer = append(rfmBuffer, uint8(r.Otaa.DevNonce&0x00FF))
	rfmBuffer = append(rfmBuffer, uint8((r.Otaa.DevNonce>>8)&0x00FF))

	// ADD Mic
	mic := r.CalculateUplinkJoinMIC(rfmBuffer, r.Otaa.AppKey)

	rfmBuffer = append(rfmBuffer, mic[:]...)

	return rfmBuffer
}

func (r *LightLW) EncryptAES(key []byte, data []byte) []byte {

	println("key size: ", len(key))
	println("data size: ", len(data))

	// create cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		println("Encrypt error 1")
		return nil
	}
	// allocate space for ciphered data
	out := make([]byte, len(data))

	/*
		iv := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
		mode := cipher.NewCBCEncrypter(block, iv)
		mode.CryptBlocks(out, data)
	*/
	// encrypt
	i := 0
	block.Encrypt(out[i:], data[i:]) //<<<<<<<<<<<<< On ne fait que le premier bloc ?
	i += aes.BlockSize
	block.Encrypt(out[i:], data[i:]) //<<<<<<<<<<<<< On ne fait que le premier bloc ?
	// return hex string
	return out
}

func revert(s []byte) []byte {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

// Decodes Join Accept Packet
/* Join Accept Payload : 204DD85AE608B87FC4889970B7D2042C9E72959B0057AED6094B16003DF12DE145
 * Application Key: B6B53F4A168A7A88BDF7EA135CE9CFCA
 * Corresponds to
 *     Message Type = Join Accept
 *         AppNonce = E5063A
 *            NetId = 000013
 *          DevAddr = 26012E43
 *       DLSettings = 03
 *          RXDelay = 01
 *           CFList = 184F84E85684B85E84886684586E8400
 *                  = decimal 8671000, 8673000, 8675000, 8677000, 8679000
 *              MIC = 55121DE0
 *
 * https://runkit.com/avbentem/deciphering-a-lorawan-otaa-join-accept
 *
 * The network server uses an AES decrypt operation in ECB mode to encrypt the join-accept
 * message so that the end-device can use an AES encrypt operation to decrypt the message.
 * This way an end-device only has to implement AES encrypt but not AES decrypt.

join received      :204dd85ae608b87fc4889970b7d2042c9e72959b0057aed6094b16003df12de145
join without 0x20: :  4dd85ae608b87fc4889970b7d2042c9e72959b0057aed6094b16003df12de145
join encoded:         3a06e5130000432e01260301184f84e85684b85e84886684586e840055121de0

Decoded:
// Size (bytes):     3       3       4         1          1     (16) Optional   4
// Join Accept:  AppNonce  NetID  DevAddr  DLSettings  RxDelay      CFList     MIC
*/
func (r *LightLW) DecodeJoinAccept(data []uint8) []uint8 {
	d := data[1:] // Remove trailing 0x20
	k := []byte(r.Otaa.AppKey[:])
	c := r.EncryptAES(k, d)

	copy(r.Otaa.AppNonce[:], revert(c[0:3]))
	copy(r.Otaa.NetID[:], revert(c[3:6]))
	copy(r.Session.DevAddr[:], revert(c[6:10]))
	r.Session.DLSettings = c[10]
	r.Session.RXDelay = c[11]

	if len(c) > 16 {
		copy(r.Session.CFList[:], c[12:28])
	}

	return c
}

//----------------------STOP HERE ---------------

func (r *LightLW) SendData(payload []uint8) {

	var rfmBuffer []uint8
	// Initialise message
	lmsg := &LoraMsg{}

	lmsg.MacHeader = 0x00
	lmsg.FramePort = 0x01
	lmsg.FrameControl = 0x00

	lmsg.DevAddr[0] = r.Session.DevAddr[0]
	lmsg.DevAddr[1] = r.Session.DevAddr[1]
	lmsg.DevAddr[2] = r.Session.DevAddr[2]
	lmsg.DevAddr[3] = r.Session.DevAddr[3]

	lmsg.Direction = 0x00
	lmsg.FrameCounter = r.Session.FrameCounter
	if r.Settings.Confirm == 0x00 {
		lmsg.MacHeader = lmsg.MacHeader | 0x40 // Unconfirmed
	} else {
		lmsg.MacHeader = lmsg.MacHeader | 0x80 // Confirmed
	}

	// Fill RFM Buffer
	rfmBuffer = append(rfmBuffer, lmsg.MacHeader)
	rfmBuffer = append(rfmBuffer, lmsg.DevAddr[3], lmsg.DevAddr[2], lmsg.DevAddr[1], lmsg.DevAddr[0])
	rfmBuffer = append(rfmBuffer, uint8(r.Session.FrameCounter&0x00FF))
	rfmBuffer = append(rfmBuffer, uint8(r.Session.FrameCounter>>8&0x00FF))

	if len(payload) > 0 {
		rfmBuffer = append(rfmBuffer, r.Settings.Mport)
	}
	/*

		EncryptPayload(payload, r.Session.AppSKey, lmsg)
			rfmBuffer[1] = lmsg.DevAddr[3]
			rfmBuffer[2] = lmsg.DevAddr[2]
			rfmBuffer[3] = lmsg.DevAddr[1]
			rfmBuffer[4] = lmsg.DevAddr[0]

			rfmBuffer[5] = lmsg.FrameControl

			rfmBuffer[6] = uint8(r.Session.FrameCounter & 0x00FF)
			rfmBuffer[7] = uint8((r.Session.FrameCounter >> 8) & 0xFF)
			rfmBuffer[]
	*/
}

// encryptPayload encrypts given payload
func (r *LightLW) encryptPayload(payload []uint8, lMsg LoraMsg) {

	i := uint8(0)
	j := uint8(0)
	numBlock := uint8(len(payload) / 16)
	incompleteBlockSize := len(payload) % 16
	var blocA [16]uint8

	if incompleteBlockSize > 0 {
		numBlock++
	}

	for i = 0; i < numBlock; i++ {
		blocA[0] = 0x01
		blocA[1] = 0x00
		blocA[2] = 0x00
		blocA[3] = 0x00
		blocA[4] = 0x00
		blocA[5] = lMsg.Direction
		blocA[6] = lMsg.DevAddr[3]
		blocA[7] = lMsg.DevAddr[2]
		blocA[8] = lMsg.DevAddr[1]
		blocA[9] = lMsg.DevAddr[0]
		blocA[10] = uint8(r.Session.FrameCounter & 0x00FF)
		blocA[11] = uint8((r.Session.FrameCounter >> 8) & 0x00FF)

		blocA[12] = 0x00
		blocA[13] = 0x00
		blocA[14] = 0x00

		blocA[15] = i + 1

		if i != (numBlock - 1) {
			for j = 0; j < 16; j++ {

			}
		}

	}

}

/* Compute MIC for Join Request
 * cmac = aes128_cmac(AppKey, MHDR | AppEUI | DevEUI | DevNonce)
 * MIC = cmac[0..3]
 */
func (r *LightLW) CalculateUplinkJoinMIC(micBytes []uint8, key [16]uint8) [4]uint8 {
	var mic [4]uint8

	hash, _ := cmac.New(key[:])
	hash.Write(micBytes)
	hb := hash.Sum([]byte{})

	copy(mic[:], hb[0:4])
	return mic
}
