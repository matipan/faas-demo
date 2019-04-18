package function

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type result struct {
	Text      string `json:"text"`
	IconEmoji string `json:"icon_emoji"`
}

// Handle handles the result of a particular action.
func Handle(w http.ResponseWriter, r *http.Request) {
	var (
		msg  string
		code = http.StatusOK
	)
	defer func() {
		w.WriteHeader(code)
		w.Write([]byte(msg))
	}()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Unable to read request body: %s", err)
		msg = "could not process command"
		code = http.StatusBadRequest
	}

	msg = string(b)
	b, _ = json.Marshal(result{Text: msg, IconEmoji: ":openfaas:"})
	_, err = http.Post(os.Getenv("slack_webhook"), "application/json", bytes.NewBuffer(b))
	if err != nil {
		log.Printf("Unable to post slack webhook: %s", err)
		msg = "could not deliver command results to slack"
		code = http.StatusInternalServerError
		return
	}

	msg = "successfuly delivered slack command"
}
