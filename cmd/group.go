package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var groupCmd = &cobra.Command{
	Use:   "group",
	Short: "Manage groups",
	Long: `List groups defined in the configuration file in use,
different groups are separated by a newline`,
	RunE: func(cmd *cobra.Command, args []string) error {
		stopProfiling, err := profile()
		if err != nil {
			return err
		}
		defer stopProfiling()

		// List groups
		for _, g := range cfg.Groups {
			fmt.Printf("%s\n", g.Name)
		}

		return nil
	},
}

func init() {
	RootCmd.AddCommand(groupCmd)
}
