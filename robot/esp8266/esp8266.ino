// WEB SOCKET AND WIFI INCLUDES
#include <ESP8266WiFi.h>
#include "wpa2_enterprise.h"
#include <WebSocketsClient.h>
#include "wifi.h"

WebSocketsClient webSocket;

void connectToEduroam() {
  const char ssid[32] = "eduroam";
  const char password[64] = "";
  wifi_station_disconnect();
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

  Serial.printf("Connecting to \"%s\"...\n", ssid);
  while (WiFi.status() != WL_CONNECTED) {
    // wait for the Wi-Fi to connect
    delay(1000);
    Serial.print('-');
  }

  Serial.println("\nConnection to shitroam established!");
  Serial.print("IP address:\t");
  Serial.println(WiFi.localIP());
}

// WEB SOCKET CONNNECTION
void webSocketEvent(WStype_t type, uint8_t * payload, size_t length) {
  switch (type) {
    case WStype_DISCONNECTED:
      Serial.printf("[WSc] Disconnected!\n");
      break;
    case WStype_CONNECTED:
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

void connectToWebsocket() {
  webSocket.begin("robot.test.halspals.co.uk", 80, "/connect/venus");
  webSocket.onEvent(webSocketEvent);
}

// INTERRUPT BUTTON //
#define LEFT_BUTTON D3
#define RIGHT_BUTTON D5

// BOLT LOCK/ LIGHT SWITCH BUTTONS//
bool lighton;
volatile long debounce_timeout = 0;

void leftButtonDown() {
  detachInterrupt(LEFT_BUTTON);
  attachInterrupt(LEFT_BUTTON, leftButtonUp, RISING);
  debounce_timeout = millis();
}

void leftButtonUp() {
  detachInterrupt(LEFT_BUTTON);
  attachInterrupt(LEFT_BUTTON, leftButtonDown, FALLING);
  if (debounce_timeout + 50 < millis()) {
    Serial.println("left");
    webSocket.sendTXT("left");
  }
}

void rightButtonDown() {
  detachInterrupt(RIGHT_BUTTON);
  attachInterrupt(RIGHT_BUTTON, rightButtonUp, RISING);
  debounce_timeout = millis();
}

void rightButtonUp() {
  detachInterrupt(RIGHT_BUTTON);
  attachInterrupt(RIGHT_BUTTON, rightButtonDown, FALLING);
  if (debounce_timeout + 50 < millis()) {
    Serial.println("right");
    webSocket.sendTXT("right");
  }
}

// THERMO //
int servo_angle = 0;
int maxServo = 180;
int minServo = 0;
int l = 0;
int r = 0;

// ~~~~~~ SETUP ~~~~~~ //
#include <Servo.h>
#define SERVO D0
#define SERVO_SHDN 15
Servo servo;

void executeCommand(String command) {
  Serial.print("Received command: ");
  Serial.print(command);
  Serial.println();
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
  // begin serial communications
  Serial.begin(115200);
  while (!Serial) {}
  Serial.println("Serial connected");
  pinMode(SERVO, OUTPUT);
  pinMode(SERVO_SHDN, OUTPUT);
  pinMode(LEFT_BUTTON, INPUT);
  pinMode(RIGHT_BUTTON, INPUT_PULLUP);
  servo.attach(SERVO);
  servo.write(0);
  SERVO_SHDN == LOW;

  // connect to the network and eduroam
  connectToEduroam();
  connectToWebsocket();

  // button interrupts
  attachInterrupt(digitalPinToInterrupt(LEFT_BUTTON), leftButtonDown, FALLING);
  attachInterrupt(digitalPinToInterrupt(RIGHT_BUTTON), rightButtonDown, FALLING);
}

void loop() {
  // //THERMO HARD CODED
  // r = digitalRead(RIGHT_BUTTON);
  // l = digitalRead(LEFT_BUTTON);

  // //left is pressed (turn on)
  // if (l == LOW ) {
  //   Serial.println("left is pressed");
  //   servo_angle = servo_angle + 7;
  //   if (servo_angle >= maxServo) {
  //     servo_angle = maxServo;
  //   }
  //   servo.write(servo_angle);
  //   delay(10);
  // }

  // //right is pressed (turn on)
  // if (r == LOW ) {
  //   Serial.println("right is pressed");
  //   servo_angle = servo_angle - 7;
  //   if (servo_angle <= minServo) {
  //     servo_angle = minServo;
  //   }

  //   servo.write(servo_angle);
  //   delay(10);
  // }

  //WEB SOCKET CODE
  webSocket.loop();
}