package function

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/machinebox/sdk-go/facebox"
)

// Handle receives a URL and trys to recognize who
// is at the picture.
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

	log.Printf("Finding similar faces at %s", u)
	client := facebox.New(os.Getenv("machinebox_url"))
	faces, err := client.SimilarsURL(u, 5)
	if err != nil {
		log.Printf("Unable to find similar faces: %s", err)
		msg = fmt.Sprintf("could not find similar faces")
		code = http.StatusInternalServerError
		return
	}

	for _, face := range faces {
		msg = fmt.Sprintf("Results from %s\n", u.String())
		for _, f := range face.SimilarFaces {
			msg += fmt.Sprintf("%d%% confident that %s is in the picture\n", int(f.Confidence*100), f.Name)
			return
		}
	}
	msg = "there were no known faces in the picture"
}
