package functions

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	iotClient "cloud.google.com/go/iot/apiv1"
	"google.golang.org/genproto/googleapis/cloud/iot/v1"
)

// UpdateIOTDeviceConfig updates the device config for a iot device
func UpdateIOTDeviceConfig(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		msg := fmt.Sprintf("Error When reading body: %s", err.Error())
		fmt.Println(msg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	ctx := context.Background()
	client, err := iotClient.NewDeviceManagerClient(ctx)
	if err != nil {
		msg := fmt.Sprintf("Error When Setting Up Client: %s", err.Error())
		fmt.Println(msg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if _, err := client.ModifyCloudToDeviceConfig(ctx, &iot.ModifyCloudToDeviceConfigRequest{
		Name:       "projects/iot-klutzer/locations/us-central1/registries/devices-klutzer/devices/test-device",
		BinaryData: b,
	}); err != nil {
		msg := fmt.Sprintf("Error When Modifying Device: %s", err.Error())
		fmt.Println(msg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

}
