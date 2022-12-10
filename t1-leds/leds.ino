// LEDs
int ledRed = D0;
int ledGreen = D1;

// serial data variable
int led = 0;

void setup() {
  delay(1000);
  Serial.begin(115200);
  delay(500);

  // set led pins mode on output
  pinMode(ledRed, OUTPUT);
  pinMode(ledGreen, OUTPUT);
}

void loop() {
  if (Serial.available() > 0) {
    led = Serial.read();
    // 49 -> 1 , red
    // 50 -> 2 , green
    if (led == 50 ){
      digitalWrite(ledGreen, HIGH);
      delay(1000);
      digitalWrite(ledGreen, LOW);
    } else if (led == 49) {
      digitalWrite(ledRed, HIGH);
      delay(1000);
      digitalWrite(ledRed, LOW);
    }
  }
  delay(500);
}
