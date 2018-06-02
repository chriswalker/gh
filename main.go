package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/fatih/color"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

var (
	token  = flag.String("token", "", "Personal Access Token fo Github")
	commit = flag.String("commit", "", "Commit to search from")
)

const (
	GithubQLURL = "https://api.github.com/graphql"
)

// Alias for fatih/color SprintFunc() returns
type SprintfFunc func(format string, a ...interface{}) string

type Query struct {
	Repository `graphql:"repository(owner: $owner, name: $repo)"`
}

type Repository struct {
	Object `graphql:"object(expression: $default_branch)"`
}

type Object struct {
	Commit `graphql:"... on Commit"`
}

type Commit struct {
	History `graphql:"history(first: 10)"`
}

type History struct {
	Edges
}

type Edges struct {
	Nodes []CommitNode
}

type CommitNode struct {
	Oid, CommittedDate, MessageHeadline githubv4.String
}

func main() {
	flag.Parse()
	if *token == "" {
		log.Fatal("token required")
	}
	if *commit == "" {
		log.Fatal("commit required")
	}

	src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: *token})

	httpClient := oauth2.NewClient(context.Background(), src)
	client := githubv4.NewClient(httpClient)

	vars := map[string]interface{}{
		"default_branch": githubv4.String("master"),
		"repo":           githubv4.String("dotfiles"),
		"owner":          githubv4.String("chriswalker"),
	}

	query := Query{}
	err := client.Query(context.Background(), &query, vars)
	if err != nil {
		log.Fatal(err)
	}

	for i, commitNode := range query.Repository.History.Edges.Nodes {
		oid := commitNode.Oid
		if commitNode.Oid == githubv4.String(*commit) {
			fn := GetColourFunc(i)
			fmt.Printf("Deployed: [%s] %s\n", color.BlueString(*commit), fn("Master: [%s], ahead by %d commits", oid, i))
		}
	}
}

func GetColourFunc(num int) func(format string, a ...interface{}) string {
	if num < 3 {
		return color.New(color.FgGreen).SprintfFunc()
	} else if num < 6 {
		return color.New(color.FgYellow).SprintfFunc()
	}
	return color.New(color.FgRed).SprintfFunc()
}
