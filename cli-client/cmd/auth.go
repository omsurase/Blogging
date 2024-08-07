package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var (
	username string
	password string
	email    string
	token    string
)

func getInput(prompt string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

func getPassword(prompt string) (string, error) {
	fmt.Print(prompt)
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println() // Add a newline after the password input
	return string(password), err
}

func register(cmd *cobra.Command, args []string) {
	username = getInput("Enter username: ")
	email = getInput("Enter email: ")
	password, err := getPassword("Enter password: ")
	if err != nil {
		fmt.Println("Error reading password:", err)
		return
	}

	reqBody, err := json.Marshal(map[string]string{
		"username": username,
		"password": password,
		"email":    email,
	})
	if err != nil {
		fmt.Println("Error preparing request:", err)
		return
	}

	resp, err := http.Post("http://localhost:8080/auth/register", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	var result map[string]interface{}
	json.Unmarshal(body, &result)

	if resp.StatusCode == http.StatusCreated {
		fmt.Println("Registration successful:", result["message"])
		if token, ok := result["token"].(string); ok {
			fmt.Println("JWT Token:", token)
		} else {
			fmt.Println("No token was returned with the registration")
		}
	} else {
		fmt.Println("Registration failed:", result["message"])
	}
}

func login(cmd *cobra.Command, args []string) {
	username = getInput("Enter username: ")
	password, err := getPassword("Enter password: ")
	if err != nil {
		fmt.Println("Error reading password:", err)
		return
	}

	reqBody, err := json.Marshal(map[string]string{
		"username": username,
		"password": password,
	})
	if err != nil {
		fmt.Println("Error preparing request:", err)
		return
	}

	resp, err := http.Post("http://localhost:8080/auth/login", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	var result map[string]interface{}
	json.Unmarshal(body, &result)

	if resp.StatusCode == http.StatusOK {
		token = result["token"].(string)
		fmt.Println("Login successful. Token:", token)
	} else {
		fmt.Println("Login failed:", result["message"])
	}
}

func logout(cmd *cobra.Command, args []string) {
	req, err := http.NewRequest("POST", "http://localhost:8080/auth/logout", nil)
	if err != nil {
		fmt.Println("Error preparing request:", err)
		return
	}

	req.Header.Add("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	var result map[string]interface{}
	json.Unmarshal(body, &result)

	if resp.StatusCode == http.StatusOK {
		token = "" // Clear the token
		fmt.Println("Logout successful")
	} else {
		fmt.Println("Logout failed:", result["message"])
	}
}
