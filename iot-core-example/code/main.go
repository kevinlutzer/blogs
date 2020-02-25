package main

import (
	"fmt"

	MQTT "github.com/eclipse/paho.mqtt.golang"
)

func main() {
	c, err := newClient()
	if err != nil {
		panic(err)
	}

	fmt.Println("Setup Google IOT Core Config subscription")
	if err := c.Subsribe(configTopic, func(_ MQTT.Client, m MQTT.Message) {
		if len(m.Payload()) == 0 {
			return
		}
		fmt.Printf("Recieved configuration message: %s \n", m.Payload())
	}); err != nil {
		panic(err)
	}

	// Block indefinitely
	b := make(chan struct{})
	<-b
}
