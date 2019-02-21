#include <ESP8266WiFi.h>
#include <Servo.h>
#include <WebSocketsClient.h>

#include "user_interface.h"
#include "wpa2_enterprise.h"

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

 // connect to eduroam
   wifi_station_disconnect();

   Serial.begin(115200);
   Serial.setDebugOutput(true);

   Serial.print("Trying to connect to ");
   Serial.println(SSID);

     {
       char ssid[32] = SSID;
       char password[64] = PASSWORD;
       struct station_config sta_conf;// = { 0 };

       os_memset(&sta_conf, 0, sizeof(sta_conf));
       os_memcpy(sta_conf.ssid, ssid, 32);
       os_memcpy(sta_conf.password, password, 64);
       wifi_station_set_config(&sta_conf);

     }

     {
       const char *identity = WPA2_IDENTITY;
       const char *username = WPA2_USERNAME;
       const char *password = WPA2_PASSWORD;

       wifi_station_set_wpa2_enterprise_auth(1);

       wifi_station_set_enterprise_identity((u8*)(void*)identity, os_strlen(identity));

       wifi_station_set_enterprise_username((u8*)(void*)username, os_strlen(username));
       wifi_station_set_enterprise_password((u8*)(void*)password, os_strlen(password));

     }

   wifi_station_connect();


 // Check if connected
   int i = 0;
   while (WiFi.status() != WL_CONNECTED) {
     // wait for the Wi-Fi to connect
     if (i <= 30) {
       digitalWrite(STATUS_PIN, HIGH);
       delay(500);
       digitalWrite(STATUS_PIN, LOW);
       delay(500);
       Serial.print('-');
       i++;
     } else {
       Serial.print("UNABLE TO CONNECT");
     }
   }

   while(1) {
     Serial.println('\n');
     Serial.println("Connection established!");
     Serial.print("IP address:\t");
     Serial.println(WiFi.localIP());
     digitalWrite(STATUS_PIN, HIGH);
     delay(1000);
   }

 // Connect to websocket server
 //  Serial.println("Connecting to WebSocket server");
 //  socket.begin("", 8080, "/connect");
 // socket.onEvent(webSocketEvent);
 //  socket.setReconnectInterval(1000);
}

void loop() {
//  socket.loop();
}
