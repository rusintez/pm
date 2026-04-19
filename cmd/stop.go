package cmd

import (
	"github.com/rusintez/pm/internal/tmux"
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop <session>",
	Short: "Stop a running tmux session",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return tmux.Stop(args[0])
	},
}
