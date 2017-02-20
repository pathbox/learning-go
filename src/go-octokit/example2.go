package main

import (
	"github.com/octokit/go-octokit/octokit"
)

func main() {
	rootURL, _ := client.RootURL.Expand(nil)
	root, _ := client.Root(rootURL).One()

	userURL, _ := root.UserURL.Expand(octokit.M{"users": "jingweno"})
	user, _ := client.Users(userURL).One()
}

