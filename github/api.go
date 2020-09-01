package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

// Issue contains issue data
type Issue struct {
	Number int
	Title  string
	Body   string
}

// CreateIssue create a new issue
func CreateIssue(owner, repo, title, body string) {
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

	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(respBody))
}

// UpdateIssue update an issue
func UpdateIssue(owner, repo, title, body, issueNumber string) {
	requestBody, err := json.Marshal(map[string]string{
		"title": title,
		"body":  body,
	})
	if err != nil {
		log.Fatal(err)
	}

	request, err := http.NewRequest("PATCH", buildIssuesURL(owner, repo, issueNumber), bytes.NewBuffer(requestBody))
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

	respBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	log.Println(string(respBody))
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

func setAuthorization(request *http.Request) {
	request.Header.Set("Authorization", fmt.Sprintf("token %s", os.Getenv("GITHUB_TOKEN")))
}

func buildIssueURL(owner, repo string) string {
	return fmt.Sprintf("https://api.github.com/repos/%s/%s/issues", owner, repo)
}

func buildIssuesURL(owner, repo, issueNumber string) string {
	return fmt.Sprintf("%s/%s", buildIssueURL(owner, repo), issueNumber)
}
