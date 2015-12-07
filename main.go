package main

import (
	"os"
	"log"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"regexp"
	"bytes"
)

const SIGNED_OFF_BY = "Signed-off-by:"

func main() {

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	http.HandleFunc("/", payloadHandler)
	err := http.ListenAndServe(":" + port, nil)
	if err != nil {
		panic(err)
	}
}

func payloadHandler(rw http.ResponseWriter, req *http.Request) {

	event := req.Header.Get("X-Github-Event")
	if req.Method == "GET" {
		succeed(rw, event)
		return
	}

	if event != "pull_request" {
		succeed(rw, event)
		log.Print("Request got non pull_request")
		return
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fail(rw, event, err)
		return
	}

	payload := PullEventRequest{}
	if err := json.Unmarshal(body, &payload); err != nil {
		fail(rw, event, err)
		return
	}

	commit_url := payload.PullRequest.CommitsUrl
	if !isSignOffPreset(commit_url) {
		log.Print("Some commits missing ", SIGNED_OFF_BY)
		setPoolRequestStatus(payload.PullRequest.StatusUrl, "failure", "Some commits missing " + SIGNED_OFF_BY )
	} else {
		log.Print("All commits have ", SIGNED_OFF_BY)
		setPoolRequestStatus(payload.PullRequest.StatusUrl, "success", "All commits contain " + SIGNED_OFF_BY)
	}
	succeed(rw, event)
}

func setPoolRequestStatus(status_url string, status string, description string) {

	access_token := os.Getenv("GITHUB_TOKEN")
	hgStatus := Status{
		State: status,
		Context: "Signed-off-by validator",
		Description:  description,
	}

	msg, err := json.Marshal(hgStatus)
	if err != nil {
		log.Fatal("Failed to encode status update")
		return
	}

	req, err := http.NewRequest("POST", status_url, bytes.NewBuffer(msg))
	if err != nil {
		log.Fatal("can't create a new request")
		return
	}
	client := http.Client{}
	req.Header.Add("Authorization", "token " + access_token)
	res, err := client.Do(req)
	if err != nil {
		log.Fatal("Failed to update pull request:", err.Error())
		return
	}

	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("failed to update pull request: ", err.Error())
	}
}

func isSignOffPreset(commit_url string) bool {

	l, err := http.Get(commit_url)
	if err != nil {
		return false
	}
	body, err := ioutil.ReadAll(l.Body)
	commits := []PullRequestCommit{}
	json.Unmarshal(body, &commits)
	matcher := regexp.MustCompile(SIGNED_OFF_BY)
	for i := range commits {
		if !matcher.MatchString(commits[i].Commit.Message) {
			return false
		}
	}
	return true
}

func succeed(w http.ResponseWriter, event string) {
	render(w, PayloadPong{
		Ok:    true,
		Event: event,
	})
}

func fail(w http.ResponseWriter, event string, err error) {
	w.WriteHeader(500)
	render(w, PayloadPong{
		Ok:    false,
		Event: event,
		Error: err.Error(),
	})
}

func render(w http.ResponseWriter, v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	w.Write(data)
}