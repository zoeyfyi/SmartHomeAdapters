from time import sleep
import serial
ser = serial.Serial('/dev/ttyACM0', 9600)
on = 0
while True:
	if on == 1: 
		on = 0 
	else:
		on = 1
	ser.write(str(chr(on)).encode())
	sleep(.1)
