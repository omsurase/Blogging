package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
)

var postCmd = &cobra.Command{
	Use:   "post",
	Short: "Manage blog posts",
	Long:  `Create, read, update, and delete blog posts.`,
}

var createPostCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new blog post",
	Run:   createPost,
}

var readPostCmd = &cobra.Command{
	Use:   "read [id]",
	Short: "Read a blog post",
	Args:  cobra.ExactArgs(1),
	Run:   readPost,
}

var updatePostCmd = &cobra.Command{
	Use:   "update [id]",
	Short: "Update a blog post",
	Args:  cobra.ExactArgs(1),
	Run:   updatePost,
}

var deletePostCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete a blog post",
	Args:  cobra.ExactArgs(1),
	Run:   deletePost,
}

var listPostsCmd = &cobra.Command{
	Use:   "list",
	Short: "List all blog posts",
	Run:   listPosts,
}

func init() {
	postCmd.AddCommand(createPostCmd)
	postCmd.AddCommand(readPostCmd)
	postCmd.AddCommand(updatePostCmd)
	postCmd.AddCommand(deletePostCmd)
	postCmd.AddCommand(listPostsCmd)
}

// func getInput(prompt string) string {
// 	reader := bufio.NewReader(os.Stdin)
// 	fmt.Print(prompt)
// 	input, _ := reader.ReadString('\n')
// 	return strings.TrimSpace(input)
// }

func createPost(cmd *cobra.Command, args []string) {
	title := getInput("Enter post title: ")
	content := getInput("Enter post content: ")
	author := getInput("Enter post author: ")

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

	resp, err := http.Post("http://localhost:8080/posts", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating post:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Post created successfully")
}

func readPost(cmd *cobra.Command, args []string) {
	id := args[0]
	resp, err := http.Get(fmt.Sprintf("http://localhost:8080/posts/%s", id))
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
	title := getInput("Enter new post title (press Enter to keep current): ")
	content := getInput("Enter new post content (press Enter to keep current): ")
	author := getInput("Enter new post author (press Enter to keep current): ")

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

	client := &http.Client{}
	req, err := http.NewRequest("PUT", fmt.Sprintf("http://localhost:8080/posts/%s", id), bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
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
	client := &http.Client{}
	req, err := http.NewRequest("DELETE", fmt.Sprintf("http://localhost:8080/posts/%s", id), nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error deleting post:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Post deleted successfully")
}

func listPosts(cmd *cobra.Command, args []string) {
	resp, err := http.Get("http://localhost:8080/posts")
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
