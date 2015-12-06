package main

import (
//	"os"
//	"net/http"
	"os"
	"fmt"
	"log"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"github.com/GitbookIO/go-github-webhook"
)

const user_agent = "kofemann-go-agent/github-webhook"

func main() {

	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	http.HandleFunc("/", payloadHandler)
	fmt.Println("Hello ", port)
	err := http.ListenAndServe(":" + port, nil)
	if err != nil {
		panic(err)
	}
}


type Repository struct {
	Name string `json:"name"`
	FullName string `json:"full_name"`
}

type PullRequest struct {
	CommitsUrl string `json:"commits_url"`
}

type PullEventRequest struct {
 	Action		string	`json:"action"`
	Number 		int 	`json:"number"`
 	PullRequest PullRequest `json:"pull_request"`
	Repository	Repository `json:"repository"`
}

func payloadHandler(rw http.ResponseWriter, req *http.Request) {

	if req.Method == "GET" {
		succeed(rw, "Hello World")
		return
	}

	if req.Header.Get("X-Github-Event") != "pull_request" {
		succeed(rw, "Nothing to do (not a pull request)")
		log.Print("Request got non pull_request")
		return
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fail(rw, "Cant read", err)
		return
	}

	payload := PullEventRequest{}
	if err := json.Unmarshal(body, &payload); err != nil {
		fail(rw, "Could not deserialize payload", err)
		return
	}

	fmt.Println("Received", payload.Action, "for", payload.Repository.FullName)
	commit_url := payload.PullRequest.CommitsUrl
	get_change_list(commit_url)
	succeed(rw, "All checks pass")
}


func get_change_list(commit_url string) {
	fmt.Println("checking url:", commit_url)
	l, err := http.Get(commit_url)
	if err != nil {
		fmt.Println("Failed to get list of commits", err.Error())
		return
	}
	fmt.Print("%s", l)
}

func succeed(w http.ResponseWriter, event string) {
	render(w, github.PayloadPong{
		Ok:    true,
		Event: event,
	})
}

func fail(w http.ResponseWriter, event string, err error) {
	w.WriteHeader(500)
	render(w, github.PayloadPong{
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
