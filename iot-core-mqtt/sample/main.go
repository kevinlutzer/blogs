package main

import (
	"fmt"
	"log"
	"os"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

var (
	logger *log.Logger
)

const (
	brokerURI  = "test.mosquitto.org"
	brokerPort = "1883"
	topic      = "mqtt/test/klutzer"
)

func init() {
	logger = log.New(os.Stdout, "", log.Ldate|log.Ltime)
}

func main() {
	resp := Subsribe()
	fmt.Println(resp)
	Publish("Hello World")
}

func getMQTTClient() MQTT.Client {

	clientID := "F`/hty$3{+JQ%,j9"
	broker := fmt.Sprintf("tcp://%v:%v", brokerURI, brokerPort)

	opts := MQTT.NewClientOptions()
	opts.AddBroker(broker)
	opts.SetClientID(clientID)

	return MQTT.NewClient(opts)
}

func Publish(msg string) {
	client := getMQTTClient()
	token := client.Publish(topic, 0, false, msg)
	token.WaitTimeout(100)
	err := token.Error()
	if err != nil {
		log.Fatalf("failed to publish the payload: %v\n", err.Error())
	}
	client.Disconnect(50)
}

func Subsribe() string {
	choke := make(chan [2]string)
	client := getMQTTClient()

	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatalf("Failed to subscribe with error %s", token.Error().Error())
	}

	token := client.Subscribe(
		topic,
		0,
		func(client MQTT.Client, msg MQTT.Message) {
			choke <- [2]string{msg.Topic(), string(msg.Payload())}
		},
	)

	token.WaitTimeout(100)
	err := token.Error()
	if err != nil {
		logger.Fatalln("failed to publish the payload")
	}

	incoming := <-choke

	client.Disconnect(250)
	return incoming[1]
}
