
** Update 17/01/2023 **

** Since a lot of Lora code had been moved to official tinygo-driver repository, this project needs a major cleanup/rewrite.  **

** Please have a look to the v2 branch for up-to-date news: **

https://github.com/ofauchon/zaza-tracker/tree/v2


# Description

ZAZA tracker project is an open-source alternative firmware for LGT92 Dragino Tracking device. It's 100% Golang/TinyGO code. 

Two communication modes will be implemented:

  * Lora mode : Mid-range positionning through some kind of Peer to peer Lora communication between two trackers)
  * Lorawan mode : For long distance tracking ( GPS position are pushed on a Lorawan Network)
  * Dual mode :  Short and Long range communication at the same time

At the moment, It's composed of two modules you can flash independently on a LGT92 device :  

  * The "tracker" which determine location and send radio packets
  * The "receiver" which listen to radio packets and help locate the tracker

The code is under heavy developpement and highly instable, please be kind. 

Please send me a mail at ofauchon2204@gmail.com if you want to work on the project. 

# Status (on 13/11/2020)

|Category|Task|Status|
|---|----|----|
|Hardware|TinyGO support for STM32L0x|  **DONE** but not yet merged in upstream Tinygo|
|Hardware|LPUART1 for serial console |**DONE** Serial communication OK| 
|Hardware|UART1 for L70 GPS|**DONE** : Serial communication OK|
|Hardware|RGB Led|**DONE** Individual LED control with GPIOs OK |
|Hardware|SX1276 SPI driver |**DONE**, SX1276 read/write register OK |
|Hardware|Push Button|**DONE**, Push and release events handled through external interrupt OK|
|Hardware|STM32L0 Eeprom support|**DONE** Eeprom Read/Write is OK. |
|Hardware|Reduce power consumption with Low power mode| **TO BE DONE** |
|Hardware|Watchdog and hardware reset| **TO BE DONE** |
|Radio|Lorawan lightweight stack| **WIP** First Lorawan Join Requests packet implementation OK. Next: Receive Join Accept packets and send real data  |
|Protocols|GPS Sentence decoding|**DONE** Getting a GPS Fix in about 30-60s|
|Protocols|serial console CLI and AT Commands| **TO BE DONE** |
|Other|Code cleanup| **TO BE DONE** |
|Other|Documentation (Build, flash, contribute)| **TO BE DONE** |

# How to run  

  * Prerequisite
```
  - A working TinyGO environnement with LGT92/STM32L0x support (*), with https://github.com/tinygo-org/tinygo/pull/1430 patch
  - SWD interface (STLink or BlackMagic) to flash the firmware on LGT92
  - TTL USART interface to connect to serial console 
 
(*) Still work in progress, not yet merged in TinyGo official repositories
```


  * Build and flash 

```
cd receiver/
Run "make" to build the firmware
Run "make flash" to flash (with BlackMagic Probe)
Run "make term" to open serial console (with blackMagic Probe)
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
 


