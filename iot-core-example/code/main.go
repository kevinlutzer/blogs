package main

import (
	"time"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
)

func main() {

	r := raspi.NewAdaptor()
	led := gpio.NewLedDriver(r, "15")

	work := func() {
		gobot.Every(1*time.Second, func() {
			led.Toggle()
		})
	}

	robot := gobot.NewRobot("blinkBot",
		[]gobot.Connection{r},
		[]gobot.Device{led},
		work,
	)

	robot.Start()

	// c, err := newClient()
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println("Setup Google IOT Core Config subscription")
	// if err := c.Subsribe(configTopic, func(_ MQTT.Client, m MQTT.Message) {
	// 	if len(m.Payload()) == 0 {
	// 		return
	// 	}
	// 	fmt.Printf("Recieved configuration message: %s \n", m.Payload())
	// }); err != nil {
	// 	panic(err)
	// }

	// // Block indefinitely
	// b := make(chan struct{})
	// <-b
}
