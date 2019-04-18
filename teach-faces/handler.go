package function

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/google/uuid"
	"github.com/machinebox/sdk-go/facebox"
)

// Handle receives a url to a face and a name
// and sends the name and url to machine box.
func Handle(w http.ResponseWriter, r *http.Request) {
	var (
		msg  string
		code = http.StatusOK
	)
	defer func() {
		w.WriteHeader(code)
		w.Write([]byte(msg))
	}()

	u, err := url.Parse(r.FormValue("url"))
	if err != nil {
		log.Printf("Invalid url: %s. Err: %s", r.FormValue("url"), err)
		msg = fmt.Sprintf("url %s is invalid", r.FormValue("url"))
		code = http.StatusBadRequest
		return
	}

	name := r.FormValue("name")

	log.Printf("Teaching that %s can be found at %s", name, u)
	client := facebox.New(os.Getenv("machinebox_url"))
	if err := client.TeachURL(u, uuid.New().String(), name); err != nil {
		log.Printf("Unable to teach for url %s. Err: %s", u.String(), err)
		msg = "could not associate name with the picture"
		code = http.StatusInternalServerError
		return
	}
	msg = fmt.Sprintf("Learned who %s is from %s", name, u.String())
}
