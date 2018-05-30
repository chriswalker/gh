package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const (
	GithubToken = "c39cf226885fabfe918911726afec955a900816e"
	GithubQLURL = "https://api.github.com/graphql"
	Body        = `{
		"query": "query ($default_branch: String!, $commit: String!) {
  repository(owner: \"chriswalker\", name: \"dotfiles\") {
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

variables {
  \"default_branch\": \"master\",
  \"commit\": \"4ef34f27fdd5c39d7f3cfb9012251d325c9900e9\"
}"
	}`

	Body2 = "{ \"query\": \"query { viewer { login } }\" }"
)

func main() {
	client := &http.Client{}

	req, err := http.NewRequest("POST", GithubQLURL, strings.NewReader(Body))
	req.Header.Add("Authorization", "bearer "+GithubToken)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// else, dump response
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(string(bytes))
}
