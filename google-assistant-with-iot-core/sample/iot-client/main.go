package main

import (
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"time"
)

func main() {

	certs, err := getSSLCerts()
	if err != nil {
		panic(err)
	}

	c, _ := newClient(certs)

	fmt.Println("Setup subscription")
	if err := c.Subsribe(func(_ MQTT.Client, m MQTT.Message) {
		fmt.Printf("Message: %s \n", m.Payload())
		fmt.Printf("Topic: %s \n", m.Topic())
	}); err != nil {
		panic(err)
	}

	time.Sleep(time.Minute * 100)
}
