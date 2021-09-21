package core

import (
	"errors"
	"machine"
	"time"

	"github.com/ofauchon/zaza-tracker/libs"
	"tinygo.org/x/drivers/lora/sx126x"
)

// configureLora initialize
func loraConfig(radio sx126x.Device) {

	println("Lora radio init")
	radio.SetStandby()

	radio.SetPacketType(sx126x.SX126X_PACKET_TYPE_LORA)
	radio.SetRfFrequency(LORA_FREQ)

	radio.SetBufferBaseAddress(0, 0)

	radio.ClearIrqStatus(sx126x.SX126X_IRQ_ALL)
	radio.SetDioIrqParams(sx126x.SX126X_IRQ_TX_DONE|sx126x.SX126X_IRQ_TIMEOUT|sx126x.SX126X_IRQ_RX_DONE, 0x00, 0x00, 0x00)

	radio.Calibrate(sx126x.SX126X_CALIBRATE_ALL)
	time.Sleep(10 * time.Millisecond)

	radio.SetCurrentLimit(60)

	radio.SetModulationParams(LORA_SF, sx126x.SX126X_LORA_BW_125_0, sx126x.SX126X_LORA_CR_4_7, sx126x.SX126X_LORA_LOW_DATA_RATE_OPTIMIZE_OFF)
	radio.SetPaConfig(0x04, 0x07, 0x00, 0x01)
	radio.SetTxParams(0x16, sx126x.SX126X_PA_RAMP_200U)
	radio.SetBufferBaseAddress(0, 0)

	radio.SetLoraPublicNetwork(true)

	radio.ClearDeviceErrors()
	radio.ClearIrqStatus(sx126x.SX126X_IRQ_ALL)

	// Configure RF GPIO
	// LoRa-E5 module ONLY transmits through RFO_HP:
	// Receive: PA4=1, PB5=0
	// Transmit(high output power, SMPS mode): PA4=0, PB5=1
	machine.PA4.Configure(machine.PinConfig{Mode: machine.PinOutput})
	machine.PB5.Configure(machine.PinConfig{Mode: machine.PinOutput})
}

// LoraTx
func loraTx(radio sx126x.Device, pkt []uint8) error {
	radio.ClearIrqStatus(sx126x.SX126X_IRQ_ALL)
	timeout := (uint32)(1000000 / 15.625) // 1sec

	// Set correct output (LoraE5 specific)
	machine.PA4.Set(false)
	machine.PB5.Set(true)

	// Define packet and modulation configuration (CRC ON, IQ OFF)
	radio.SetModulationParams(LORA_SF, sx126x.SX126X_LORA_BW_125_0, sx126x.SX126X_LORA_CR_4_7, sx126x.SX126X_LORA_LOW_DATA_RATE_OPTIMIZE_OFF)
	radio.SetPacketParam(LORA_PREAMBLE_TX, sx126x.SX126X_LORA_HEADER_EXPLICIT, sx126x.SX126X_LORA_CRC_ON, uint8(len(pkt)), sx126x.SX126X_LORA_IQ_STANDARD)

	// Copy and send packet
	radio.SetBufferBaseAddress(0, 0)
	radio.WriteBuffer(pkt)
	radio.SetTx(timeout)

	for {
		st := radio.GetStatus()
		irq := radio.GetIrqStatus()
		println("Status:", st, libs.IntToBinString(int(st), 8), "Irq:", libs.IntToBinString(int(irq), 8))

		if irq&sx126x.SX126X_IRQ_TX_DONE == sx126x.SX126X_IRQ_TX_DONE {
			return nil
		} else if irq&sx126x.SX126X_IRQ_TIMEOUT == sx126x.SX126X_IRQ_TIMEOUT {
			return errors.New("Tx timeout")
		} else if irq > 0 {
			println("IRQ value", irq)
			return errors.New("Unexpected IRQ value")
		}

		time.Sleep(time.Second * 1) // Check status every 100ms
	}
	return nil
}

// LoraRx() waits for incoming lora packet
// timeoutMs is the RX Timeout in milliseconds
// error returns error if any
func loraRx(radio sx126x.Device, timeoutMs int) ([]uint8, error) {

	radio.ClearIrqStatus(sx126x.SX126X_IRQ_ALL)

	//timeout := uint32(float32(timeoutMs) * 1000 / 15.625)

	// Wait RX
	machine.PA4.Set(true)
	machine.PB5.Set(false)

	// Define packet and modulation configuration (CRC OFF, IQ ON)
	radio.SetModulationParams(LORA_SF, sx126x.SX126X_LORA_BW_125_0, sx126x.SX126X_LORA_CR_4_7, sx126x.SX126X_LORA_LOW_DATA_RATE_OPTIMIZE_OFF)
	radio.SetPacketParam(LORA_PREAMBLE_RX, sx126x.SX126X_LORA_HEADER_EXPLICIT, sx126x.SX126X_LORA_CRC_OFF, 1, sx126x.SX126X_LORA_IQ_STANDARD)

	radio.SetRx(0) // RX Forever : FIXME, use RX timeouts

	elapsedMs := 0
	for elapsedMs < timeoutMs {
		//st := radio.GetStatus()
		//irq := radio.GetIrqStatus()

		//println("Status:", st, libs.IntToBinString(int(st), 8), "Irq:", irq, libs.IntToBinString(int(irq), 8))
		if irq&sx126x.SX126X_IRQ_RX_DONE == sx126x.SX126X_IRQ_RX_DONE {
			radio.ClearIrqStatus(sx126x.SX126X_IRQ_ALL)

			leng, offs := radio.GetRxBufferStatus()
			radio.SetBufferBaseAddress(0, offs) // Skip first byte
			pkt := radio.ReadBuffer(leng + 1)
			pkt = pkt[1:] // Skip first char ??? checkthat
			return pkt, nil
		} else if irq&sx126x.SX126X_IRQ_TIMEOUT == sx126x.SX126X_IRQ_TIMEOUT {
			return nil, nil
		} else if irq > 0 {
			println("IRQ value", irq)
			return nil, errors.New("RX:Unexpected IRQ value")
		}
		time.Sleep(time.Millisecond * 100) // Fixme, use interrupts
		elapsedMs += 100
	}
	return nil, nil

}
