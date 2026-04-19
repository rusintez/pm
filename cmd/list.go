package cmd

import (
	"fmt"
	"sort"

	"github.com/rusintez/pm/internal/config"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List processes and layouts",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		fmt.Println("Processes:")
		names := sortedKeys(cfg.Processes)
		for _, name := range names {
			fmt.Printf("  %-16s %s\n", name, cfg.Processes[name].Cmd)
		}

		fmt.Println("\nLayouts:")
		layouts := sortedKeys(cfg.Layouts)
		for _, name := range layouts {
			fmt.Printf("  %-16s %v\n", name, cfg.Layouts[name].Processes)
		}
		return nil
	},
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
