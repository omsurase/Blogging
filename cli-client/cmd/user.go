package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type Config3 struct {
	UserID string `yaml:"user_id"`
}

func loadConfig3() (Config3, error) {
	var config Config3
	configFile, err := ioutil.ReadFile("user_config.yml")
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal(configFile, &config)
	return config, err
}

func getUser(cmd *cobra.Command, args []string) {
	config, err := loadConfig3()
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
	config, err := loadConfig3()
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
	config, err := loadConfig3()
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
	config, err := loadConfig3()
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
	config, err := loadConfig3()
	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/users/%s/followers", config.UserID))
	if err != nil {
		fmt.Println("Error getting followers list:", err)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("Followers list:", string(body))
}
