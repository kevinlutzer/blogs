package main

import MQTT "github.com/eclipse/paho.mqtt.golang"

type client struct {
	mqttClient MQTT.Client
}

func newClient() (*client, error) {
	// This method creates some default options for us, most notably it sets the auto reconnect option to be true, and the default port to `1883`. Auto reconnect is really useful in IOT applications as the internet connection may not always be extremely strong.
	opts := MQTT.NewClientOptions()

	// The specified The connection type we are using is just plain unencrypted TCP/IP
	opts.AddBroker("tcp://test.mosquitto.org:1883")
	// The client id needs to be unique, The argument passed was generated through a random number generator to avoid collisions.
	opts.SetClientID("F`/hty$3{+JQ%,j9")

	mqttClient := MQTT.NewClient(opts)

	// We have to create the connection to the broker manually and verify that there is no error.
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return &client{
		mqttClient,
	}, nil
}

// Publish publishes a message on a specific topic. An error is returned if there was problem. This function will publish with a QOS of 1.
func (c *client) Publish(msg, topic string) error {
	if token := c.mqttClient.Publish(topic, 1, false, msg); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}

// Subscribe creates a subscription for the passed topic. The callBack function is used to process any messages that the client recieves on that topic. The subscription created will have a QOS of 1.
func (c *client) Subsribe(topic string, f MQTT.MessageHandler) error {
	if token := c.mqttClient.Subscribe(topic, 0, f); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}
