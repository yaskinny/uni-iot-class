package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/castai/promwrite"
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type sensors struct {
	Ds  string
	Ldr string
	Us  string
}

var (
	sensorsData = sensors{
		Ds:  "0",
		Ldr: "0",
		Us:  "0",
	}
	// set kardane tanzimat baraye vasl shodan be mqtt server
	opts = mqtt.NewClientOptions().AddBroker("tcp://localhost:1883").SetClientID("golang")

	// in function zamani ke payami ruye topic daryaft mishe estefadeh mishe
	f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
		fmt.Printf("payam rooye topic %v daryaft shod\n", msg.Topic())
		// [0] -> ldr data
		// [1] -> ultrasonic data
		// [2] -> temperature
		d := strings.Split(string(msg.Payload()), ",")

		sensorsData.Ldr = d[0]
		sensorsData.Us = d[1]
		sensorsData.Ds = d[2]

		fmt.Printf("%v - ldr: %v , ultrasonic: %v , ds: %v \n", time.Now(), sensorsData.Ldr, sensorsData.Us, sensorsData.Ds)
		writeMetrics()
	}

	// mqtt topic
	topic = "sensors/#"

	// omqe taanker aab ro 564cm dar nazar migirim
	tankerHeight = 564
)

func main() {
	// tanzimante mortabet be mqtt connection
	opts.SetKeepAlive(60 * time.Second)
	opts.SetDefaultPublishHandler(f)
	opts.SetPingTimeout(1 * time.Second)

	// saakhtane yek kaarbar baraye mqtt server ba tanzimati ke az qabl tanzim kardim
	mqttClient := mqtt.NewClient(opts)

	log.Println("Dar haale vasl shodan be mqtt server...")

	// vasl shodan be mqtt server
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Panic(token.Error())
	}
	log.Println("Be mqtt server vasl shodim...")

	log.Printf("subscribe kardan be topic %v dar mqtt server...\n", topic)
	// subscribe kardan be topic e sensor ha dakhele mqtt server
	if token := mqttClient.Subscribe(topic, 1, nil); token.Wait() && token.Error() != nil {
		log.Panic((token.Error()))
	}
	log.Printf("be topic %v dar mqtt server subscribe kardim...\n", topic)

	// montazer baashim ta darkhaste baste shodan application ijad beshe ba feshordane ctrl+c
	exitApp(mqttClient)
}

// in function baraye neveshtane data darune database estefadeh mishavad
func writeMetrics() {
	// saakhtane yek kaarbar baraye vasl shodan be database
	client := promwrite.NewClient("http://localhost:9090/api/v1/write")

	// hesab kardane inke chand darsad az tanker khaali hast
	ius, _ := strconv.Atoi(sensorsData.Us)
	twp := (ius * 100) / tankerHeight

	fds, _ := strconv.ParseFloat(sensorsData.Ds, 64)
	fus := float64(twp)
	fldr, _ := strconv.ParseFloat(sensorsData.Ldr, 64)
	data := promwrite.WriteRequest{
		TimeSeries: []promwrite.TimeSeries{
			{
				Labels: []promwrite.Label{
					{
						Name:  "__name__",
						Value: "sensor_ds",
					},
				},
				Sample: promwrite.Sample{
					Time:  time.Now(),
					Value: fds,
				},
			},
			{
				Labels: []promwrite.Label{
					{
						Name:  "__name__",
						Value: "sensor_ldr",
					},
				},
				Sample: promwrite.Sample{
					Time:  time.Now(),
					Value: fldr,
				},
			},
			{
				Labels: []promwrite.Label{
					{
						Name:  "__name__",
						Value: "sensor_us",
					},
				},
				Sample: promwrite.Sample{
					Time:  time.Now(),
					Value: fus,
				},
			},
		},
	}
	client.Write(context.TODO(), &data)
}

// in function jahate khoruj az barnaameh estefadeh mishavad
func exitApp(c mqtt.Client) {
	s := make(chan os.Signal)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM)
	<-s
	log.Println("Shoma dar khaaste baste shodan application ra anjam daadid, bad az qat shodan application az mqtt server, application baste khaahad shod.")

	// unsubscribe kardan az topic sensor ha
	if token := c.Unsubscribe(topic); token.Wait() && token.Error() != nil {
		log.Panic((token.Error()))
	}

	// qat kardane connection az mqtt server
	c.Disconnect(500)

	// bastane barname
	os.Exit(0)
}
