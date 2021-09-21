package core

import (
	"fmt"
	"machine"
	"strings"
	"time"

	"github.com/ofauchon/zaza-tracker/libs"
)

// ModeTracker is used to GPS Position over Lora
func ModeTracker() {

	for {
		// Do we have a fix ?
		if currentState.lastValidFix.Valid {

			pCurrent := NewPoint(float64(currentState.lastValidFix.Longitude), float64(currentState.lastValidFix.Latitude))

			// Marshal and send
			data, err := pCurrent.MarshalBinary()
			if err == nil {

				pkt := append([]byte("zaza:"), data...)

				machine.LED1.High()
				err = loraTx(radio, pkt)
				time.Sleep(time.Second)
				machine.LED1.Low()

				println("TX: ", libs.BytesToHexString(pkt[:]))

			}

		}

		time.Sleep(time.Second * 15)

	}

}

// ModeReveive() Listen for Lora data, and helps locating target
func ModeReceive() {

	for {
		// Try to get a packet
		println("wait for RX Packet")
		resp, err := loraRx(radio, 10000)
		println("Rx done")
		if err != nil {
			println("RX Error:", err)
			time.Sleep(time.Second)
			continue
		}
		machine.LED2.High()
		time.Sleep(time.Millisecond * 250)
		machine.LED1.Low()

		println("RX: ", libs.BytesToHexString(resp[:]))

		// Decode it
		s := string(resp)
		if strings.HasPrefix(s, "zaza:") {
			pTarget := &Point{}
			err := pTarget.UnmarshalBinary(resp[5:])
			if err != nil {
				println("Can't unmarshall packet")
				continue
			}
			fmt.Printf("Position: %f %f \r\n", pTarget.Lat(), pTarget.Lng())

			if currentState.lastValidFix.Valid {
				pCurrent := NewPoint(float64(currentState.lastValidFix.Longitude), float64(currentState.lastValidFix.Latitude))

				// Display distance & bearing
				dist := pCurrent.GreatCircleDistance(pTarget)
				fmt.Printf("great circle distance: %f\r\n", dist)
				bear := pCurrent.BearingTo(pTarget)
				fmt.Printf("Bearing: %f\r\n", bear)
			}
		}

	}

}
