package cmd

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/lukasmoellerch/mensa-cli/internal/base"
	"github.com/lukasmoellerch/mensa-cli/internal/protobuf/storage"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

var mealsDateFlag string
var mealsDaytimeFlag string
var mealsFilterFlag string
var mealsGroupFlag string
var mealsDinnerFlag bool
var mealsLunchFlag bool
var mealsTomorrowFlag bool

// mealsCmd represents the get command
var mealsCmd = &cobra.Command{
	Use:   "meals",
	Short: "Fetches the list of meals for a given date",
	Long: `Fetches the list of meals for a given date.
If no date is specified, the current date is used.
If no daytime is specified, the current daytime is used.
If the filter is empty, all facilities are fetched. The filter is matched against the name of the facility using a substring check.`,
	Args: cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		stopProfiling, err := profile()
		if err != nil {
			return err
		}
		defer stopProfiling()

		if mealsGroupFlag != "" && mealsFilterFlag != "" {
			return fmt.Errorf("cannot use both --group and --filter")
		}
		if mealsDinnerFlag && mealsDaytimeFlag != "" {
			return fmt.Errorf("cannot use both --dinner and --daytime")
		}
		if mealsLunchFlag && mealsDaytimeFlag != "" {
			return fmt.Errorf("cannot use both --lunch and --daytime")
		}
		if mealsLunchFlag && mealsDinnerFlag {
			return fmt.Errorf("cannot use both --lunch and --dinner")
		}
		if mealsTomorrowFlag && mealsDateFlag != "" {
			return fmt.Errorf("cannot use both --tomorrow and --date")
		}

		if mealsGroupFlag == "" && mealsFilterFlag == "" {
			mealsGroupFlag = "default"
		}

		ctx := cmd.Context()

		store, err := base.NewStore(storageDirectory)
		if err != nil {
			return err
		}
		if err := syncStore(ctx, store); err != nil {
			return err
		}

		now := time.Now()

		if mealsTomorrowFlag {
			mealsDateFlag = now.Add(time.Hour * 24).Format("2006-01-02")
		} else if mealsDateFlag == "" {
			mealsDateFlag = now.Format("2006-01-02")
		}

		if mealsLunchFlag {
			mealsDaytimeFlag = "lunch"
		} else if mealsDinnerFlag {
			mealsDaytimeFlag = "dinner"
		} else if mealsDaytimeFlag == "" {
			if now.Hour() < 14 {
				mealsDaytimeFlag = "lunch"
			} else {
				mealsDaytimeFlag = "dinner"
			}
		}

		filtered := []*storage.CanteenData{}
		if mealsFilterFlag != "" {
			filtered, err = store.Filter(ctx, langFlag, mealsFilterFlag)
			if err != nil {
				return err
			}
		} else if mealsGroupFlag != "" {
			var found *group
			for _, g := range cfg.Groups {
				if g.Name == mealsGroupFlag {
					found = &g
					break
				}
			}

			if found == nil {
				return fmt.Errorf("group %s not found", mealsGroupFlag)
			}

			filtered = make([]*storage.CanteenData, 0, len(found.Refs))
			for a, b := range store.Data.Canteens {
				for _, ref := range found.Refs {
					if ref.Id == b.Id && ref.Provider == b.Provider {
						filtered = append(filtered, store.Data.Canteens[a])
					}
				}
			}
		}

		// Built index of id -> label
		canteenLabels := map[string]string{}
		for _, c := range filtered {
			canteenLabels[c.Id] = c.Label[langFlag]
		}

		eg, ctx := errgroup.WithContext(ctx)

		results := make(map[string]base.CanteenMenu, len(filtered))
		lock := sync.Mutex{}
		for _, provider := range providers {
			provider := provider
			eg.Go(func() error {
				refs := make([]base.CanteenRef, 0)
				id := provider.Id()
				for _, r := range filtered {
					if r.Provider == id {
						refs = append(refs, base.CanteenRef{
							ID:   r.Id,
							Meta: r.Meta,
						})
					}
				}
				if len(refs) == 0 {
					return nil
				}

				menu, err := provider.FetchMenus(ctx, refs, mealsDateFlag, mealsDaytimeFlag, langFlag)
				if err != nil {
					fmt.Printf("An error occurred fetching %s: %s", provider.Label(), err)
					return nil
				}

				lock.Lock()
				for _, r := range menu {
					if len(r.Meals) == 0 {
						continue
					}
					results[id+r.Canteen] = r
				}
				lock.Unlock()

				return nil
			})
		}

		if err := eg.Wait(); err != nil {
			return err
		}
		renderResult(results, canteenLabels)

		return nil
	},
}

var subtle = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#989898"}
var special = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#9Db6FF"}
var headingStyle = lipgloss.NewStyle().
	Bold(true).
	PaddingLeft(1).
	Foreground(special).
	Border(lipgloss.ThickBorder(), false, false, false, true)

var canteenNameStyle = lipgloss.NewStyle().
	Bold(true).
	PaddingLeft(1).
	Underline(true)

var descriptionStyle = lipgloss.NewStyle().
	PaddingLeft(1).
	Border(lipgloss.ThickBorder(), false, false, false, true)

var priceStyle = lipgloss.NewStyle().
	Foreground(subtle)

func formatPrice(price int64) string {
	p := int(price)
	rappen := p % 100
	chf := p / 100
	return fmt.Sprintf("%d.%02d", chf, rappen)
}

func renderResult(results map[string]base.CanteenMenu, canteenLabels map[string]string) {
	s := ""
	for _, result := range results {
		label, ok := canteenLabels[result.Canteen]
		if !ok {
			label = result.Canteen
		}
		s += fmt.Sprintf("\n%s\n\n", canteenNameStyle.Render(label))
		for i, meal := range result.Meals {
			if i != 0 {
				s += "\n"
			}
			s += fmt.Sprintf("%s\n", headingStyle.Render(meal.Label))
			for _, line := range meal.Description {
				s += fmt.Sprintf("%s\n", descriptionStyle.Render(strings.TrimSpace(line)))
			}

			p := ""
			p += priceStyle.Render(formatPrice(meal.Prices.Student))
			p += priceStyle.Render(" / ")
			p += priceStyle.Render(formatPrice(meal.Prices.Staff))
			p += priceStyle.Render(" / ")
			p += priceStyle.Render(formatPrice(meal.Prices.Extern))
			s += descriptionStyle.Render(p)
			s += "\n"
		}
	}

	fmt.Print(s)
}

func init() {
	RootCmd.AddCommand(mealsCmd)
	mealsCmd.Flags().StringVarP(&mealsDateFlag, "date", "d", "", "Date to fetch meals for (format YYYY-MM-DD)")
	mealsCmd.Flags().StringVarP(&mealsDaytimeFlag, "daytime", "t", "", "Daytime (dinner, lunch, current)")
	mealsCmd.Flags().StringVarP(&mealsFilterFlag, "filter", "f", "", "Filter to apply to the canteen name")
	mealsCmd.Flags().StringVarP(&mealsGroupFlag, "group", "g", "", "Group to fetch canteens from")
	mealsCmd.Flags().BoolVarP(&mealsDinnerFlag, "dinner", "i", false, "Fetch dinner meals")
	mealsCmd.Flags().BoolVarP(&mealsLunchFlag, "lunch", "l", false, "Fetch lunch meals")
	mealsCmd.Flags().BoolVarP(&mealsTomorrowFlag, "tomorrow", "m", false, "Fetch tomorrow's meals")
}
