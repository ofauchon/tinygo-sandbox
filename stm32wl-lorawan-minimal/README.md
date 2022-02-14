# STM32WL Lorawan Connect

This simple application can be used with Nucleo WL55JC (STM32WL) boards. 
The code will try to join a Lorawan network and send a simple upload message.


# Configuration

Edit config_prod.go and replace with your DEVEUI/APPEUI/APPKEY.

# Build

  go flash -target=nucleo-wl55jc 

# Connect serial console 

  picocom -b 115200 /dev/ttyACM1
