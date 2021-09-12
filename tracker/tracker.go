package main

import (
	"machine"
	"time"

	//	"github.com/ofauchon/go-lorawan-stack"
	//	"tinygo.org/x/TRASH/drivers.lora.works/lora/sx127x"
	core "./core"
)

var (
	//	loraRadio sx127x.Device
	//	loraStack lorawan.LoraWanStack

	cycle          uint32
	btnCount       uint32
	enableLoraInts bool
)

// main is where the program begins :-)
func main() {

	// 3 Blinks at poweron
	machine.LED.Configure(machine.PinConfig{Mode: machine.PinOutput})
	for i := uint8(0); i < 6; i++ {
		machine.LED.Set((i % 2) == 0)
		time.Sleep(250 * time.Millisecond)

	}

	// Initialize all hardware
	core.HwInit()

	// Start GPS and Console routines
	core.StartTasks()

	for {
		time.Sleep(time.Second)
	}
}
