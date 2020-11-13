# Description

ZAZA tracker project is an open-source alternative firmware for LGT92 Dragino Tracking device. It's 100% Golang/TinyGO code. 

Two communication modes will be implemented:

  * Lora mode : Mid-range positionning through some kind of Peer to peer Lora communication between two trackers)
  * Lorawan mode : For long distance tracking ( GPS position are pushed on a Lorawan Network)
  * Dual mode :  Short and Long range communication at the same time

At the moment, It's composed of two modules you can flash independently on a LGT92 device :  

  * The "tracker" which determine location and send radio packets
  * The "receiver" which listen to radio packets and help locate the tracker


# Status (on 13/11/2020)

|Task|Status|
|----|----|
|Build and run TinyGO on STM32L0x|  **DONE** but not yet merged in upstream Tinygo|
|Serial console interface on LPUART1 |**DONE**| 
|Power up, read serial, decode sentences from L70 GPS|**DONE**|
|RGB Led|**DONE**|
|Push Button|**DONE**, with GPIO interrupt handler|
|SX1276 minimal SPI driver |**DONE**, RX/TX packets OK|
|Eeprom to store configuration|**WIP**|
|Lorawan lightweight stack| **WIP**|
|Code cleanup| TODO |
|Documentation| TODO |

# How to run  

  * Prerequisite
```
  - A working TinyGO environnement with LGT92/STM32L0x support (*)
  - SWD interface (STLink or Bluepill) to flash the firmware on LGT92
  - TTL USART interface to connect to serial console (upcoming advanced features)
 
(*) Still work in progress, not yet merged in TinyGo official repositories
```


  * Build and flash 

```
cd receiver/
Run "make" to build the firmaware
Run "make flash" to flash (with bluepill)
Run "make term" to open serial console (with bluepill)
```

# Work in progress

|Task |
|-|
|Improve serial CLI commands
|Save non-volatile configuration in eeprom 
|Receiver/Tracker merge or code factorisation + cleanup| 
|Implement low power modes to save battery|
|Use accelerometer to improve location or moving detection|
|Tracker listening for notification packets|
|Over the air (Lora) firmware updates ? |
 


