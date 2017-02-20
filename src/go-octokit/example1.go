package main

import (
	"github.com/octokit/go-octokit/octokit"
)

func main() {
	client := octokit.NewClient(nil)

	url, err := octokit.UserURL.Expand(octokit.M{"user": "jingweno"})
	if err != nil {
		// Handle error
	}

	user, result := client.Users(url).One()
	if result.HasError() {
		// Handle error
	}

	fmt.Println(user.ReposURL)  // https://api.github.com/users/jingweno/repos
}

