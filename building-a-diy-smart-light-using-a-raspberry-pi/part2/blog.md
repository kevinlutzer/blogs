# Building a DIY Smart Light Using A Raspberry Pi W/ Go — Part 2/2

In [Part 1](https://medium.com/@kevinlutzer9/building-a-diy-smart-light-using-a-raspberry-pi-w-go-part-1-2-501efadcd36a) we connected a Raspberry Pi to Google Cloud Platform via IoT Core. In this part, I will demonstrate how to set the state of the light connected to the Rasberry Pi using Google Assistant.

**Prerequisite:** The instructions in this blog are going to be building on the steps described in [Part 1](https://medium.com/@kevinlutzer9/building-a-diy-smart-light-using-a-raspberry-pi-w-go-part-1-2-501efadcd36a). If you haven’t checked out that blog yet, I would recommend taking a look before reading this one!

To do this, we are going to use an IFTTT applet. It will allow us to send the desired light state to an API based on a custom command we say to Google Assistant. This API will then update the device configuration within Google IoT Core. Based on what we have currently built, the Raspberry Pi will pull down the updated configuration data and set the light’s state correctly.

## The API

We need an API that will allow us to update the IoT Core Configuration with the value of a passed parameter. This API will be publicly accessible so that it will be callable by IFTTT. Google provides a serverless technology called [Cloud Functions](https://cloud.google.com/functions/docs/concepts/overview) that will allow us to easily build this API.

To be consistent with the first part, we will be using Go 1.13 as the runtime environment. Let’s start by creating a workspace to write the code for our Cloud Function.

``` bash 
mkdir -p server
cd server
```

Next, create a file in the ‘server’ directory called ‘handle_light_state.go’. This file will contain the following snippet of code.

``` go
package handlelightstate

import (
	api "cloud.google.com/go/iot/apiv1"
	"context"
	"fmt"
	iotpb "google.golang.org/genproto/googleapis/cloud/iot/v1"
	"net/http"
	"strings"
)

const (
	projectID = "<PROJECT_ID>"
)

// HandleLightState state
func HandleLightState(w http.ResponseWriter, r *http.Request) {

	ctx := context.Background()
	c, err := api.NewDeviceManagerClient(ctx)
	if err != nil {
		msg := fmt.Sprintf("failed to setup device manager: %s", err.Error())
		fmt.Println(msg)
		w.Write([]byte(msg))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	q := r.URL.Query()
	state := strings.ToUpper(
		strings.Replace(
			strings.Replace(q.Get("state"), " ", "", -1),
			"\n", "", -1),
	)
	if !(state == "OFF" || state == "ON") {
		msg := fmt.Sprintf("the state query param must be specified and it must be \"OFF\" or \"ON\". specified value was \"%s\"", state)
		fmt.Println(msg)
		w.Write([]byte(msg))
		w.WriteHeader(http.StatusPreconditionFailed)
	}

	deviceName := fmt.Sprintf("projects/%s/locations/us-central1/registries/devices/devices/test-device", projectID)
	_, err = c.ModifyCloudToDeviceConfig(ctx, &iotpb.ModifyCloudToDeviceConfigRequest{
		Name:       deviceName,
		BinaryData: []byte(state),
	})

	if err != nil {
		msg := fmt.Sprintf("failed to update device configuration: %s", err.Error())
		fmt.Println(msg)
		w.Write([]byte(msg))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write([]byte("success"))
	return
}
```

**Make that you replace the constant ‘projectID’ with your actual project ID.** If you don’t, you will receive a response containing a forbidden (403) status code.

The function is looking at the query param ‘state’, validating that it is either ‘ON’ or ‘OFF’, and updating our device configuration with that value. Note that this is what we were manually doing at the end of the last blog using the IoT Core user interface.

Lets set up the dependencies for this function. Run the following commands in a terminal:

``` bash
go mod init server
go build ./... # This should pull in all of the dependencies
```

We are now ready to upload the code to the Google Cloud Platform. To do this, run the following command:

``` bash
gcloud functions deploy HandleLightState --entry-point HandleLightState --runtime go111 --trigger-http
```

This process **might** fail, and require you to run the following command:

``` bash
gcloud alpha functions add-iam-policy-binding HandleLightStateasd --member=allUsers --role=roles/cloudfunctions.invoker
```

All this does is make the function publically accessible. At the time of this writing, IFTTT has no authentication support. To make this a little bit more secure you could use an API key.

When the Cloud Function is done uploading the code, you should see some output about the configuration for the API. In that block of text, there should be a specified URL in the form of ‘https://us-central1-<PROJECT_ID>.cloudfunctions.net/HandleLightState’. That is the URL you have to call to hit the Cloud Function.

## Testing the API

Trying turning on your light by running the following command in your terminal.

``` bash
curl https://us-central1-<PROJECT_ID>.cloudfunctions.net/HandleLightState?state=ON
```

Make sure that the IoT client code is still running on the Raspberry Pi! Now if you want to turn the light off, just specify ‘OFF’ for the ‘state’ query param. Note the **state query param is not dependent on the case of the letters.** So ‘On’ will have the same effect as ‘ON’.

## Using IFTTT To Connect the API and Google Assistant

IFTTT stands for If This Then That. It is an API Platform that allows you to build simple apps by connecting functionality from 3rd party platforms together. IFTTT is based on the idea that you have some sort of trigger, that when activated will cause IFTTT to execute some sort of action. For us, the trigger is Google Assistant and the action is calling our API.

Start by navigating to https://ifttt.com/ and then signing up. After you create an account, navigate to the create page. From there, you should be able to click on the ‘+ IF’ button. Take a minute and look through the list of triggers that you can use.

![Choose Service](https://raw.githubusercontent.com/kevinlutzer/blogs/master/building-a-diy-smart-light-using-a-raspberry-pi/part2/images/choose_service.gif "Choose Service")

A ton of different companies have integrations with IFTTT, namely Google and Samsung. Use the search to find the Google Assistant. Click on that card, then click on the card with the text, “Say a phrase with a text ingredient.” You can now specify the parameters of the trigger. Most importantly you want to add the command we will use to set the state of the light, “Turn $ custom light.” You can specify ‘Done’ for the form field labeled, “What do you want the Assistant to say in response.” Note that the dollar sign will represent where our ‘On’ or ‘Off’ word will be.

![Setup App](https://raw.githubusercontent.com/kevinlutzer/blogs/master/building-a-diy-smart-light-using-a-raspberry-pi/part2/images/setup_app.gif "Setup App")

The next step is to set up the action. After you have finished creating the trigger, you should be back to the page that has the text, “If This Then That.” Click on the ‘+ That’ button and search for the Webhook action. Click on the card, that says “Make a webhook request”. You can specify the URL of the API we created before. Note that you don’t have to fill out any other information for this to work. The URL will be in the format of:

```bash
https://us-central1-<PROJECT_ID>.cloudfunctions.net/HandleLightState?state= {{TextField}}`
```
In this case, the `{{TextField}}` template variable will have the value of the ‘$’ we specified when creating the trigger. Once you click ‘Create Action’ and ‘Finish’, your IFTTT applet is complete!

![Setup Action](https://raw.githubusercontent.com/kevinlutzer/blogs/master/building-a-diy-smart-light-using-a-raspberry-pi/part2/images/setup_action.gif "Setup Action")

Through whatever device you have Google Assistant installed, say or type the command, “Turn on custom light.” As long as your Raspberry Pi is still running the client code you should see the light turn on. In real-time, this is what it looks for me.

[![Demo](https://img.youtube.com/vi/8xJSLQvz6fs/0.jpg)](https://www.youtube.com/watch?v=8xJSLQvz6fs)

