package function

import (
	"fmt"
	"net/http"
)

// Handle handles an incoming request from slack.
func Handle(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Form: %+v\n", r.Form)
	fmt.Printf("Message: %s\n", r.FormValue("text"))
}
