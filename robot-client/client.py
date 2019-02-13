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


# Connect to websocket
ws = websocket.WebSocket()
ws.on_open = on_open
ws = websocket.WebSocketApp("ws://robot.halspals.co.uk/connect",
	on_message = on_message,
	on_error = on_error,
	on_close = on_close)
ws.run_forever()

