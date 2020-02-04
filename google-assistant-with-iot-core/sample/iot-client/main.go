package main

import (
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"time"
)

func main() {
	topic := "/mqtt-explanation-klutzer"

	c, _ := newClient()
	
	if err := c.Subsribe(topic, func(_ MQTT.Client, m MQTT.Message) {
		fmt.Printf("Message: %s \n", m.Payload())
		fmt.Printf("Topic: %s \n", m.Topic())
	}); err != nil {
		panic(err)
	}

	if err := c.Publish("Hello World", topic); err != nil {
		panic(err)
	}

	time.Sleep(time.Second * 6)
}
