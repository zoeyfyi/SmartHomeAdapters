#include <ESP8266WiFi.h>        // Include the Wi-Fi library
#include "wifi.h"

void setup() {
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
}

void loop() { 
}
