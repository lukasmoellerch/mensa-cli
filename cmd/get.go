package cmd

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/lukasmoellerch/mensa-cli/internal/components"
	"github.com/lukasmoellerch/mensa-cli/internal/zuerichmensa"
	"github.com/spf13/cobra"
)

var dateFlag string
var daytimeFlag string

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Fetches the list of meals for a given date",
	Long: `A CLI tool which fetches the list of meals for a given date.

If no date is specified, the current date is used.

If no daytime is specified, the current daytime is used.

If no filter is specified, all facilities are fetched.

The filter is matched against the name of the facility using a substring check.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		store, err := zuerichmensa.NewStore("/tmp/facilties")
		if err != nil {
			return err
		}

		if store.IsEmpty() {
			err = store.Sync(cmd.Context())
			if err != nil {
				return err
			}
		}

		now := time.Now()
		if dateFlag == "" {
			dateFlag = now.Format("2006-01-02")
		}
		if daytimeFlag == "" {
			if now.Hour() < 14 {
				daytimeFlag = "lunch"
			} else {
				daytimeFlag = "dinner"
			}
		}
		filter := ""
		if len(args) == 1 {
			filter = args[0]
		}

		facilities := make([]int, 0)
		res, err := store.Filter(cmd.Context(), filter)
		if err != nil {
			return err
		}
		for _, f := range res {
			facilities = append(facilities, f.ID)
		}

		p := tea.NewProgram(components.LoaderModel{
			Ctx:      cmd.Context(),
			MensaIds: facilities,
			Lang:     "de",
			Date:     dateFlag,
			Daytime:  daytimeFlag,
		})
		if err := p.Start(); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(getCmd)
	getCmd.Flags().StringVarP(&dateFlag, "date", "d", "", "Date in format YYYY-MM-DD")
	getCmd.Flags().StringVarP(&daytimeFlag, "daytime", "t", "", "Daytime (dinner, lunch, current)")
}
