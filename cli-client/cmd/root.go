package cmd

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "blog-cli",
	Short: "A CLI for interacting with the blogging platform",
	Long:  `This CLI allows you to interact with various microservices of the blogging platform.`,
}

func init() {
	RootCmd.AddCommand(authCmd)
	RootCmd.AddCommand(postCmd)
}
