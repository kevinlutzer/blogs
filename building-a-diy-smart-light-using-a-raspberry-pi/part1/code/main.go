package main

import (
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
)

func main() {
	r := raspi.NewAdaptor()
	led := gpio.NewLedDriver(r, "10")

	robot := gobot.NewRobot("unused",
		[]gobot.Connection{r},
		[]gobot.Device{led},
	)

	c, err := newClient()
	if err != nil {
		panic(err)
	}

	fmt.Println("Setup Google IOT Core Config subscription")
	err = c.Subsribe(configTopic, func(_ MQTT.Client, m MQTT.Message) {
		if string(m.Payload()) == "ON" {
			led.On()
		} else {
			led.Off()
		}
	})

	if err != nil {
		panic(err)
	}

	robot.Start()
}
