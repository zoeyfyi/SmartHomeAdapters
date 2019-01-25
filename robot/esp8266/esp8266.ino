#include <WebSocketsClient.h>
#include <ESP8266WiFi.h>
#include "wifi.h"

char path[] = "/";
char host[] = "echo.websocket.org";
  
WebSocketsClient socket;
WiFiClient client;

void webSocketEvent(WStype_t type, uint8_t * payload, size_t length) {
	switch(type) {
		case WStype_DISCONNECTED:
      digitalWrite(D8, LOW);
			Serial.printf("[WSc] Disconnected!\n");
			break;
		case WStype_CONNECTED: {
      digitalWrite(D8, HIGH);
			Serial.printf("[WSc] Connected to url: %s\n", payload);
      socket.sendTXT("Connected");
		}
			break;
		case WStype_TEXT:
			Serial.printf("[WSc] get text: %s\n", payload);
      break;
		case WStype_BIN:
			Serial.printf("[WSc] get binary length: %u\n", length);
			hexdump(payload, length);
      break;
	}
}

void setup() {
  // use pin 7/8 as our LED indicators
  pinMode(D7, OUTPUT);
  pinMode(D8, OUTPUT);

  // begin serial communications
  Serial.begin(9600);
  delay(10);
  Serial.println('\n');

  // connect to the network
  WiFi.begin(ssid, password);             
  Serial.print("Connecting to \"");
  Serial.print(ssid); 
  Serial.print("\"");
  Serial.println("...");

  int i = 0;
  while (WiFi.status() != WL_CONNECTED) {
    // wait for the Wi-Fi to connect
    digitalWrite(D7, HIGH);
    delay(500);
    digitalWrite(D7, LOW);
    delay(500);
    Serial.print('-');
  }
  
  Serial.println('\n');
  Serial.println("Connection established!");  
  
  Serial.print("IP address:\t");
  Serial.println(WiFi.localIP());         // Send the IP address of the ESP8266 to the computer

  digitalWrite(D7, HIGH);

  delay(5000);

  // connect to websocket server
  socket.begin("echo.websocket.org", 80, "/");
  socket.onEvent(webSocketEvent);
  socket.setReconnectInterval(5000);
}

void loop() { 
  socket.loop();
}
