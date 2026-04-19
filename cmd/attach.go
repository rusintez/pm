package cmd

import (
	"github.com/rusintez/pm/internal/tmux"
	"github.com/spf13/cobra"
)

var attachCmd = &cobra.Command{
	Use:   "attach <session>",
	Short: "Attach to a running tmux session",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return tmux.Attach(args[0])
	},
}
