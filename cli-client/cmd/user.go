package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/cobra"
)

func getUser(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Println("Please provide a user ID")
		return
	}

	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/users/%s", args[0]))
	if err != nil {
		fmt.Println("Error getting user:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("User details:", string(body))
}

func followUser(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		fmt.Println("Please provide follower ID and followee ID")
		return
	}

	url := fmt.Sprintf("http://localhost:8080/users/%s/follow/%s", args[0], args[1])
	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		fmt.Println("Error following user:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Successfully followed user")
	} else {
		fmt.Println("Failed to follow user")
	}
}

func unfollowUser(cmd *cobra.Command, args []string) {
	if len(args) < 2 {
		fmt.Println("Please provide follower ID and followee ID")
		return
	}

	url := fmt.Sprintf("http://localhost:8080/users/%s/unfollow/%s", args[0], args[1])
	resp, err := http.Post(url, "application/json", nil)
	if err != nil {
		fmt.Println("Error unfollowing user:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Successfully unfollowed user")
	} else {
		fmt.Println("Failed to unfollow user")
	}
}

func getFollowing(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Println("Please provide a user ID")
		return
	}

	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/users/%s/following", args[0]))
	if err != nil {
		fmt.Println("Error getting following list:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Following list:", string(body))
}

func getFollowers(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Println("Please provide a user ID")
		return
	}

	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/users/%s/followers", args[0]))
	if err != nil {
		fmt.Println("Error getting followers list:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Followers list:", string(body))
}
