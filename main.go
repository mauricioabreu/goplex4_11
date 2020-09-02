package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mauricioabreu/goplex4_11/editor"
	"github.com/mauricioabreu/goplex4_11/github"
)

func main() {
	printUsage()
	args := os.Args[1:]
	action := args[0]
	owner := args[1]
	repo := args[2]
	issueNumber := args[3]
	if action == "create" {
		input, err := editor.Edit("")
		if err != nil {
			log.Fatal(err)
		}
		title, body := parseText(input)
		github.CreateIssue(owner, repo, title, body)
	} else if action == "update" {
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
		github.UpdateIssue(owner, repo, title, body, issueNumber)
	}
}

func parseText(input []byte) (string, string) {
	content := strings.Split(string(input), "\n")
	title := content[0]
	body := strings.Join(content[1:], "\n")
	return title, body
}

func printUsage() {
	fmt.Println("Usage: ./gcrud create|read|update|delete <args>")
}
