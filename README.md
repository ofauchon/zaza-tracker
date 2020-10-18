# Description

ZAZA tracker project is an open-source alternative firmware for LGT92 Dragino Tracking device. It's 100% written in TinyGO. 

Two communication modes will be implemented:

  * Short range transmisson (Peer to peer Lora communication between two devices. This mode don't require Lorawan subscriptions)
  * Long range  mode (Lorawan network communication for long range between tracker and other devices : phones, computers.)
  * Dual mode :  Short and Long range communication at the same time

At the moment, It's composed of two modules you can flash independently on a LGT92 device :  

  * The "tracker" which determine location and send radio packets
  * The "receiver" which listen to radio packets and help locate the tracker


# Features (on 18/10/2020)


  * TinyGO port on LGT92/STM32L0x **DONE, not yet merged**
  * Serial console interface on LPUART1 **DONE** 
  * L70 GPS sentences decoding**DONE**
  * LGB Led support**DONE**
  * Button support **DONE** 
  * SX1276 LORA data transmission (peer to peer, no lorawan yet) **Work in progress**

# Pre-reqisites

  * A working TinyGO environnement with LGT92/STM32L0x support (*)
  * SWD interface (STLink or Bluepill) to flash the firmware on LGT92
  * TTL USART interface to connect to serial console (upcoming advanced features)
 
(*) Still work in progress, not yet merged in TinyGo official repositories

# How to use  

  * Build and run Tracker 

```
cd tracker/
Run "make" to build the firmaware
Run "make flash" to flash (with bluepill)
```

  * Build and run the Receiver 

```
cd receiver/
Run "make" to build the firmaware
Run "make flash" to flash (with bluepill)
Run "make term" to open serial console (with bluepill)
```

# Next steps 

  * Merge tracker and receiver application in a single one ? 
  * Implement low power modes to save battery 
  * Use interrupts for button events
  * Use EEPROM to save current configuration
  * Use accelerometer to improve location or moving detection
  * Tracker listening for notification packets
  * Over the air (Lora) firmware updates ? 
 


