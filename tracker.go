package main

import (
	"machine"
	"time"

	"github.com/ofauchon/zaza-tracker/core"
	//	"github.com/ofauchon/go-lorawan-stack"
	//	"tinygo.org/x/TRASH/drivers.lora.works/lora/sx127x"
	//core "./core"
)

var (
	//	loraRadio sx127x.Device
	//	loraStack lorawan.LoraWanStack

	cycle          uint32
	btnCount       uint32
	enableLoraInts bool
	currentMode    uint
)

// main is where the program begins :-)
func main() {

	// 3 Blinks at poweron
	machine.LED.Configure(machine.PinConfig{Mode: machine.PinOutput})
	for i := uint8(0); i < 6; i++ {
		machine.LED.Set((i % 2) == 0)
		time.Sleep(500 * time.Millisecond)

	}

	// Initialize basic hardware (UART, LED, BUTTONS)
	core.HwInit1()

	println("\n\n")
	println("################")
	println("# ZAZA TRACKER #")
	println("################")

	core.RunMode = core.RUNMODE_TRACKER
	if !core.BUTTON.Get() {
		core.RunMode = core.RUNMODE_RECEIVER
		println("XXX RECEIVER MODE XXX ")
		core.LED.Set(true)
	}

	if !core.BUTTON.Get() {
		core.RunMode = core.RUNMODE_CONSOLE
		println("XXX CONSOLE MODE XXX ")
		core.LED.Set(true)
	}

	// Initialize all hardware
	core.HwInit2()

	// Start Loops
	go core.ConsoleTaskLoop()
	go core.LoraWanTask()
	go core.GpsTaskLoop()
	go core.MonitorLoop()
	go core.TrackerLoop()

	for {
		time.Sleep(time.Millisecond * 250)
	}
}
