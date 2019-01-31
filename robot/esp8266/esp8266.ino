#include <Servo.h>

#include <WebSocketsClient.h>
#include <ESP8266WiFi.h>
#include "wifi.h"

#define SERVO_PIN D0
#define STATUS_PIN D7
  
WebSocketsClient socket;
WiFiClient client;

Servo servo;

void webSocketEvent(WStype_t type, uint8_t * payload, size_t length) {
	switch(type) {
		case WStype_DISCONNECTED:
      digitalWrite(STATUS_PIN, LOW);
			Serial.printf("[WSc] Disconnected!\n");
			break;
		case WStype_CONNECTED:
      digitalWrite(STATUS_PIN, HIGH);
			Serial.printf("[WSc] Connected to url: %s\n", payload);
			break;
		case WStype_TEXT:
			Serial.printf("[WSc] get text: %s\n", payload);

      {
        String command = String((char *) payload);
        command.trim();
        executeCommand(command);
      }
      
      break;
		case WStype_BIN:
			Serial.printf("[WSc] get binary length: %u\n", length);
			hexdump(payload, length);
      break;
	}
}

void executeCommand(String command) {
      Serial.print("Received command: ");
      Serial.print(command);
      Serial.println();

      // process LED commands
      if (command.startsWith("led")) {
        if(command == "led on") {
            Serial.println("Turning LED on");
            digitalWrite(LED_BUILTIN, LOW);
        } else if (command == "led off") {
            Serial.println("Turning LED off");
            digitalWrite(LED_BUILTIN, HIGH);
        }
      }

      if (command.startsWith("servo")) {
        // convert characters after "servo " to int
        int angle = command.substring(6).toInt();

        Serial.print("Setting servo to ");
        Serial.print(angle);
        Serial.println(" degrees");

        servo.write(angle);
      }
}

void setup() {
  // setup servo
  servo.attach(SERVO_PIN);
  
  // use pin 7/8 as our LED indicators
  pinMode(STATUS_PIN, OUTPUT);
  pinMode(LED_BUILTIN, OUTPUT);

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
    digitalWrite(STATUS_PIN, HIGH);
    delay(500);
    digitalWrite(STATUS_PIN, LOW);
    delay(500);
    Serial.print('-');
  }
  
  Serial.println('\n');
  Serial.println("Connection established!");  
  
  Serial.print("IP address:\t");
  Serial.println(WiFi.localIP());

  digitalWrite(STATUS_PIN, HIGH);

  delay(1000);

  // connect to websocket server
  Serial.println("Connecting to WebSocket server");
  socket.begin("192.168.0.12", 8080, "/connect");
  socket.onEvent(webSocketEvent);
  socket.setReconnectInterval(1000);
}

void loop() { 
  socket.loop();
}
