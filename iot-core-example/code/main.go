package main

import (
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
	"time"
)

func main() {
	r := raspi.NewAdaptor()
	led := gpio.NewLedDriver(r, "10")

	robot := gobot.NewRobot("unused",
		[]gobot.Connection{r},
		[]gobot.Device{led},
	)

	led.On()
	time.Sleep(5 * time.Second)
	led.Off()

	robot.Start()
}
