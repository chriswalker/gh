package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/shurcooL/githubv4"
	"github.com/shurcooL/graphql"
	"golang.org/x/oauth2"
)

var (
	token = flag.String("token", "", "Personal Access Token fo Github")
)

const (
	GithubQLURL = "https://api.github.com/graphql"
)

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
	History `graphql:"history(first: 10, after: $commit)"`
}

type History struct {
	Edges
}

type Edges struct {
	Node
}

type Node struct {
	Oid, CommitDate, MessageHeadline githubv4.String
}

func main() {
	flag.Parse()
	if *token == "" {
		log.Fatal("token required")
	}

	src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: *token})

	httpClient := oauth2.NewClient(context.Background(), src)
	client := graphql.NewClient(GithubQLURL, httpClient)
	/*
	   query ($default_branch: String!, $commit: String!) {
	     repository(owner: "chriswalker", name: "dotfiles") {
	       object(expression: $default_branch) {
	         ... on Commit {
	           history(first: 10, after: $commit) {
	             edges {
	               node {
	                 id
	                 committedDate
	                 messageHeadline
	               }
	             }
	           }
	         }
	       }
	     }
	   }
	*/
	vars := map[string]interface{}{
		"default_branch": "master",
		"commit":         "4ef34f27fdd5c39d7f3cfb9012251d325c9900e9",
		"repo":           "dotfiles",
		"owner":          "chriswalker",
	}

	query := Query{}
	// update to pass in vars when ready
	err := client.Query(context.Background(), &query, vars)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(query)
}
