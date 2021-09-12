package core

import (
	"time"

	"tinygo.org/x/drivers/gps"
)

/*
func GpsEnable() {
	// GPS Standby_L (PB3)
	GPS_STANDBY_PIN.Configure(machine.PinConfig{Mode: machine.PinOutput})
	GPS_STANDBY_PIN.Set(false)
	// GPS Reset OFF (PB4)
	GPS_RESET_PIN.Configure(machine.PinConfig{Mode: machine.PinOutput})
	GPS_RESET_PIN.Set(false)
	// GPS Power ON (PB5)
	GPS_POWER_PIN.Configure(machine.PinConfig{Mode: machine.PinOutput})
	GPS_POWER_PIN.Set(true)
}

func GpsDisable() {
	// GPS Standby_L (PB3)
	GPS_STANDBY_PIN.Configure(machine.PinConfig{Mode: machine.PinOutput})
	GPS_STANDBY_PIN.Set(false)
	// GPS Reset OFF (PB4)
	GPS_RESET_PIN.Configure(machine.PinConfig{Mode: machine.PinOutput})
	GPS_RESET_PIN.Set(false)
	// GPS Power ON (PB5)
	GPS_POWER_PIN.Configure(machine.PinConfig{Mode: machine.PinOutput})
	GPS_POWER_PIN.Set(false)
}
*/
// gpsTask() processes incoming GPS sentences from the driver
func GpsTask(pGps gps.Device, pParser gps.Parser) {
	var fix gps.Fix
	println("GpsTask start")
	for {
		s, err := pGps.NextSentence()
		if err != nil {
			//println(err)
			continue
		}

		if (currentState.debug & DBG_GPS) > 0 {
			println("DGB:", s)
		}

		fix, err = pParser.Parse(s)
		if err != nil {
			//println(err)
			continue
		}

		if fix.Valid {
			currentState.lastValidFix = fix
		}

		time.Sleep(500 * time.Millisecond)

	}
}
