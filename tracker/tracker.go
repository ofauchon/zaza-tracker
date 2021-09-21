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
		time.Sleep(100 * time.Millisecond)

	}

	// Initialize all hardware
	core.HwInit1()

	println("\n\n")
	println("################")
	println("# ZAZA TRACKER #")
	println("################")

	core.RunMode = core.RUNMODE_TRACKER
	if !machine.BTN1.Get() {
		core.RunMode = core.RUNMODE_RECEIVER
		println("XXX RECEIVER MODE XXX ")
		machine.LED1.Set(true)
	}

	if !machine.BTN2.Get() {
		core.RunMode = core.RUNMODE_CONSOLE
		println("XXX CONSOLE MODE XXX ")
		machine.LED2.Set(true)
	}

	// Initialize all hardware
	core.HwInit2()

	switch core.RunMode {
	case core.RUNMODE_RECEIVER:
		println("Start RECEIVER mode")
		go core.ConsoleTask()
		//go core.GpsTask()
		core.ModeReceive()

	case core.RUNMODE_TRACKER:
		println("Start TRACKER mode")
		go core.ConsoleTask()
		go core.GpsTask()
		core.ModeTracker()

	case core.RUNMODE_CONSOLE:
		println("Start CONSOLE mode")
		go core.ConsoleTask()

	default:
		println("Nothing to do")

	}

	for {
		time.Sleep(time.Millisecond * 250)
	}
}
