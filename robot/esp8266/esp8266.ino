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
      digitalWrite(D7, LOW);
			Serial.printf("[WSc] Disconnected!\n");
			break;
		case WStype_CONNECTED:
      digitalWrite(D7, HIGH);
			Serial.printf("[WSc] Connected to url: %s\n", payload);
			break;
		case WStype_TEXT:
			Serial.printf("[WSc] get text: %s\n", payload);
      
      // led on
      if (strcmp((char *)payload, "led on") == 0) {
        digitalWrite(D8, HIGH);
      } 
      
      // led off
      if (strcmp((char *)payload, "led off") == 0) {
        digitalWrite(D8, LOW);
      }
      
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
  Serial.println(WiFi.localIP());

  digitalWrite(D7, HIGH);

  delay(1000);

  // connect to websocket server
  Serial.println("Connecting to WebSocket server");
  socket.begin("192.168.0.2", 8080, "/connect");
  socket.onEvent(webSocketEvent);
  socket.setReconnectInterval(1000);
}

void loop() { 
  socket.loop();
}
