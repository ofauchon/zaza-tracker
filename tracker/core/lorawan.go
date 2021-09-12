package core

import (
	"time"

	"github.com/ofauchon/zaza-tracker/libs"
)

// LoraWanTask() routing deals with the LoraWan
func LoraWanTask() {

	println("Lorawan Configuration ")
	println("  DevEUI : ", libs.BytesToHexString(loraStack.Otaa.DevEUI[:]))
	println("  AppEUI : ", libs.BytesToHexString(loraStack.Otaa.AppEUI[:]))
	println("  AppKey : ", libs.BytesToHexString(loraStack.Otaa.AppKey[:]))

	for {

		// Send join packet
		println("lorawan: Start joining")
		payload, err := loraStack.GenerateJoinRequest()
		if err != nil {
			println("lorawan: Error generating join request", err)
		}
		println("lorawan: UP_JOINREQUEST: ", libs.BytesToHexString(payload))

		// Send join
		loraTx(radio, payload)

		// Receive join Accept (Timeout 10s)
		resp, err := loraRx(radio, 10)
		if err != nil {
			println("lorawan: Error Rx ", err)
		}

		err = loraStack.DecodeJoinAccept(resp)
		if (err) == nil {
			println("Lorawan Network Joined !")
			println("  DevAddr: ", libs.BytesToHexString(loraStack.Session.DevAddr[:]), " (LSB)")
			println("  NetID  : ", libs.BytesToHexString(loraStack.Otaa.NetID[:]))
			println("  NwkSKey: ", libs.BytesToHexString(loraStack.Session.NwkSKey[:]))
			println("  AppSKey: ", libs.BytesToHexString(loraStack.Session.AppSKey[:]))
			// Sent sample message
			payload, err := loraStack.GenMessage(0, []byte("TinyGoLora"))
			if err == nil {
				println("TX_	UPMSG: --appkey ", libs.BytesToHexString(loraStack.Session.AppSKey[:]), " --nwkkey ", libs.BytesToHexString(loraStack.Session.NwkSKey[:]), " --hex", libs.BytesToHexString(payload))
				loraTx(radio, payload)
			} else {
				println(err)
			}
		} else {
			println("lorawan: Cant' decode join accept ", err)
		}

		// Wait 60s
		time.Sleep(time.Second * 60)
	} //for
}
