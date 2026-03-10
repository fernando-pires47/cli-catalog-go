package main

import (
	"fmt"
	"os"

	"command-cli/internal/cli"
)

func main() {
	app, err := cli.NewApp()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	exitCode := app.Run(os.Args[1:])
	os.Exit(exitCode)
}
