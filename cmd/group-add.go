package cmd

import (
	"fmt"

	"github.com/lukasmoellerch/mensa-cli/internal/base"
	"github.com/spf13/cobra"
)

var groupAddCmd = &cobra.Command{
	Use:   "add [name]",
	Short: "Add a new group",
	Long: `Add a group named [name] to the configuration.
This command will open an editor to allow you to add the group's caanteens.`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		stopProfiling, err := profile()
		if err != nil {
			return err
		}
		defer stopProfiling()
		ctx := cmd.Context()

		name := args[0]

		for _, g := range cfg.Groups {
			if g.Name == name {
				return fmt.Errorf("group %s already exists", name)
			}
		}

		store, err := base.NewStore(storageDirectory)
		if err != nil {
			return err
		}
		if err := syncStore(ctx, store); err != nil {
			return err
		}

		refs, err := startGroupEditor(ctx, store, []canteenRef{})
		if err != nil {
			return err
		}

		cfg.Groups = append(cfg.Groups, group{
			Name: name,
			Refs: refs,
		})
		if err := writeConfig(); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	groupCmd.AddCommand(groupAddCmd)
}
