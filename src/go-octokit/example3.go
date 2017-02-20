package main

import "github.com/octokit/go-octokit/octokit"

func main() {
	url, err := octokit.UserURL.Expand(nil)
	if err != nil  {
		// Handle error
	}

	users, result := client.Users(url).All()
	if result.HasError() {
		// Handle error
	}

	// Do something with users

	// Next page
	nextPageURL, _ := result.NextPage.Expand(nil)
	users, result := client.Users(nextPageURL).All()
	if result.HasError() {
		// Handle error
	}

	// Do something with users
}