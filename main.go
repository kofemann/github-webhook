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

func payloadHandler(rw http.ResponseWriter, req *http.Request) {

	if req.Method == "GET" {
		succeed(rw, "Hello World")
		return
	}

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fail(rw, "Cant read", err)
		return
	}

	payload := github.GitHubPayload{}
	if err := json.Unmarshal(body, &payload); err != nil {
		fail(rw, "Could not deserialize payload", err)
		return
	}

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
