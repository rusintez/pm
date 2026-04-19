package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/rusintez/pm/internal/config"
	"github.com/rusintez/pm/internal/runner"
	"github.com/rusintez/pm/internal/tmux"
	"github.com/spf13/cobra"
)

var useTmux bool

var runCmd = &cobra.Command{
	Use:   "run <layout | process...>",
	Short: "Run a layout or set of processes",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		names, layoutEnv, err := cfg.Resolve(args)
		if err != nil {
			return err
		}

		if useTmux {
			session := strings.Join(args, "-")
			return tmux.Run(cfg, session, names, layoutEnv)
		}
		return runner.Run(cfg, names, layoutEnv)
	},
}

func init() {
	runCmd.Flags().BoolVar(&useTmux, "tmux", false, "run in a tmux session")
}
