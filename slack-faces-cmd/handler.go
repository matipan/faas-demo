package function

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	whodis = "whodis"
	whodat = "whodat"
	teach  = "teach"
)

// Handle handles an incoming request from slack.
func Handle(w http.ResponseWriter, r *http.Request) {
	if token := r.FormValue("token"); token != slackToken() {
		http.Error(w, "Token is not valid", http.StatusForbidden)
		return
	}

	var msg string

	defer func() {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(msg))
	}()

	cmd := strings.Split(r.FormValue("text"), " ")
	if len(cmd) < 2 {
		msg = fmt.Sprintf("Invalid command: %s. Supported commands are: \n* `whodis [url]`\n* `teach [name] [url]`\n", r.FormValue("text"))
		return
	}
	action := cmd[0]

	var (
		form        = url.Values{}
		functionURL string
	)
	switch action {
	case whodis, whodat:
		url := cmd[1]
		if !isValidURL(url) {
			msg = fmt.Sprintf("%s: invalid URL %s", action, cmd[1])
			return
		}
		msg = fmt.Sprintf("%s: analyzing image at %s. Hang tight", action, url)
		form.Add("url", url)
		functionURL = os.Getenv("predict_function")
	case teach:
		if len(cmd) < 3 {
			msg = "teach: invalid use of command. Example: teach [name] [url]"
			return
		}
		url := cmd[len(cmd)-1]
		if !isValidURL(url) {
			msg = fmt.Sprintf("%s: invalid URL %s", action, cmd[1])
			return
		}
		msg = fmt.Sprintf("teach: teaching that %s is found at %s. Hang tight", cmd[1], url)
		form.Add("url", url)
		form.Add("name", strings.Join(cmd[1:len(cmd)-1], " "))
		functionURL = os.Getenv("teach_function")
	default:
		msg = fmt.Sprintf("Action %s is not supported", action)
		return
	}

	go func() {
		req, _ := http.NewRequest(http.MethodPost, functionURL, strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("X-Callback-Url", os.Getenv("results_callback"))
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Printf("Unable to post request to function %s. Err: %s", functionURL, err)
			return
		}
		if res.StatusCode >= 400 && res.StatusCode <= 600 {
			log.Printf("Function call responded with: %d", res.StatusCode)
		}
	}()
}

func isValidURL(uri string) bool {
	u, err := url.Parse(uri)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") {
		return false
	}
	return true
}

func slackToken() string {
	f, err := os.Open(os.Getenv("slack_token_secret"))
	if err != nil {
		log.Printf("Could not read secrets file: %s", err)
		return err.Error()
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return err.Error()
	}

	return strings.Trim(string(b), "\n")
}
