#include <Servo.h>

#include <WebSocketsClient.h>
#include <ESP8266WiFi.h>
#include "ESP8266HTTPClient.h"
#include "wifi.h"

#include "user_interface.h"
#include "wpa2_enterprise.h"

#define SERVO_PIN D6
#define STATUS_PIN D7
#define CONTROL_LED 3
  
WebSocketsClient socket;
WiFiClient client;
HTTPClient http;

Servo servo;

void connectToEduroam(){
  wifi_station_disconnect();

  const char ssid[32] = "eduroam";
  const char password[64] = EDUROAM_PASSWORD;
  
  struct station_config sta_conf;
  os_memset(&sta_conf, 0, sizeof(sta_conf));
  os_memcpy(sta_conf.ssid, ssid, 32);
  os_memcpy(sta_conf.password, password, 64);
  wifi_station_set_config(&sta_conf);

  wifi_station_set_wpa2_enterprise_auth(1);
  wifi_station_set_enterprise_identity((u8*)(void*)EDUROAM_USERNAME, os_strlen(EDUROAM_USERNAME));
  wifi_station_set_enterprise_username((u8*)(void*)EDUROAM_USERNAME, os_strlen(EDUROAM_USERNAME));
  wifi_station_set_enterprise_password((u8*)(void*)EDUROAM_PASSWORD, os_strlen(EDUROAM_PASSWORD));

  wifi_station_connect();
}

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
        
        int first = cmd.indexOf('(');
        int last = cmd.indexOf(')');
        String id = cmd.substring (first+1, last-first);

        Serial.printf("Command ID: %s\n", id.c_str());

        cmd = cmd.substring(last-first+2, cmd.length() - first);

        if(cmd == "led on") {
          Serial.println("Turning LED on");
          digitalWrite(CONTROL_LED, LOW);
        } else if (cmd == "led off") {
          Serial.println("Turning LED off");
          digitalWrite(CONTROL_LED, HIGH);
        } else if (cmd.startsWith("servo")) {
          // convert characters after "servo " to int
          int angle = cmd.substring(6).toInt();
  
          Serial.printf("Setting servo to %d degrees\n", angle);

          // write to the servo so it moves to this angle when it is attached
          Serial.println("s w 1");
          servo.write(angle);
          
          // attach the servo
          Serial.println("s a 1");
          servo.attach(SERVO_PIN);
          delay(1000); // wait for movement

          // detach the servo
          Serial.println("s d 1");
          servo.detach();
          delay(500);
        } else if (cmd.startsWith("delay")) {
          int microseconds = cmd.substring(6).toInt();
          Serial.printf("Delaying for %d ms\n", microseconds);
          delay(microseconds);
        } else {
          Serial.printf("Invalid command: %s\n", cmd.c_str());
        }

        // acknowledge command
        http.begin("http://robot.test.halspals.co.uk/123abc/acknowledge/" + id);
        int httpCode = http.POST("");
        if (httpCode != 200) { 
          Serial.printf("http error acknowledging command, code: %d\n", httpCode);
        }
        http.end();

        cmdtok = strtok(nullptr, delim);
      } 

}

void setup() {
  pinMode(STATUS_PIN, OUTPUT);
  pinMode(CONTROL_LED, OUTPUT);
  pinMode(SERVO_PIN, OUTPUT);

  // begin serial communications
  Serial.begin(9600);
  while(!Serial) {}
  Serial.println("Serial connected");

  // connect to the network
  // WiFi.begin(ssid, password);
  connectToEduroam();             
//  Serial.printf("Connecting to \"%s\"...\n", ssid);

  while (WiFi.status() != WL_CONNECTED) {
    // wait for the Wi-Fi to connect
    delay(1000);
    Serial.print('-');
  }
  Serial.println("\nConnection established!");  
  
//  Serial.printf("IP address:\t%s\n", WiFi.localIP());

  digitalWrite(STATUS_PIN, HIGH);

  delay(1000);

  // connect to websocket server
//  Serial.println("Connecting to WebSocket server");
//  socket.begin("192.168.0.2", 8080, "/");
//  socket.onEvent(webSocketEvent);
//  socket.setReconnectInterval(1000);
}

void loop() { 

  
  http.begin("http://robot.test.halspals.co.uk/123abc/commands");
  int httpCode = http.GET();

  if (httpCode == 200) { //Check the returning code
    String payload = http.getString();
    Serial.printf("http payload: %s\n", payload.c_str());

    payload.trim();
    executeCommandSequence(payload);
    
  } else {
    Serial.printf("http error, code: %d\n", httpCode);
  }

  http.end();
  
  String command = Serial.readString();

  delay(2000);
}
