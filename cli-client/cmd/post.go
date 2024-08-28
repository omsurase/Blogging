package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Token string `yaml:"token"`
}

func loadConfig() (Config, error) {
	var config Config
	configFile, err := ioutil.ReadFile("user_config.yml")
	if err != nil {
		return config, err
	}
	err = yaml.Unmarshal(configFile, &config)
	return config, err
}

func createAuthenticatedRequest(method, url string, body []byte) (*http.Request, error) {
	config, err := loadConfig()
	if err != nil {
		return nil, fmt.Errorf("error loading config: %v", err)
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.Token)
	return req, nil
}

func createPost(cmd *cobra.Command, args []string) {
	title := getInput2("Enter post title: ")
	content := getInput2("Enter post content: ")
	author := getInput2("Enter post author: ")
	post := map[string]string{
		"title":   title,
		"content": content,
		"author":  author,
	}
	jsonData, err := json.Marshal(post)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	req, err := createAuthenticatedRequest("POST", "http://localhost:8080/posts", jsonData)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error creating post:", err)
		return
	}
	defer resp.Body.Close()
	fmt.Println("Post created successfully")
}

func readPost(cmd *cobra.Command, args []string) {
	id := args[0]
	req, err := createAuthenticatedRequest("GET", fmt.Sprintf("http://localhost:8080/posts/%s", id), nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error reading post:", err)
		return
	}
	defer resp.Body.Close()

	var post map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&post)
	fmt.Printf("Title: %s\nContent: %s\nAuthor: %s\n", post["title"], post["content"], post["author"])
}

func updatePost(cmd *cobra.Command, args []string) {
	id := args[0]
	title := getInput2("Enter new post title (press Enter to keep current): ")
	content := getInput2("Enter new post content (press Enter to keep current): ")
	author := getInput2("Enter new post author (press Enter to keep current): ")
	post := map[string]string{}
	if title != "" {
		post["title"] = title
	}
	if content != "" {
		post["content"] = content
	}
	if author != "" {
		post["author"] = author
	}
	jsonData, err := json.Marshal(post)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	req, err := createAuthenticatedRequest("PUT", fmt.Sprintf("http://localhost:8080/posts/%s", id), jsonData)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error updating post:", err)
		return
	}
	defer resp.Body.Close()
	fmt.Println("Post updated successfully")
}

func deletePost(cmd *cobra.Command, args []string) {
	id := args[0]
	req, err := createAuthenticatedRequest("DELETE", fmt.Sprintf("http://localhost:8080/posts/%s", id), nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error deleting post:", err)
		return
	}
	defer resp.Body.Close()
	fmt.Println("Post deleted successfully")
}

func listPosts(cmd *cobra.Command, args []string) {
	req, err := createAuthenticatedRequest("GET", "http://localhost:8080/posts", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error listing posts:", err)
		return
	}
	defer resp.Body.Close()

	var posts []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&posts)
	for _, post := range posts {
		fmt.Printf("ID: %s\nTitle: %s\nAuthor: %s\n\n", post["id"], post["title"], post["author"])
	}
}

func getInput2(prompt string) string {
	fmt.Print(prompt)
	var input string
	fmt.Scanln(&input)
	return input
}
