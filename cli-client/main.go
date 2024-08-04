package main

import (
	"fmt"
	"github.com/omsurase/Blogging/cli-client/cmd"
	"os"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
