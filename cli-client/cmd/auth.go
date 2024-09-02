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
	"gopkg.in/yaml.v2"
)

type UserConfig struct {
	Username string `yaml:"username"`
	Token    string `yaml:"token"`
	UserID   string `yaml:"user_id"` // Added field for user ID
}

var (
	username string
	password string
	email    string
	token    string
	userID   string
)

const configFile = "user_config.yml"

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

func saveConfig(config UserConfig) error {
	data, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}
	return os.WriteFile(configFile, data, 0600)
}

func loadConfig() (*UserConfig, error) {
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	var config UserConfig
	err = yaml.Unmarshal(data, &config)
	return &config, err
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
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Error parsing response:", err)
		return
	}

	if resp.StatusCode == http.StatusCreated {
		fmt.Println("Registration successful:", result["message"])
		if token, ok := result["token"].(string); ok {
			if userID, ok := result["user_id"].(string); ok {
				config := UserConfig{
					Username: username,
					Token:    token,
					UserID:   userID,
				}
				if err := saveConfig(config); err != nil {
					fmt.Println("Error saving config:", err)
				} else {
					fmt.Println("User configuration saved successfully.")
				}
			} else {
				fmt.Println("No user ID was returned with the registration")
			}
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
	if err := json.Unmarshal(body, &result); err != nil {
		fmt.Println("Error parsing response:", err)
		return
	}

	if resp.StatusCode == http.StatusOK {
		if token, ok := result["token"].(string); ok {
			if userID, ok := result["user_id"].(string); ok {
				config := UserConfig{
					Username: username,
					Token:    token,
					UserID:   userID,
				}
				if err := saveConfig(config); err != nil {
					fmt.Println("Error saving config:", err)
				} else {
					fmt.Println("Login successful. User configuration saved.")
				}
			} else {
				fmt.Println("No user ID was returned with the login")
			}
		} else {
			fmt.Println("No token was returned with the login")
		}
	} else {
		fmt.Println("Login failed:", result["message"])
	}
}

func logout(cmd *cobra.Command, args []string) {
	if err := os.Remove(configFile); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("No active session found.")
		} else {
			fmt.Println("Error removing config file:", err)
		}
	} else {
		fmt.Println("Logout successful. User configuration removed.")
	}
}
