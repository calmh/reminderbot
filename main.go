package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/calmh/github"
)

var (
	repo               = "syncthing/syncthing"
	maxUnclassifiedAge = 30 * 24 * time.Hour
	minIdle            = 7 * 24 * time.Hour
)

func main() {
	flag.StringVar(&repo, "repo", repo, "Repository")
	flag.Parse()
	triage()
	oldBugs()
}

func triage() {
	header("Issues needing triage")

	query := make(url.Values)
	query.Add("milestone", "none")
	issues, err := github.LoadIssues(repo, query)
	if err != nil {
		log.Fatal(err)
	}
	for _, i := range issues {
		if !shouldProcess(i) {
			continue
		}
		if time.Since(i.Created) > maxUnclassifiedAge {
			issue(i)
		}
	}
}

func oldBugs() {
	header("Old bugs")

	query := make(url.Values)
	query.Add("labels", "bug")
	issues, err := github.LoadIssues(repo, query)
	if err != nil {
		log.Fatal(err)
	}
	for _, i := range issues {
		if !shouldProcess(i) {
			continue
		}
		if time.Since(i.Created) > 365*24*time.Hour {
			issue(i)
		}
	}
}

func issue(i github.Issue) {
	fmt.Printf("#%d: %s\n> %s\n\n", i.Number, i.Title, i.HTMLURL)
}

func header(s string) {
	underline := make([]byte, len(s))
	for i := range underline {
		underline[i] = '='
	}
	fmt.Printf("%s\n%s\n\n", s, underline)
}

func shouldProcess(i github.Issue) bool {
	if i.PullRequest.URL != "" {
		return false
	}
	return time.Since(i.Updated) > minIdle
}
