package core

import (
	"machine"
	"time"

	"github.com/ofauchon/zaza-tracker/libs"
)

var (
	TrackerLoopEnabled bool
	MonitorLoopEnabled bool
)

func TrackerLoopEnable() {
	TrackerLoopEnabled = true
}
func TrackerLoopDisable() {
	TrackerLoopEnabled = false
}

func MonitorLoopEnable() {
	MonitorLoopEnabled = true
}
func MonitorLoopDisable() {
	MonitorLoopEnabled = false
}

// ModeTracker waits for a valid position, and send a Lora packet
func TrackerLoop() {

	for {

		for TrackerLoopEnabled == false {
			time.Sleep(time.Second)
		}

		// Do we have a fix ?
		if currentState.lastValidFix.Valid {

			// Marshal to binary data
			pCurrentLocation := NewPoint(float64(currentState.lastValidFix.Longitude), float64(currentState.lastValidFix.Latitude))
			data, err := pCurrentLocation.MarshalBinary()

			if err == nil {
				pkt := append([]byte("zaza:"), data...)

				machine.LED.High()
				//err = loraTx(radio, pkt)
				time.Sleep(time.Second)
				machine.LED.Low()

				println("> TX: ", libs.BytesToHexString(pkt[:]))
			}
		}
		time.Sleep(time.Second * 15)
	}

}

// ModeReveive() Listen for Lora data, and helps locating target
func MonitorLoop() {
	/*
		pTarget := &Point{}
		for {

			for MonitorLoopEnabled == false {
				time.Sleep(time.Second)
			}

			// Try to get a packet
			println("wait for RX Packet (10s)")
			resp, err := loraRx(radio, 10000)
			//println("Rx done")
			if err != nil {
				println("RX Error:", err)
				time.Sleep(time.Second)
				continue
			}

			if resp == nil {
				//println("No data")
				continue
			}
			// LED Blink
			machine.LED.High()
			time.Sleep(time.Millisecond * 250)
			machine.LED.Low()
			println("> RX: ", libs.BytesToHexString(resp[:]))

			// Try to decode
			s := string(resp)
			if strings.HasPrefix(s, "zaza:") {
				err := pTarget.UnmarshalBinary(resp[5:])
				if err != nil {
					println("Can't unmarshall packet")
					continue
				}
			}
			fmt.Printf("> POS: %f %f \r\n", pTarget.Lat(), pTarget.Lng())

			if currentState.lastValidFix.Valid {
				pCurrent := NewPoint(float64(currentState.lastValidFix.Longitude), float64(currentState.lastValidFix.Latitude))

				// Display distance & bearing
				dist := pCurrent.GreatCircleDistance(pTarget)
				fmt.Printf("> Distance: %f\r\n", dist)
				bear := pCurrent.BearingTo(pTarget)
				fmt.Printf("> Bearing: %f\r\n", bear)
			}

		}

	*/

}
