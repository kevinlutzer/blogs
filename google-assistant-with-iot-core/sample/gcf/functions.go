package functions

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	iotClient "cloud.google.com/go/iot/apiv1"
	"google.golang.org/genproto/googleapis/cloud/iot/v1"
)

// HelloWorld
func HelloWorld(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	client, err := iotClient.NewDeviceManagerClient(ctx)
	if err != nil {
		msg := fmt.Sprintf("Error When Setting Up Client: %s", err.Error())
		fmt.Println(msg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	bs, _ := json.Marshal(struct {
		Key string `json:"key"`
	}{Key: "value"})

	if _, err := client.ModifyCloudToDeviceConfig(ctx, &iot.ModifyCloudToDeviceConfigRequest{
		Name:       "projects/iot-klutzer/locations/us-central1/registries/devices-klutzer/devices/room-environment-monitor-personal",
		BinaryData: bs,
	}); err != nil {
		msg := fmt.Sprintf("Error When Modifying Device: %s", err.Error())
		fmt.Println(msg)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

}
