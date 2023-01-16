package core

import (
	"encoding/hex"
	"errors"
	"time"

	"tinygo.org/x/drivers/examples/lora/lorawan/common"
	"tinygo.org/x/drivers/lora"
	"tinygo.org/x/drivers/lora/lorawan"

	cayennelpp "github.com/TheThingsNetwork/go-cayenne-lib"
)

const (
	LORAWAN_JOIN_TIMEOUT_SEC    = 180
	LORAWAN_RECONNECT_DELAY_SEC = 15
	LORAWAN_UPLINK_DELAY_SEC    = 360
)

func failMessage(err error) {
	println("FATAL:", err)
	for {
	}
}

func loraConnect() error {
	start := time.Now()
	var err error
	for time.Since(start) < LORAWAN_JOIN_TIMEOUT_SEC*time.Second {
		println("Trying to join network")
		err = lorawan.Join(otaa, session)
		if err == nil {
			println("Connected to network !")
			return nil
		}
		println("Join error:", err, "retrying in", LORAWAN_RECONNECT_DELAY_SEC, "sec")
		time.Sleep(time.Second * LORAWAN_RECONNECT_DELAY_SEC)
	}

	err = errors.New("Unable to join Lorawan network")
	println(err.Error())
	return err
}

// LoraWanTask() routing deals with the LoraWan
func LoraWanTask() {

	// Board specific Lorawan initialization
	var err error
	radio, err = common.SetupLora()
	if err != nil {
		failMessage(err)
	}

	// Required for LoraWan operations
	session = &lorawan.Session{}
	otaa = &lorawan.Otaa{}

	// Initial Lora modulation configuration
	loraConf := lora.Config{
		Freq:           868100000,
		Bw:             lora.Bandwidth_125_0,
		Sf:             lora.SpreadingFactor9,
		Cr:             lora.CodingRate4_7,
		HeaderType:     lora.HeaderExplicit,
		Preamble:       12,
		Ldr:            lora.LowDataRateOptimizeOff,
		Iq:             lora.IQStandard,
		Crc:            lora.CRCOn,
		SyncWord:       lora.SyncPublic,
		LoraTxPowerDBm: 20,
	}
	radio.LoraConfig(loraConf)

	// Connect the lorawan with the Lora Radio device.
	lorawan.UseRadio(radio)

	// Configure AppEUI, DevEUI, APPKey
	setLorawanKeys()

	// Try to connect Lorawan network
	if err := loraConnect(); err != nil {
		failMessage(err)
	}

	encoder := cayennelpp.NewEncoder()

	// Try to periodicaly send an uplink sample message
	upCount := 1
	for {
		if currentState.lastValidFix.Latitude != 0 {
			println("lorawan: We have a fix to send")

			lat := float64(currentState.lastValidFix.Latitude)
			lon := float64(currentState.lastValidFix.Longitude)
			alt := float64(currentState.lastValidFix.Altitude)

			encoder.Reset()
			encoder.AddGPS(1, lat, lon, alt)
			cayBytes := encoder.Bytes()
			println("lorawan: Encoded payload:", hex.EncodeToString(cayBytes))

			if err := lorawan.SendUplink(cayBytes, session); err != nil {
				println("lorawan: Uplink transmission error:", err)
			} else {
				println("lorawan: Uplink transmission success")
			}
		} else {
			println("lorawan: No valid fix to send")
		}

		println("Sleeping for", LORAWAN_UPLINK_DELAY_SEC, "sec")
		time.Sleep(time.Second * LORAWAN_UPLINK_DELAY_SEC)
		upCount++
	}
}
