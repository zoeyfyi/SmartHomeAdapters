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
      String command = Serial.readString();
      command.remove(command.length() - 1); // remove new line

      Serial.print("Received command: ");
      Serial.print(command);
      Serial.println();

      // process LED commands
      if (command.startsWith("led")) {
        if(command == "led on") {
            Serial.println("Turning LED on");
            digitalWrite(LED_BUILTIN, HIGH);
        } else if (command == "led off") {
            Serial.println("Turning LED off");
            digitalWrite(LED_BUILTIN, LOW);
        }
      }
  }
}
