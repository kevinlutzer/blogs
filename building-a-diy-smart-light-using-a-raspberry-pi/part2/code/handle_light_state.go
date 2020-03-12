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
