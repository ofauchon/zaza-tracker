package core

import (
	"time"
	//"github.com/ofauchon/go-lorawan-stack"
	"github.com/ofauchon/zaza-tracker/libs"
)

// LoraWanTask() routing deals with the LoraWan
func LoraWanTask() {

	msg := []byte("TinyGoLora")
	r := GetRand32()
	loraStack.Otaa.DevNonce[0] = uint8(r & 0xFF)
	loraStack.Otaa.DevNonce[1] = uint8((r >> 8) & 0xFF)

	cpu := GetCpuId()
	println("Lorawan Configuration ")
	println("lorawan:  CPUID    : ", cpu[2], "/", cpu[1], "/", cpu[0])
	println("lorawan:  DevEUI   : ", libs.BytesToHexString(loraStack.Otaa.DevEUI[:]))
	println("lorawan:  AppEUI   : ", libs.BytesToHexString(loraStack.Otaa.AppEUI[:]))
	println("lorawan:  AppKey   : ", libs.BytesToHexString(loraStack.Otaa.AppKey[:]))
	println("lorawan:  DevNounce: ", libs.BytesToHexString(loraStack.Otaa.DevNonce[:]))

	for {

		// Send join packet
		println("lorawan: Start JOIN sequence")
		payload, err := loraStack.GenerateJoinRequest()
		if err != nil {
			println("lorawan: Error generating join request", err)
		}
		println("lorawan: Send JOIN request ", libs.BytesToHexString(payload))

		// Send join
		loraTx(radio, payload)

		println("lorawan: Wait for JOINACCEPT for 10s")
		// Receive join Accept (Timeout 10s)
		resp, err := loraRx(radio, 10000)
		if err != nil {
			println("lorawan: Error loraRx: ", err)
		}
		println("lorawan: Received a frame ")
		err = loraStack.DecodeJoinAccept(resp)
		if (err) == nil {
			println("lorawan: Valid JOINACCEPT, now connected")

			println("lorawan:   DevAddr: ", libs.BytesToHexString(loraStack.Session.DevAddr[:]), " (LSB)")
			println("lorawan:   NetID  : ", libs.BytesToHexString(loraStack.Otaa.NetID[:]))
			println("lorawan:   NwkSKey: ", libs.BytesToHexString(loraStack.Session.NwkSKey[:]))
			println("lorawan:   AppSKey: ", libs.BytesToHexString(loraStack.Session.AppSKey[:]))
			// Sent sample message
			payload, err := loraStack.GenMessage(0, msg)
			if err == nil {
				//println("TX_	UPMSG: --appkey ", libs.BytesToHexString(loraStack.Session.AppSKey[:]), " --nwkkey ", libs.BytesToHexString(loraStack.Session.NwkSKey[:]), " --hex", libs.BytesToHexString(payload))
				println("lorawan: Sending payload ", string(msg))
				err = loraTx(radio, payload)
				if err == nil {
					println("lorawan: loraTx OK")
				} else {
					println("lorawan: loraTx Error:", err)
				}
			} else {
				println("lorawan: Error building uplink message")
			}

		} else {
			println("lorawan: Cant' decode message (join accept expected) ", err)
		}

		// Wait 60s
		println("SLEEP 60s")
		time.Sleep(time.Second * 60)
	} //for
}
