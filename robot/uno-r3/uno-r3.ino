// the setup function runs once when you press reset or power the board
void setup() {
  // begin serial communications
  Serial.begin(9600);
  delay(10);
  
  // initialize digital pin LED_BUILTIN as an output.
  pinMode(LED_BUILTIN, OUTPUT);
  digitalWrite(LED_BUILTIN, LOW);
}

// the loop function runs over and over again forever
void loop() {
  if(Serial.available() > 0){
      // turn on when we receive a 1, and off when we receive a 0
      int value = Serial.read();
      if(value == '1') {
          digitalWrite(LED_BUILTIN, HIGH);
      } else if (value == '0') {
          digitalWrite(LED_BUILTIN, LOW);
      }
  }
}
