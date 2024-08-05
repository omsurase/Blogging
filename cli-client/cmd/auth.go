package cmd

//./blog-cli auth register -u your_username -p your_password -e your_email@example.com
//./blog-cli auth login -u yourusername -p yourpassword
import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/spf13/cobra"
)

var (
	username string
	password string
	email    string
	token    string
)

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authenticate with the blogging platform",
	Long:  `Use this command to log in or log out of the blogging platform.`,
}

var registerCmd = &cobra.Command{
	Use:   "register",
	Short: "Register a new user",
	Run:   register,
}

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Log in to the blogging platform",
	Run:   login,
}

var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "Log out from the blogging platform",
	Run:   logout,
}

func init() {
	registerCmd.Flags().StringVarP(&username, "username", "u", "", "Username for registration")
	registerCmd.Flags().StringVarP(&password, "password", "p", "", "Password for registration")
	registerCmd.Flags().StringVarP(&email, "email", "e", "", "Email for registration")
	registerCmd.MarkFlagRequired("username")
	registerCmd.MarkFlagRequired("password")
	registerCmd.MarkFlagRequired("email")

	loginCmd.Flags().StringVarP(&username, "username", "u", "", "Username for authentication")
	loginCmd.Flags().StringVarP(&password, "password", "p", "", "Password for authentication")
	loginCmd.MarkFlagRequired("username")
	loginCmd.MarkFlagRequired("password")

	authCmd.AddCommand(registerCmd)
	authCmd.AddCommand(loginCmd)
	authCmd.AddCommand(logoutCmd)
}

func register(cmd *cobra.Command, args []string) {
	// Prepare the request body
	reqBody, err := json.Marshal(map[string]string{
		"username": username,
		"password": password,
		"email":    email,
	})
	if err != nil {
		fmt.Println("Error preparing request:", err)
		return
	}
	log.Printf("cli app")
	// Send POST request to the auth microservice
	resp, err := http.Post("http://localhost:8080/auth/register", "application/json", bytes.NewBuffer(reqBody))

	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	// Parse the response
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	// Check if registration was successful
	if resp.StatusCode == http.StatusCreated {
		fmt.Println("Registration successful:", result["message"])

		// Check if a token was returned
		if token, ok := result["token"].(string); ok {
			fmt.Println("JWT Token:", token)
			// Here you might want to save the token securely for future use
		} else {
			fmt.Println("No token was returned with the registration")
		}
	} else {
		fmt.Println("Registration failed:", result["message"])
	}
}

func login(cmd *cobra.Command, args []string) {
	// Prepare the request body
	reqBody, err := json.Marshal(map[string]string{
		"username": username,
		"password": password,
	})
	if err != nil {
		fmt.Println("Error preparing request:", err)
		return
	}

	// Send POST request to the auth microservice
	resp, err := http.Post("http://localhost:8080/auth/login", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	// Parse the response
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	// Check if login was successful
	if resp.StatusCode == http.StatusOK {
		token = result["token"].(string)
		fmt.Println("Login successful. Token:", token)
		// Here you might want to save the token securely for future use
	} else {
		fmt.Println("Login failed:", result["message"])
	}
}

func logout(cmd *cobra.Command, args []string) {
	// Prepare the request
	req, err := http.NewRequest("POST", "http://localhost:8080/auth/logout", nil)
	if err != nil {
		fmt.Println("Error preparing request:", err)
		return
	}

	// Add the token to the request header
	req.Header.Add("Authorization", "Bearer "+token)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	// Parse the response
	var result map[string]interface{}
	json.Unmarshal(body, &result)

	// Check if logout was successful
	if resp.StatusCode == http.StatusOK {
		token = "" // Clear the token
		fmt.Println("Logout successful")
	} else {
		fmt.Println("Logout failed:", result["message"])
	}
}
