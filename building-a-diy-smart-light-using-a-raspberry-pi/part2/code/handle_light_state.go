package handlelightstate

import (
	api "cloud.google.com/go/iot/apiv1"
	"context"
	"fmt"
	iotpb "google.golang.org/genproto/googleapis/cloud/iot/v1"
	"net/http"
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

	_, err = c.ModifyCloudToDeviceConfig(ctx, &iotpb.ModifyCloudToDeviceConfigRequest{
		Name:       "projects/iot-core-sample-klutzer/locations/us-central1/registries/devices/devices/test-device",
		BinaryData: []byte("ON"),
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
