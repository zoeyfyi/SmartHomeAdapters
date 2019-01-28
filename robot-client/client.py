from time import sleep
import serial

import websocket

ser = serial.Serial('/dev/ttyACM0', 9600)

def on_open(ws):
	print("Websocket connected")

def on_error(ws, error):
	print("Websocket error", error)

def on_close(ws):
	print("Websocket closed")

def on_message(ws, msg):
	if msg == "led on":
		print("Sending 1")
		ser.write(str(chr(1)).encode())
	elif msg == "led off":
		print("Sending 0")
		ser.write(str(chr(0)).encode())
	else:
		print("Unknown message")

ws = websocket.WebSocket()
ws.on_open = on_open
ws = websocket.WebSocketApp("ws://0.0.0.0:8080/connect",
                              on_message = on_message,
                              on_error = on_error,
                              on_close = on_close)

ws.run_forever()

