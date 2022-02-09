package cmd

import (
	"fmt"

	"github.com/lukasmoellerch/mensa-cli/internal/base"
	"github.com/spf13/cobra"
)

var groupUpdateCmd = &cobra.Command{
	Use:   "update [name]",
	Short: "Update a group",
	Long:  `Opens the group editor for the group with the given name and allows you to make changes to it.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name := args[0]

		stopProfiling, err := profile()
		if err != nil {
			return err
		}
		defer stopProfiling()
		ctx := cmd.Context()

		store, err := base.NewStore(storageDirectory)
		if err != nil {
			panic(err)
		}
		if err := syncStore(ctx, store); err != nil {
			return err
		}

		var index int
		var found *group
		for i, g := range cfg.Groups {
			if g.Name == name {
				index = i
				found = &g
				break
			}
		}
		if found == nil {
			return fmt.Errorf("group %s not found", name)
		}

		refs, err := startGroupEditor(ctx, store, found.Refs)
		if err != nil {
			return err
		}

		cfg.Groups[index].Refs = refs
		if err := writeConfig(); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	groupCmd.AddCommand(groupUpdateCmd)
}
