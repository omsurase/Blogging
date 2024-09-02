package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// Config2 struct to hold the YAML config
type Config2 struct {
	UserID string `yaml:"user_id"`
}

var config Config2

func init() {
	// Load the configuration when the package is initialized
	loadConfig3()
}

func loadConfig3() {
	file, err := os.Open("user_config.yml")
	if err != nil {
		fmt.Println("Error opening config file:", err)
		return
	}
	defer file.Close()

	// Read the file content for debugging
	content, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Error reading config file:", err)
		return
	}
	fmt.Println("Config file content:", string(content))

	// Decode the YAML content
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		fmt.Println("Error decoding config file:", err)
		return
	}

	// Print the loaded config for debugging
	fmt.Printf("Loaded Config: %+v\n", config)
}

func getUser(cmd *cobra.Command, args []string) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/users/%s", config.UserID))
	if err != nil {
		fmt.Println("Error getting user:", err)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("User details:", string(body))
}

func followUser(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		fmt.Println("Please provide the followee ID")
		return
	}

	url := fmt.Sprintf("http://localhost:8080/users/%s/follow/%s", config.UserID, args[0])
	fmt.Println(url)
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
	if len(args) < 1 {
		fmt.Println("Please provide the followee ID")
		return
	}

	url := fmt.Sprintf("http://localhost:8080/users/%s/unfollow/%s", config.UserID, args[0])
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
	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/users/%s/following", config.UserID))
	if err != nil {
		fmt.Println("Error getting following list:", err)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Following list:", string(body))
}

func getFollowers(cmd *cobra.Command, args []string) {
	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/users/%s/followers", config.UserID))
	if err != nil {
		fmt.Println("Error getting followers list:", err)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Followers list:", string(body))
}
