#include <OneWire.h>
#include <DallasTemperature.h>
#include <PubSubClient.h>
#include <ESP8266WiFi.h>


// WIFI
const char* ssid = "TODO";
const char* password = "TODO";
WiFiClient espWifi;

// DS
#define ONE_WIRE_BUS D4
OneWire oneWire(ONE_WIRE_BUS);
DallasTemperature sensors(&oneWire);


// emqx server
const char *emqx_host = "TODO";
const char *topic = "TODO";
const char *emqx_username = "emqx";
const char *emqx_password = "emqx";
const int emqx_port = 1883;
PubSubClient client(espWifi);
struct msg {
  int distance;
  int light;
  float temperature;
};

msg data;
    

// LDR
#define LDRpin A0
int LDRValue = 0;

// ultra sonic
#define trigPin D8
#define echoPin D7
long duration;

void setup(void) {
  // set ultrasonic pins mode
  pinMode(trigPin, OUTPUT);
  pinMode(echoPin, INPUT);

  // begin serial port
  Serial.begin(115200);
  delay(100);


  // connect to wifi
  WiFi.begin(ssid, password);
  while (WiFi.status() != WL_CONNECTED) {
    delay(1000);
    Serial.println("Haven't connected yet");
  }
  Serial.println("Connected!");
  Serial.print("IP: ");
  Serial.println(WiFi.localIP());
  

  // set emqx server conf
  client.setServer(emqx_host, emqx_port);

  // connect to emqx
  while (!client.connected()) {
      String client_id = "esp826";
      Serial.println("conecting to emqx server");
      if (client.connect(client_id.c_str(), emqx_username, emqx_password)) {
          Serial.println("connected to emqx server");
      } else {
          Serial.print("failed with state ");
          Serial.print(client.state());
          delay(2000);
      }
  }

  // begin ds sensor
  sensors.begin();

}

void loop(void){
  // clear garbage on trig pin
  digitalWrite(trigPin, LOW);
  delayMicroseconds(2);
  // sets the trigPin on HIGH state for 10 micro seconds
  digitalWrite(trigPin, HIGH);
  delayMicroseconds(10);
  digitalWrite(trigPin, LOW);
  // reads the echoPin, returns the sound wave travel time in microseconds
  duration = pulseIn(echoPin, HIGH);

  // Calculating the distance
  data.distance = duration * 0.034 / 2;
  Serial.print("Distance: ");
  Serial.println(data.distance);

  
  sensors.requestTemperatures(); 
  // since we have only one ds sensor, read first index in sensors
  data.temperature = sensors.getTempCByIndex(0);
  Serial.print("Celsius temperature: ");
  Serial.println(data.temperature);


  data.light = analogRead(LDRpin);
  Serial.print("LDR: ");
  Serial.println(data.light);
  char dataTemp[32];
  sprintf(dataTemp, "%d,%d,%2f", data.light, data.distance, data.temperature);
  client.publish(topic, dataTemp);
  delay(6000);
}
