package stringadderhandler

import (
	"context"
	"net/http"
	"google.golang.org/api/cloudiot/v1"
)
// StringAdderHandler is a handler that api-ifies the Add method in github.com/kevinlutzer/string-adder
func StringAdderHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	cloudiot.NewService(ctx)
}
