# Overview

MQTT is a lightweight network messaging standard based on the publish-subscribe design pattern. It is primarily used in IOT applications as there is built in safeguards to ensure message delivery under poor internet connections. Plain TCP, SSL/TLS sockets and websockets are all supported which gives us lots of flexable in the security of our application!

In this tutorial I will take you through creating a simple MQTT application that will publish a message and then subscribe to it. Lets get started!

## Architecture

In MQTT terms the `broker` is a server which recieves messages on a specific topic and publishes those messages to all other devices subscribing to that same topic. A topic is just a filter used to direct messages to the correct recepients. This allows for the clients in an MQTT application to send and recieve different messages without any collisions.

Lets take a look at the interactions between a server and the two clients during message transmission.

INsert image

So whats happening here? Client A tells the server that it is publishing a message on the topic `/mqtt-explanation-klutzer`. The server then sends the message to each subscriber of topic `/mqtt-explanation-klutzer` which in this case is just Client B. 

What if Client B looses internet access before it can recieve and process the message? How can we gaurentee that Client B recieves the message? That is where quality of service (QOS) comes in!

### Quality of Service (QOS)

Clients can define a QOS for both publish and subcribe actions. This setting describes the level of effort the broker will use to ensure message delivery. Here are the different QOS defintions from the [manual page](https://mosquitto.org/man/mqtt-7.html) for MQTT.

0: The broker/client will deliver the message once, with no confirmation.
1: The broker/client will deliver the message at least once, with confirmation required.
2: The broker/client will deliver the message exactly once by using a four step handshake.

As you increase in QOS, the reliability and rhobustness of the message transmissions increases, but more process power and internet bandwidth is required. Note that even though both the publisher and subscriber can set their own QOS, the lower number is the one that is honoured by the server. So if Client A sent the message with a QOS of 2, and Client B is subscribing with 0, 0 will be used.

### Message Retention

Lets go back to our previous example. What if we add Client C to our application? What if we want it to recieve the same message that Client A originally sent. That is where message rention comes in! When Client A publishes the message, it cam configure it to be retained by the server. Now when Client C subscribes to the topic `/mqtt-explanation-klutzer`, the server will send the message to it.

# Sample

Lets write a simple application that can subscribe to a specific topic and publish messages to it. So that means that we will only have one client interacting with the server. The code will be written in Golang, so make sure you have the latest version.

## Quick Start

If you want to run the sample code right away, download it from [this](https://raw.githubusercontent.com/kevinlutzer/blogs/master/iot-core-mqtt/sample) link. From the root of the directory, run `go run ./...` from your terminal. You should see the following text in your terminal

```
Message: Hello World 
Topic: mqtt/test/iot-mqtt-blog 
```

## The Code

The MQTT [client](https://github.com/eclipse/paho.mqtt.golang) that we are using is from the [Eclipse Foundation](https://www.eclipse.org/org/foundation/), I have found their project to be up to date and documented very well. The Eclipse Foundation provides an open access server operating as a MQTT broker for testing purposes. You can find the documentation for this server [here](`http://test.mosquitto.org/`).

Lets build a small wrapper around the MQTT Client to allow us to publish and subscribe to messages. We are going to use a struct to hold an instance of the client that has connections to the test MQTT broker. This struct will contain the different method-recievers we need to perform our publish and subscribe operations. Start by creating a directory called `sample` with an file `mqtt_client.go` containing the following code.

``` golang
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
```

Note that most of the methods that require some sort of interaction with the server like `Subsribe`, `Unsubscribe`, `Publish` and `Connect` return a token. This token has a `Wait` method which blocks until the operation has been completed. The `Error` method on the token will hold an error caused during the operations' execution.


Now that we have our client wrapper, lets built a small program to test it. In a new `main.go` file under the same directory add the following code.

``` golang
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

	// This Sleep gives enough time for us to recieve the message before we exit out of the application.
	time.Sleep(time.Second * 6)
}

```

I recommend making this workspace a go module. To do this, and install the depencies by running the following commands in your terminal.

``` bash
	go mod init sample
	go run ./... # This parse the .go files for the needed depencencies and build your go.sum and go.mod
	go mod vendor
```

You should now be able to execute this code by running `go run ./...` from the base directory, and see the following result

```
Message: Hello World 
Topic: mqtt/test/iot-mqtt-blog 
```