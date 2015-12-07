package main


type Commit struct {
	Message string `json"message"`
}

type PullRequestCommit struct {
	Sha string `json:"sha"`
	Commit Commit `json"commit"`
}

type Repository struct {
	Name string `json:"name"`
	FullName string `json:"full_name"`
}

type PullRequest struct {
	CommitsUrl string `json:"commits_url"`
	StatusUrl   string `json:"statuses_url"`
}

type PullEventRequest struct {
	Action		string	`json:"action"`
	Number 		int 	`json:"number"`
	PullRequest PullRequest `json:"pull_request"`
	Repository	Repository `json:"repository"`
}

type PayloadPong struct {
	Ok    bool   `json:"ok"`
	Event string `json:"event"`
	Error string `json:"error,omitempty"`
}

type Status struct {
	State string `json:"state"`
	Description string `json:"description"`
	Context string `json:"context"`
	TargetUrl string `json:"target_url"`
}
