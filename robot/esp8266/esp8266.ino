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
        executeCommandSequence(command);
      }
      
      break;
		case WStype_BIN:
			Serial.printf("[WSc] get binary length: %u\n", length);
			hexdump(payload, length);
      break;
	}
}

void executeCommandSequence(String command) {
      const char* delim = ";";
      char *cmdtok = strtok(const_cast<char*>(command.c_str()), delim);
      while (cmdtok != nullptr) {
        String cmd = String(cmdtok);
        
        if(cmd == "led on") {
          Serial.println("Turning LED on");
          digitalWrite(LED_BUILTIN, LOW);
        } else if (cmd == "led off") {
          Serial.println("Turning LED off");
          digitalWrite(LED_BUILTIN, HIGH);
        } else if (cmd.startsWith("servo")) {
          // convert characters after "servo " to int
          int angle = cmd.substring(6).toInt();
  
          Serial.print("Setting servo to %d degrees\n", angle);

          // write to the servo so it moves to this angle when it is attached
          servo.write(angle);

          // attach the servo
          servo.attach(SERVO_PIN);
          delay(500); // wait for movement

          // this cuts the signal to the servo, and prevents it from moving after the detach
          // this took _way_ to long to figure out
          servo.writeMicroseconds(0);
          delay(500);

          // detach the servo
          servo.detach();
        } else if (cmd.startsWith("delay")) {
          int microseconds = cmd.substring(6).toInt();
          Serial.print("Delaying for %d ms\n", microseconds);
          delay(microseconds);
        } else {
          Serial.printf("Invalid command: %s\n", cmd);
        }

        cmdtok = strtok(nullptr, delim);
      } 

}

void setup() {
  pinMode(STATUS_PIN, OUTPUT);
  pinMode(LED_BUILTIN, OUTPUT);

  // begin serial communications
  Serial.begin(9600);
  while(!Serial) {}
  Serial.println("Serial connected");

  // connect to the network
  WiFi.begin(ssid, password);             
  Serial.printf("Connecting to \"%s\"...\n", ssid);

  while (WiFi.status() != WL_CONNECTED) {
    // wait for the Wi-Fi to connect
    delay(1000);
    Serial.print('-');
  }
  Serial.println("\nConnection established!");  
  
  Serial.printf("IP address:\t%s\n", WiFi.localIP());

  digitalWrite(STATUS_PIN, HIGH);

  delay(1000);

  // connect to websocket server
  Serial.println("Connecting to WebSocket server");
  socket.begin("192.168.0.2", 8080, "/");
  socket.onEvent(webSocketEvent);
  socket.setReconnectInterval(1000);
}

void loop() { 
  socket.loop();
}
