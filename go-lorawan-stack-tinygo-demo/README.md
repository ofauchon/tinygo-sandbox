
This is basic example of running go-lorawan-stack with tinygo.


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

