package functions

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// HelloWorld
func HelloWorld(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}
	fmt.Printf("Body: %s", b)
}
