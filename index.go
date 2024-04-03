package imgShareAPI

import (
	"encoding/json"
	"fmt"
	"html"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
	functions.HTTP("img-share-api-func", imgShareAPIFunc)
}

func imgShareAPIFunc(w http.ResponseWriter, r *http.Request) {
	var d struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		fmt.Fprint(w, "Hello, World! Success")
		return
	}
	if d.Name == "" {
		fmt.Fprint(w, "Hello, World! Success 2")
		return
	}
	fmt.Fprintf(w, "Hello, %s! Success 3", html.EscapeString(d.Name))
}
