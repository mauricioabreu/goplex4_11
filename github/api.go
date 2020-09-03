package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// Issue contains issue data
type Issue struct {
	Number    int
	Title     string
	Body      string
	State     string
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CreateIssue create a new issue
func CreateIssue(owner, repo, title, body string) error {
	requestBody, err := json.Marshal(map[string]string{
		"title": title,
		"body":  body,
	})
	if err != nil {
		log.Fatal(err)
	}

	request, err := http.NewRequest("POST", buildIssueURL(owner, repo), bytes.NewBuffer(requestBody))
	request.Header.Set("Accept", "application/vnd.github.v3+json")
	request.Header.Set("Content-Type", "application/json")
	setAuthorization(request)
	if err != nil {
		log.Fatal(err)
	}

	timeout := time.Duration(5 * time.Second)
	client := http.Client{Timeout: timeout}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusCreated {
		return fmt.Errorf("failted to create issue: %s", response.Status)
	}
	return nil
}

func patchIssue(owner, repo, issueNumber string, values map[string]string) error {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.Encode(values)

	request, err := http.NewRequest("PATCH", buildIssuesURL(owner, repo, issueNumber), buffer)
	request.Header.Set("Accept", "application/vnd.github.v3+json")
	request.Header.Set("Content-Type", "application/json")
	setAuthorization(request)
	if err != nil {
		log.Fatal(err)
	}

	timeout := time.Duration(5 * time.Second)
	client := http.Client{Timeout: timeout}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to update issue %s: %s", issueNumber, response.Status)
	}
	return nil
}

// UpdateIssue update an issue
func UpdateIssue(owner, repo, title, body, issueNumber string) error {
	values := map[string]string{
		"title": title,
		"body":  body,
	}
	return patchIssue(owner, repo, issueNumber, values)
}

// GetIssue fetch issue information
func GetIssue(owner, repo, issueNumber string) (*Issue, error) {
	request, err := http.NewRequest("PATCH", buildIssuesURL(owner, repo, issueNumber), nil)
	request.Header.Set("Accept", "application/vnd.github.v3+json")
	request.Header.Set("Content-Type", "application/json")
	setAuthorization(request)
	if err != nil {
		log.Fatal(err)
	}

	timeout := time.Duration(5 * time.Second)
	client := http.Client{Timeout: timeout}
	response, err := client.Do(request)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get issue %s: %s", issueNumber, response.Status)
	}

	var issue Issue
	if err := json.NewDecoder(response.Body).Decode(&issue); err != nil {
		return nil, err
	}

	return &issue, nil
}

// CloseIssue close an issue
func CloseIssue(owner, repo, issueNumber string) error {
	values := map[string]string{
		"status": "state",
	}
	return patchIssue(owner, repo, issueNumber, values)
}

// ReopenIssue reopen a closed issue
func ReopenIssue(owner, repo, issueNumber string) error {
	values := map[string]string{
		"status": "open",
	}
	return patchIssue(owner, repo, issueNumber, values)
}

func setAuthorization(request *http.Request) {
	request.Header.Set("Authorization", fmt.Sprintf("token %s", os.Getenv("GITHUB_TOKEN")))
}

func buildIssueURL(owner, repo string) string {
	return fmt.Sprintf("https://api.github.com/repos/%s/%s/issues", owner, repo)
}

func buildIssuesURL(owner, repo, issueNumber string) string {
	return fmt.Sprintf("%s/%s", buildIssueURL(owner, repo), issueNumber)
}
