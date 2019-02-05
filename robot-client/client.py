import time
import serial
from serial.tools import list_ports

import websocket

def on_open(ws):
	print("Websocket connected")

def on_error(ws, error):
	print("Websocket error:", error)

def on_close(ws):
	print("Websocket closed")

def on_message(ws, msg):
	print("Sending \"" + msg + "\" to arduino")
	ser.write(msg.encode('ascii', 'ignore'))

while True:
	# Look for ports
	ports = list(list_ports.comports())
	print("Found", len(ports), "ports:")
	for p in ports:
		print("- ", p)

	# Try to connect to port
	for p in ports:
		print("Trying to connect to port:", p)
		try:
			ser = serial.Serial('/dev/ttyACM0', 9600)
			break
		except:
			print("/dev/ttyACM0 is not connected...")
			time.sleep(1)

	if ser == None:
		break

	print("Connected to:", ser.port)

	# Connect to websocket
	ws = websocket.WebSocket()
	ws.on_open = on_open
	ws = websocket.WebSocketApp("ws://robot.halspals.co.uk/connect",
								on_message = on_message,
								on_error = on_error,
								on_close = on_close)
	ws.run_forever()

