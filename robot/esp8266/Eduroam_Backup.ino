#include <ESP8266WiFi.h>
extern "C" {
#include "user_interface.h"
#include "wpa2_enterprise.h"
}
#define SSID          "eduroam"
#define PASSWORD      ""
#define WPA2_USERNAME "Your username"
#define WPA2_IDENTITY WPA2_USERNAME
#define WPA2_PASSWORD "Your Password"
const char* host = "www.google.co.uk";
int counter = 0;

void setup() {
  // put your setup code here, to run once:
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
    typedef enum {
      EAP_TLS,
      EAP_PEAP,
      EAP_TTLS,
    } eap_method_t;

    eap_method_t method = EAP_TTLS;
    const char *identity = WPA2_IDENTITY;
    const char *username = WPA2_USERNAME;
    const char *password = WPA2_PASSWORD;

    wifi_station_set_wpa2_enterprise_auth(1);

    wifi_station_set_enterprise_identity((u8*)(void*)identity, os_strlen(identity));

    if (method == EAP_TLS) {
      Serial.println("error");
      //wifi_station_set_enterprise_cert_key(client_cert, os_strlen(client_cert) + 1, client_key, os_strlen(client_key) + 1, NULL, 1);
      //wifi_station_set_enterprise_username(username, os_strlen(username));//This is an option for EAP_PEAP and EAP_TLS.
    }
    else if (method == EAP_PEAP || method == EAP_TTLS) {
      wifi_station_set_enterprise_username((u8*)(void*)username, os_strlen(username));
      wifi_station_set_enterprise_password((u8*)(void*)password, os_strlen(password));
      //wifi_station_set_enterprise_ca_cert(ca, os_strlen(ca)+1);//This is an option for EAP_PEAP and EAP_TTLS.
    }
  }
  wifi_station_connect();

  // Wait for connection AND IP address from DHCP
  
  while (counter <= 10)
  {
    Serial.print("Status: ");
    Serial.print(wifi_station_get_connect_status());
    
    Serial.print(" - Arduino status: ");
    Serial.print(WiFi.status());
    Serial.print(" - Local IP:");
    Serial.println(WiFi.localIP());
    delay(2000);
    counter ++ ;
  }
} // setup


void loop() {
  // put your main code here, to run repeatedly:
  Serial.print("Connecting to website: ");
  Serial.println(host);
  WiFiClient client;
  if (client.connect(host, 80)) {
    String url = "/rele/rele1.txt";
    Serial.println(String("GET ") + url + " HTTP/1.1\r\n" + "Host: " + host + "\r\n" + "User-Agent: ESP32\r\n" + "Connection: close\r\n\r\n");
    while (client.connected()) {
      String line = client.readStringUntil('\n');
      if (line == "\r") {
        break;
      }
    }
    String line = client.readStringUntil('\n');
   Serial.println(line);
  }else{
      Serial.println("Connection unsucessful");
    }  
}
