package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"github.com/shurcooL/graphql"
	"golang.org/x/oauth2"
)

var (
	token = flag.String("token", "", "Personal Access Token fo Github")
)

const (
	GithubQLURL = "https://api.github.com/graphql"
)

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
	var query struct {
		Viewer struct {
			Login graphql.String
			Name  graphql.String
		}
	}

	vars := map[string]interface{}{
		"default_branch": "master",
		"commit":         "4ef34f27fdd5c39d7f3cfb9012251d325c9900e9",
		"repo":           "dotfiles",
		"owner":          "chriswalker",
	}

	// update to pass in vars when ready
	err := client.Query(context.Background(), &query, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(query.Viewer)
}
