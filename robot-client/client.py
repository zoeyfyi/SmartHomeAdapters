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
	print("Sending \"" + msg + "\" to arduino")
	ser.write(msg.encode('ascii', 'ignore'))

ws = websocket.WebSocket()
ws.on_open = on_open
ws = websocket.WebSocketApp("ws://0.0.0.0:8080/connect",
                              on_message = on_message,
                              on_error = on_error,
                              on_close = on_close)

ws.run_forever()

