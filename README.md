# Description

Zaza tracker is a geolocation application for tinygo. 
It reads gps coordinate from GPS Modules and transmit to Lorawan TTN Netwotk


# Hardware. 

I'm currently developping with :

- LoraE5 developement board
- Beitian GPS Module (connected on UART2)
- Tinygo


# Build and run 

For developement, I use a custom go.mod file so I can relocate the packages to any location: 

It looks like this

```
go.mod

module zaza-tracker
replace tinygo.org/x/drivers => /your/workspace/tinygo-drivers
replace github.com/ofauchon/zaza-tracker => /your/workspace/zaza-tracker
replace  github.com/tinygo-org/go-cayenne-lib => /your/workspace/go-cayenne-lib-deadprogram
```

You'll need to define your OTAA Lorawan settings, you can create a mykeys.go in core folder :

```
package core

// These are sample keys, so the example builds
// Either change here, or create a new go file and use customkeys build tag
func setLorawanKeys() {
	otaa.SetAppEUI([]uint8{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
	otaa.SetDevEUI([]uint8{0xAA, 0xBB, 0xBB, 0xDD, 0xEE, 0x81, 0xB3, 0x65})
	otaa.SetAppKey([]uint8{0xA2, 0x45, 0xFC, 0xF8, 0x34, 0x22, 0x22, 0x34, 0x43, 0x43, 0x12, 0x4F, 0xCE, 0x7A, 0x32, 0x16})
}
```

then I build with : 

```
tinygo flash -work -x -target=lorae5 
```
 
