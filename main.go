package main

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"
	"time"

	"github.com/mauricioabreu/goplex4_11/editor"
	"github.com/mauricioabreu/goplex4_11/github"
)

var issueTemplate = template.Must(template.New("issue").Funcs(template.FuncMap{"formatTime": formatTime}).Parse(`
Number:	{{.Number}}
Title: {{.Title}} 
Body: {{.Body}}
Created:  {{.CreatedAt | formatTime}}
Updated:  {{.UpdatedAt | formatTime}}
`))

func main() {
	printUsage()
	args := os.Args[1:]
	action := args[0]
	owner := args[1]
	repo := args[2]
	issueNumber := args[3]

	switch action {
	case "create":
		create(owner, repo)
	case "update":
		update(owner, repo, issueNumber)
	case "read":
		read(owner, repo, issueNumber)
	case "close":
		close(owner, repo, issueNumber)
	case "reopen":
		reopen(owner, repo, issueNumber)
	}
}

func create(owner, repo string) {
	input, err := editor.Edit("")
	if err != nil {
		log.Fatal(err)
	}
	title, body := parseText(input)
	err = github.CreateIssue(owner, repo, title, body)
	if err != nil {
		log.Fatal(err)
	}
}

func update(owner, repo, issueNumber string) {
	issue, err := github.GetIssue(owner, repo, issueNumber)
	if err != nil {
		log.Fatal(err)
	}
	content := fmt.Sprintf("%s\n\n%s", string(issue.Title), string(issue.Body))
	input, err := editor.Edit(content)
	if err != nil {
		log.Fatal(err)
	}
	title, body := parseText(input)
	err = github.UpdateIssue(owner, repo, title, body, issueNumber)
	if err != nil {
		log.Fatal(err)
	}
}

func read(owner, repo, issueNumber string) {
	issue, err := github.GetIssue(owner, repo, issueNumber)
	if err != nil {
		log.Fatal(err)
	}
	err = issueTemplate.Execute(os.Stdout, issue)
	if err != nil {
		log.Fatal(err)
	}
}

func close(owner, repo, issueNumber string) {
	err := github.CloseIssue(owner, repo, issueNumber)
	if err != nil {
		log.Fatal(err)
	}
}

func reopen(owner, repo, issueNumber string) {
	err := github.ReopenIssue(owner, repo, issueNumber)
	if err != nil {
		log.Fatal(err)
	}
}

func parseText(input []byte) (string, string) {
	content := strings.Split(string(input), "\n")
	title := content[0]
	body := strings.TrimSpace(strings.Join(content[1:], "\n"))
	return title, body
}

func printUsage() {
	fmt.Println("Usage: ./gcrud create|update|delete|read|close|reopen <args>")
}

func formatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}
