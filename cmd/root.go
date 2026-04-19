package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var root = &cobra.Command{
	Use:   "pm",
	Short: "Superlight process manager",
}

func Execute() {
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	root.AddCommand(runCmd, listCmd, attachCmd, stopCmd)
}
