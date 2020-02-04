package functions

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	iotClient "cloud.google.com/go/iot/apiv1"
	"google.golang.org/api/option"
	"google.golang.org/genproto/googleapis/cloud/iot/v1"
)

// HelloWorld
func HelloWorld(w http.ResponseWriter, r *http.Request) {
	// b, err := ioutil.ReadAll(r.Body)
	// if err != nil {
	// 	fmt.Printf("Error: %s", err.Error())
	// }

	ctx := context.Background()
	client, err := iotClient.NewDeviceManagerClient(ctx, option.WithCredentialsFile("iot-klutzer-21813b0f64b4.json"))
	if err != nil {
		fmt.Printf("Error Setting Up Client: %s", err.Error())
		return
	}

	bs, _ := json.Marshal(struct {
		Key string `json:"key"`
	}{Key: "value"})

	if _, err := client.ModifyCloudToDeviceConfig(ctx, &iot.ModifyCloudToDeviceConfigRequest{
		Name:       "projects/repcore-prod/locations/us-central1/registries/devices-klutzer/devices/room-environment-monitor-personal",
		BinaryData: bs,
	}); err != nil {
		fmt.Printf("Error Modifying Config: %s", err.Error())
		return
	}

}
