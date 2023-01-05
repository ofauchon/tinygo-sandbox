
This TinyGo application was initialy written for testing sx126x and sx127x drivers interoperability.

AT-style commands sent on usart let you change Lora module modulation:

- Frequency (AT+FREQ)
- Spread Factor
- Coding Rate (AT+CR <x>)
- Header type (implicit/explicit)
- CRC (on/off)

Then you can send and receive Lora message with the defined configuration

# How to use

Just select your device (with -target) and the radio driver (with -tags): 

Then you can flash it:

```
# LoraE5 with sx126x (STM32WL's internal Lorawan radio)
tinygo flash -target=lorae5 -tags sx126x

# Bluepill with external sx127x module connected through SPI
tinygo flash -target=bluepill-clone -tags sx127x
```

Supported Lora modules : 

  * sx126x (Only tested on STM32WL5x SoCs that embedded sx126x)
  * sx127x (Present on various Lorawan modules : eg: RFM95)

