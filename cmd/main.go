package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var (
	argJSONOutput = false
	argInput      = ""
)

var rootCmd = &cobra.Command{
	Use:   CommandName,
	Short: "ChatGPT command line tool that supports pipe and repl.",
	Run:   run,
}

func run(_ *cobra.Command, args []string) {
	argInput = strings.Join(args, " ")
	_shared = newApp()
	_shared.run()
}

func main() {
	err := initConfig()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, fmt.Sprintf("error: %s\n", err.Error()))
		os.Exit(1)
		return
	}

	rootCmd.PersistentFlags().BoolVarP(&argJSONOutput, "json", "j", false, "output as json")

	if err := rootCmd.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
