# Building a DIY Smart Light Using A Raspberry Pi W/ Go - Part 1/2
IoT devices are ubiquitous in society these days. If you go into stores like Best Buy or Walmart, there are sections dedicated to internet-connected light switches, bulbs, thermostats, speakers, and doorbells. The construction of these devices is actually quite similar. Inside each one, there is some sort of processor that connects to a server in the cloud via a wifi module. This processor will pull download instructions like, "turn on light, " or "set the temperature to 21 degrees Celcius" from the server, which it will execute. On the other hand, the server will get its' instructions from different web/mobile applications like Smart CE or Google Smart Home, which would usually have been created by a human.
In this two-part series, I am going to show you how to make a smart light that you can control anywhere in the world with Google Assistant using a Raspberry Pi, [Google Cloud Platform](https://cloud.google.com/?hl=en), and [IFTTT](https://ifttt.com/). **All the code for this blog can be found [here](https://github.com/kevinlutzer/blogs/tree/master/building-a-diy-smart-light-using-a-raspberry-pi/part1/code)**.

**What is the Google Cloud Platform?** The Google Cloud Platform is a collection of products and tools that Google offers to build cloud based software. They have tools for webhosting, server management, database management, natural language processing, and IoT device management. Given that most of these tools are **free** to try, the platform is perfect for projects like this!

![Programming Paradigm](https://raw.githubusercontent.com/kevinlutzer/blogs/master/building-a-diy-smart-light-using-a-raspberry-pi/part1/images/tool-communication-paths.jpeg "Programming Paradigm")

The idea is, when you speak to a device using the Google Assistant saying our custom phrase, IFTTT will then call a Google Cloud Function web API. This API will then send device configuration to Google's IoT device management product, called IoT Core. Once the configuration is set, the Raspberry pi will pull that configuration from it and update the lights' state accordingly. The key phrase itself will contain the actual state the light will be set too.

**What are we building in this part?** We will be focusing on the interactions between the Google Cloud Platform and the Raspberry Pi. At the end of this blog you should have a light that you can turn on remotely using the Raspberry Pi, it just won't be controllable by Google Assistant.

## Setting Up the Raspberry Pi

For this project, you can use any Raspberry Pi. I will be using a model W. Next, you will have to set up the SD Card and connect the board to the internet. There are a ton of instructions on the internet for how to do this, but personally I always reference the [documentation](https://projects.raspberrypi.org/en/projects/raspberry-pi-setting-up/2) from the Raspberry Pi Foundation. **The main thing is that you need to be able to access the device over your local network.** We will be writing the code on our computers and use tools like Secure Copy (SCP) and Secure Shell (SSH) to copy and then run the code on the device.

**Prerequisite:** The instructions I will be giving assume that you are using a Mac or Linux based computer. If you are running a Windows machine, the terminal commands I will provide may not work for you.

## The Electric Circuit

For our light, let's use a simple led. Connect pin 10 (GPIO15) on the Raspberry pi to a 220 Ohm resistor, then connect the other side of the resistor to pin 9 (GND). You can see the schematic for the circuit in the image below.

![Schematic](https://raw.githubusercontent.com/kevinlutzer/blogs/master/building-a-diy-smart-light-using-a-raspberry-pi/part1/images/schematic.jpeg "Schematic")

To simply the construction, I soldered together a little contraption using a 2 pin female header, the resistor, and the LED. Here is an image of it on the Raspberry Pi:

![Raspberry Pi With Circuit](https://raw.githubusercontent.com/kevinlutzer/blogs/master/building-a-diy-smart-light-using-a-raspberry-pi/part1/images/raspberry_pi_with_circuit.jpeg "Raspberry Pi With Circuit")

## Testing the Circuit

At this point, you should have the circuit built, and be able to access the device remotely over your local network. Lets right some **Go code** that we can use to test that our LED is wired correctly.

**Why Go?** As Go has been rising in popularity over the last few years, frameworks and libraries like [Gobot.io](https://gobot.io) have been created to make building go applications with sensors and hardware controls a painless process. I strongly believe that for IoT devices running on a single board computer like a Raspberry Pi, Go is the best programming language. It is more performant then python has built-in concurrency concepts, and is extremely easy to use!

**Make sure that you have at least Go 1.13 installed your machine.** If you don't, you can find the install instructions [here](https://go.dev/doc/install). The first thing we need to do is to create a workspace. Run the following commands in a Linux or Mac OSX terminal.

``` bash
mkdir iot-device
cd iot-device
```

Next copy and paste the following code into a file called main.go.

``` go
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
```

All this code is doing, is toggling pin10 on for 5 seconds, then turning it off again. Now you should be able to run the following commands in a terminal to install the dependencies, build the code, copy it to the device, and execute it. Note that this requires you have both SCP and SSH installed on your machine, this should be the case for most Mac and Linux machines.

``` bash
# build go module and install dependencies
go mod init iot-client
go get gobot.io/x/gobot
go mod vendor# compile go code using 
GOARM=6 GOARCH=arm GOOS=linux go build ./... # copy the files to the raspberry pi and then execute
scp iot-client <USER>@<IP>:~ 
ssh -t <USER>@<IP> './iot-client' 
```

Please replace `<Raspberry Pi User>@<Raspberry Pi IP>` with your device’s actual IP and user name. This will look something like this `pi@192.168.1.111`.

If you don’t see the LED lighting up on your device, double-check the wiring. In your machine’s terminal, you should see the following output:

```
2020/03/10 00:21:26 Initializing connection RaspberryPi-35F5ED67 ...
2020/03/10 00:21:26 Initializing devices...
2020/03/10 00:21:26 Initializing device LED-4F9AEB15 ...
2020/03/10 00:21:26 Robot unused initialized.
2020/03/10 00:21:32 Starting Robot unused ...
2020/03/10 00:21:32 Starting connections...
2020/03/10 00:21:32 Starting connection RaspberryPi-35F5ED67...
2020/03/10 00:21:32 Starting devices...
2020/03/10 00:21:32 Starting device LED-4F9AEB15 on pin 10...
2020/03/10 00:21:32 Starting work...
```

## Using Google IoT Core

[IoT Core](https://cloud.google.com/iot-core) is a product Google offers to connect physical devices to the Google Cloud Platform (Gcloud). We can send configuration data to devices with the IoT Core Apis, and process device telemetry data with tools like [Pub/Sub](https://cloud.google.com/pubsub/docs) and [Cloud Functions](https://cloud.google.com/functions). I am going to show you how to set up Google IoT Core for a Google Cloud Project, Create IoT Core devices, and how to send configuration data to a device using IoT Core with MQTT.

**Quick Tip:** This blog builds on the concepts I described in a previous blog about MQTT, a publish-subscribe based message standard. We will be using a MQTT to communicate with IoT Core. If you are not familiar with this standard, check out my blog [here](https://medium.com/@kevinlutzer9/mqtt-with-go-fb5e97fbd394)

## Setting Up The Gcloud Project

A Google Cloud Project is just a workspace within the platform. We need to set up one to have access to IoT Core. If you don’t have a Google Cloud Project, click on this [link](https://console.cloud.google.com/projectcreate?pli=1) to create one. After you have created one, make sure you set up the Gcloud SDK on your machine. Installation and setup instructions can be found [here](https://cloud.google.com/sdk/docs/install).

Verify which project the Gcloud SDK on your machine is connected to by running the following command in a terminal instance:

``` bash
gcloud config get-value project
```

For me, this command will produce the following output:

``` bash 
Your active configuration is: [work]
iot-core-sample-klutzer
```

So “iot-core-sample-klutzer” is the project id I will use in my code, and with the CLI commands to set up the device with IoT Core.

## Enable the IoT Core

We need to enable our project to be able to access IoT Core. Without doing this step, all other steps **will not work**. Navigate to the [google cloud console](https://console.cloud.google.com/projectselector2/home/dashboard?supportedpurview=project) and enter the phrase “IoT Core” in the top search bar. You should be able to click on an autocompleted option for that product. Next when you get to the products page click “ENABLE”.

![Enable IoT Core](enable_iot_core.gif "Enable IoT Core")

## Receiving Configuration Data

Data sent to an IoT Core device is called configuration. There are two different options for sending it.

1. We can call an IoT Core API to update the device. Documentation of this API lives [here](https://cloud.google.com/iot-core). Google offers an SDK in multiple different languages to call this API. **In the next blog part, we will be using this functionality to allow Google Assistant to update device configuration.**

1. Using the Google Cloud Console UI, we can manually send data to a specific device. **For the purposes of this blog, this is what we are going to use.**

To receive the configuration, devices can subscribe to the “/devices/<DEVICE ID>/mqtt” topic.

## Setting up the IoT Device And Getting the Code

Now let's create an IoT Core Device Registry and Device. A Device Registry is just a mechanism to group similar devices. To create an IoT Device we only have to supply a unique identifier and a private key. The key allows us to have encryption when communicating with IoT Core as well as to authenticate the device.

We can use the Gcloud console to set up the device and registry configuration, but using the Gcloud SDK is easier. To create the public/private key pair used for encryption We are going to use OpenSSL to create the public/private key pair. Make sure you have that command installed. Here is the source page for [OpenSSL](https://www.openssl.org/). Below is a list of commands to run in a terminal that will create the registry, device, and private/public keys for the device. Replace any instances of “<PROJECT_ID>” with your actual project id, and make sure that this is the same project you enabled the IoT Core APIs for.

``` bash
# create cert files using open ssl
mkdir -p certs
openssl req -x509 -newkey rsa:2048 -keyout certs/rsa_private.pem -nodes -out certs/rsa_cert.pem
curl -o certs/roots.pem https://pki.google.com/roots.pem# create device registry and device
gcloud iot registries create devices --region=us-central1
gcloud iot devices create test-device --project=<PROJECT_ID> --region=us-central1 --registry=devices --public-key path=certs/rsa_cert.pem,type=rs256 # copy all cert files to the raspberry pi
scp -r certs <USER>@<IP>:certs
```

After the commands have completed, you should be able to see the device registry and the device on the IoT Core page. Navigate to the [IoT Core Registries page](https://console.cloud.google.com/iot/registries). Then select the device registry “devices”. Next, click on devices in the last navigation, and then select “test-device” from the listed devices. This page is for the **device details**. Make a note of the button that says ‘UPDATE CONFIG’, we will be using it later!

![Seeing Device Details](https://raw.githubusercontent.com/kevinlutzer/blogs/master/building-a-diy-smart-light-using-a-raspberry-pi/part1/images/seeing_device_details.gif "Seeing Device Details")

Let’s add the code required to connect to IoT Core and download device configuration. Start by adding a new file called mqtt_client.go. The code for this file is below:

``` go
package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"time"

	"io/ioutil"

	"github.com/dgrijalva/jwt-go"
	MQTT "github.com/eclipse/paho.mqtt.golang"
)

const (
	deviceID    = "test-device"
	registryID  = "devices"
	region      = "us-central1"
	configTopic = "/devices/test-device/config"
	certPath    = "certs/"
)

var projectID string

func init() {
	if projectID = os.Getenv("PROJECT_ID"); projectID == "" {
		panic("please specify the PROJECT_ID environement variable")
	}
}

func getSSLCerts() (rootsCert []byte, clientKey []byte, err error) {
	rootsCert, err = ioutil.ReadFile(certPath + "roots.pem")
	if err != nil {
		return
	}

	clientKey, err = ioutil.ReadFile(certPath + "rsa_private.pem")
	if err != nil {
		return
	}
	return
}

func getTokenString(rsaPrivate string) (string, error) {
	token := jwt.New(jwt.SigningMethodRS256)

	token.Claims = jwt.StandardClaims{
		Audience:  projectID,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(rsaPrivate))
	if err != nil {
		return "", err
	}

	tokenString, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return tokenString, err
}

func getTLSConfig(rootsCert string) *tls.Config {
	certpool := x509.NewCertPool()
	certpool.AppendCertsFromPEM([]byte(rootsCert))

	return &tls.Config{
		RootCAs:            certpool,
		ClientAuth:         tls.NoClientCert,
		ClientCAs:          nil,
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{},
		MinVersion:         tls.VersionTLS12,
	}
}

func newClient() (*client, error) {

	clientID := fmt.Sprintf("projects/%v/locations/%v/registries/%v/devices/%v",
		projectID,
		region,
		registryID,
		deviceID,
	)

	roots, clientKey, err := getSSLCerts()
	if err != nil {
		return nil, err
	}

	jwtString, err := getTokenString(string(clientKey))
	if err != nil {
		return nil, err
	}

	tlsConfig := getTLSConfig(string(roots))

	opts := MQTT.NewClientOptions()
	opts.AddBroker("ssl://mqtt.googleapis.com:8883")
	opts.SetClientID(clientID).SetTLSConfig(tlsConfig)
	opts.SetUsername("unused")
	opts.SetPassword(jwtString)

	mqttClient := MQTT.NewClient(opts)

	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		return nil, token.Error()
	}

	return &client{
		mqttClient,
	}, nil
}

type client struct {
	mqttClient MQTT.Client
}

func (c *client) Subsribe(topic string, f MQTT.MessageHandler) error {
	if token := c.mqttClient.Subscribe(topic, 0, f); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	return nil
}
```

This code wraps the Eclipse MQTT Client for Go. It contains functions to initialize the client/connection, and then subscribe to the device configuration.

Next, let's update the main.go file by replacing the contents with the following code snippet:

``` go
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
		} else if string(m.Payload()) == "OFF" {
			led.Off()
		}
	})

	if err != nil {
		panic(err)
	}

	robot.Start()
}
```

This new code creates an instance of the client struct defined in mqtt_client.go, and creates a subscription to the IoT Core config topic with a handler that updates the light state based on the values of “ON” and “OFF”. After you have updated files, your workspace should look have the following file structure:

```
+-- certs/
|   +-- roots.pem
|   +-- rsa_cert.pem
|   +-- rsa_private.pem
+-- go.mod
+-- go.sum
+-- main.go
+-- mqtt_client.go
+-- vendor/
|   +-- modules.txt
|   +-- ... etc ...
```
At this point, feel free to copy the code to your raspberry pi and execute it by running the following commands:
```
# Recompile go code, upload and execute
GOARM=6 GOARCH=arm GOOS=linux go build ./... #builds using Go ARM
scp iot-client <USER>@<IP>:~ #copy files
ssh -t <USER>@<IP> 'PROJECT_ID=<PROJECT_ID> ./iot-client' #execute
```

Make sure you replace “<PROJECT_ID>”, “<USER>”, and “<IP>” with the corresponding Google Cloud Project ID, Raspberry Pi User, and Raspberry Pi IP Address. You should be able to see some output about the subscription being setup, as well as, the application publishing data in your terminal. It will look similar to the following:

**Troubleshooting Tip:** If you see an error like “Unacceptable protocol version” that most likely means your device isn’t setup correctly. Make sure that the device id and the device registry id used in the code, are the same ones that were created in IoT Core.

## Let’s Try This Out!

Let’s now send some configuration using the Gcloud Console. Start by getting to the IoT Core device details page. We previously did this when verifying the device was created properly. Click on the button that is labeled, “UPDATE CONFIG”. This should pop open a window with a text area and a button labelled “SEND TO DEVICE”. Place the text “ON” in the text area and then click on the “SEND TO DEVICE” button. You should see the light turn on! If you don’t make sure that there are no spaces or other characters other then word “ON”. Note that character cases are important as well. Below is a GIF of me setting the config on the left, then on the right the LED Lighting up.

![Seeing Device Details](https://raw.githubusercontent.com/kevinlutzer/blogs/master/building-a-diy-smart-light-using-a-raspberry-pi/part1/images/seeing_device_details.gif "Seeing Device Details")

## What's Next

In the next blog, I will be demonstrating how to interface Google Assistant with IoT Core using Google Cloud Functions and IFTTT. Stay tuned!