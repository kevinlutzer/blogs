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

	c, err := newClient(certs)
	if err != nil {
		panic(err)
	}

	fmt.Println("Setup Google IOT Core Config subscription")
	if err := c.Subsribe(func(_ MQTT.Client, m MQTT.Message) {
		if len(m.Payload()) == 0 {
			return
		}
		fmt.Printf("Message: %s \n", m.Payload())
	}); err != nil {
		panic(err)
	}

	time.Sleep(time.Minute * 100)
}
