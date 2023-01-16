package core

import (
	"machine"

	"tinygo.org/x/drivers/gps"
	"tinygo.org/x/drivers/lora"
	"tinygo.org/x/drivers/lora/lorawan"
)

const (
	GPS_STANDBY_PIN = machine.PB0
	GPS_POWER_PIN   = machine.PB1
	GPS_RESET_PIN   = machine.PB2
	DBG_GPS         = 1
	//FREQ_LORA       = 434100000
	//	LORA_FREQ = 868100000

	LORA_FREQ        = 868100000
	LORA_SF          = 7
	LORA_PREAMBLE_TX = 12 // 12 symbols is default, 8 will be actually used
	LORA_PREAMBLE_RX = 12 // It should be the same for receiver and transmitter

	RUNMODE_CONSOLE  = 0
	RUNMODE_TRACKER  = 1
	RUNMODE_RECEIVER = 2
)

/*
	TX Test

868.1
SF8 / 150 / CR47 OK
SF6 OK
SF9 OK
SF12
*/
var (
	currentState                      status
	UartConsole, uartConsole, uartGps *machine.UART
	gps1                              gps.Device
	parser1                           gps.Parser

	RunMode int

	// Lora/Lorawan
	radio   lora.Radio
	session *lorawan.Session
	otaa    *lorawan.Otaa
)

type status struct {
	lastValidFix gps.Fix
	newData      bool
	debug        uint8
}

const (
	BUTTON    = machine.PA0
	LED       = machine.PB5
	ENABLE33V = machine.PA9
	ENABLE5V  = machine.PB10
)

func HwInit1() {

	//Enable 5V
	ENABLE5V.Configure(machine.PinConfig{Mode: machine.PinOutput})
	ENABLE5V.Set(true)

	//Enable 3.3V
	ENABLE33V.Configure(machine.PinConfig{Mode: machine.PinOutput})
	ENABLE33V.Set(true)

	// Buttons
	BUTTON.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	//machine.BTN2.Configure(machine.PinConfig{Mode: machine.PinInputPullup})
	//machine.BTN3.Configure(machine.PinConfig{Mode: machine.PinInputPullup})

	// LEDS
	LED.Configure(machine.PinConfig{Mode: machine.PinOutput})
	//machine.LED2.Configure(machine.PinConfig{Mode: machine.PinOutput})
	//machine.LED3.Configure(machine.PinConfig{Mode: machine.PinOutput})

	// UART CONSOLE
	// Console

	uartConsole = machine.UART0
	uartConsole.Configure(machine.UARTConfig{TX: machine.UART1_TX_PIN, RX: machine.UART1_RX_PIN, BaudRate: 9600})

	uartGps = machine.UART1
	uartGps.Configure(machine.UARTConfig{TX: machine.UART2_TX_PIN, RX: machine.UART2_RX_PIN, BaudRate: 9600})

}

func HwInit2() {

	// Radio
	//radio = sx126x.New(machine.SPI0)

	// Start LoraWan
	// DEVEUI : A84041000181B365
	// AppKey : 2C44FCF86C7B767B8FD3124FCE7A3216
	//loraStack.Otaa.AppEUI = [8]uint8{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}
	//	loraStack.Otaa.DevEUI = [8]uint8{0xA8, 0x40, 0x41, 0x00, 0x01, 0x81, 0xB3, 0x65}
	//	loraStack.Otaa.AppKey = [16]uint8{0x2C, 0x44, 0xFC, 0xF8, 0x6C, 0x7B, 0x76, 0x7B, 0x8F, 0xD3, 0x12, 0x4F, 0xCE, 0x7A, 0x32, 0x16}

	// Prepare HW for Lora
	//	loraConfig(radio)

}
