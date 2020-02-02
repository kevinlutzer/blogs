package main

import (
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"time"
)

func main() {

	topic := "mqtt/test/iot-mqtt-blog"

	c, _ := newClient()
	
	c.Subsribe(topic, func(_ MQTT.Client, m MQTT.Message) {
		fmt.Printf("Message: %s \n", m.Payload())
		fmt.Printf("Topic: %s \n", m.Topic())
	})
	
	c.Publish("Hello World", topic)
	
	time.Sleep(time.Second * 6)
}

type client struct {
	mqttClient MQTT.Client
}

func newClient() (*client, error) {
	opts := MQTT.NewClientOptions()
	opts.AddBroker("tcp://test.mosquitto.org:1883")
	opts.SetClientID("F`/hty$3{+JQ%,j9")

	mqttClient := MQTT.NewClient(opts)

	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return &client{
		mqttClient,
	}, nil
}

func (c *client) Cleanup(topics ...string) {
	c.mqttClient.Disconnect(250)
}

// Publish will sent Call this when we want to
func (c *client) Publish(msg, topic string) error {
	if token := c.mqttClient.Publish(topic, 1, false, msg); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

func (c *client) Subsribe(topic string, f MQTT.MessageHandler) error {
	if token := c.mqttClient.Subscribe(topic, 0, f); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}
