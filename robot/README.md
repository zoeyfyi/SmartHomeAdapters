# Smart Home Adapters - Robot

## Setup 

1. Copy wifi.h into the esp8266 directory `cp wifi.h esp8266/`
2. Enter your SSID and password into wifi.h
3. If you are going to use eduroam, enter rest of the detail required by wifi.h

## Compile & Run

Use the Arduino IDE to open:
esp8266.ino for normal wifi
Eduroam_with_control for wpa2-enterprise PEAP

# Eduroam_Backup

The Eduroam_Backup.ino is the back up, it only provide the method to connect to eduroam.
