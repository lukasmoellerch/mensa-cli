package cmd

import (
	"github.com/spf13/cobra"
)

var groupRemoveCmd = &cobra.Command{
	Use:   "remove [name]",
	Short: "Remove a group",
	Long:  `Remvoes the group with the given name from the configuration file.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]
		newGroups := []group{}
		for _, g := range cfg.Groups {
			if g.Name != name {
				newGroups = append(newGroups, g)
			}
		}
		cfg.Groups = newGroups
		if err := writeConfig(); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	groupCmd.AddCommand(groupRemoveCmd)
}
