package main

import (
	"encoding/json"
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
	if err := c.Subsribe(configTopic, func(_ MQTT.Client, m MQTT.Message) {
		if len(m.Payload()) == 0 {
			return
		}
		fmt.Printf("Recieved configuration message: %s \n", m.Payload())
	}); err != nil {
		panic(err)
	}

	errors := make(chan error)
	go func() {
		for true {
			fmt.Println("Publishing message")
			payload := struct {
				Value     string    `json:"value"`
				Timestamp time.Time `json:"timestamp"`
			}{
				Value:     "Foo Bar",
				Timestamp: time.Now().UTC(),
			}
			b, _ := json.Marshal(payload)
			if err := c.Publish(telemetryTopic, b); err != nil {
				errors <- err
			}
			time.Sleep(time.Second * 5)
		}
	}()

	if err := <-errors; err != nil {
		fmt.Errorf("Error with publishing message: %s", err.Error())
	}

	time.Sleep(time.Minute * 100)
}
