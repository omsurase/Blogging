package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "blog-cli",
	Short: "A CLI for interacting with the blogging platform",
	Long:  `This CLI allows you to interact with various microservices of the blogging platform.`,
	Run:   runCLI,
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runCLI(cmd *cobra.Command, args []string) {
	for {
		fmt.Println("\nMain Menu:")
		fmt.Println("1. Auth")
		fmt.Println("2. Post")
		fmt.Println("3. User")
		fmt.Println("4. Exit")

		choice := getInput("Enter your choice (1-4): ")

		switch choice {
		case "1":
			runAuthMenu()
		case "2":
			runPostMenu()
		case "3":
			runUserMenu()
		case "4":
			fmt.Println("Goodbye!")
			os.Exit(0)
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}

func runAuthMenu() {
	for {
		fmt.Println("\nAuth Menu:")
		fmt.Println("1. Register")
		fmt.Println("2. Login")
		fmt.Println("3. Logout")
		fmt.Println("4. Back to main menu")

		choice := getInput("Enter your choice (1-4): ")

		switch choice {
		case "1":
			register(nil, nil)
		case "2":
			login(nil, nil)
		case "3":
			logout(nil, nil)
		case "4":
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}

func runPostMenu() {
	for {
		fmt.Println("\nPost Menu:")
		fmt.Println("1. Create post")
		fmt.Println("2. Read post")
		fmt.Println("3. Update post")
		fmt.Println("4. Delete post")
		fmt.Println("5. List posts")
		fmt.Println("6. Back to main menu")

		choice := getInput("Enter your choice (1-6): ")

		switch choice {
		case "1":
			createPost(nil, nil)
		case "2":
			id := getInput("Enter post ID: ")
			readPost(nil, []string{id})
		case "3":
			id := getInput("Enter post ID: ")
			updatePost(nil, []string{id})
		case "4":
			id := getInput("Enter post ID: ")
			deletePost(nil, []string{id})
		case "5":
			listPosts(nil, nil)
		case "6":
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}

}

func runUserMenu() {
	for {
		fmt.Println("\nUser Menu:")
		fmt.Println("1. Get user")
		fmt.Println("2. Follow user")
		fmt.Println("3. Unfollow user")
		fmt.Println("4. Get following")
		fmt.Println("5. Get followers")
		fmt.Println("6. Back to main menu")
		choice := getInput("Enter your choice (1-7): ")
		switch choice {
		case "1":
			id := getInput("Enter user ID: ")
			getUser(nil, []string{id})
		case "2":
			followerID := getInput("Enter follower ID: ")
			followeeID := getInput("Enter followee ID: ")
			followUser(nil, []string{followerID, followeeID})
		case "3":
			followerID := getInput("Enter follower ID: ")
			followeeID := getInput("Enter followee ID: ")
			unfollowUser(nil, []string{followerID, followeeID})
		case "4":
			id := getInput("Enter user ID: ")
			getFollowing(nil, []string{id})
		case "5":
			id := getInput("Enter user ID: ")
			getFollowers(nil, []string{id})
		case "6":
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}
