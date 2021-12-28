package components

import (
	"context"
	"fmt"
	"strings"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lukasmoellerch/mensa-cli/internal/zuerichmensa"
)

type LoaderModel struct {
	Ctx context.Context

	MensaIds []int
	Lang     string
	Date     string
	Daytime  string

	loaded  bool
	results map[int]zuerichmensa.MensaMenuResponse
	err     map[int]error
}

type resultMsg struct {
	results map[int]zuerichmensa.MensaMenuResponse
}

type errMsg struct {
	errs map[int]error
}

func (e errMsg) Error() string { return "Error loading menuplan" }

func (m LoaderModel) Init() tea.Cmd {
	return func() tea.Msg { return checkServer(m) }
}

func (m LoaderModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		default:
			return m, nil
		}

	case resultMsg:
		m.results = msg.results
		m.loaded = true
		return m, tea.Quit

	case errMsg:
		m.err = msg.errs
		return m, nil

	default:
		return m, nil
	}
}

var subtle = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#989898"}
var special = lipgloss.AdaptiveColor{Light: "#874BFD", Dark: "#9Db6FF"}
var headingStyle = lipgloss.NewStyle().
	Bold(true).
	Underline(true).
	PaddingLeft(1).
	Foreground(special).
	Border(lipgloss.ThickBorder(), false, false, false, true)

var descriptionStyle = lipgloss.NewStyle().
	PaddingLeft(1).
	Border(lipgloss.ThickBorder(), false, false, false, true)

var priceStyle = lipgloss.NewStyle().
	Foreground(subtle)

func (m LoaderModel) View() string {
	if m.err != nil {
		return "something went wrong"
	} else if m.loaded {
		s := ""
		for _, result := range m.results {
			s += fmt.Sprintf("> %s\n\n", result.Mensa)
			for i, meal := range result.Menu.Meals {
				if i != 0 {
					s += "\n"
				}
				s += fmt.Sprintf("%s\n", headingStyle.Render(meal.Label))
				for _, line := range meal.Description {
					s += fmt.Sprintf("%s\n", descriptionStyle.Render(strings.TrimSpace(line)))
				}

				p := ""
				p += priceStyle.Render(meal.Prices.Student)
				p += priceStyle.Render(" / ")
				p += priceStyle.Render(meal.Prices.Staff)
				p += priceStyle.Render(" / ")
				p += priceStyle.Render(meal.Prices.Extern)
				s += descriptionStyle.Render(p)
				s += "\n"
			}
		}
		return s
	} else {
		return "Loading...\n"
	}
}

func checkServer(m LoaderModel) tea.Msg {
	wg := &sync.WaitGroup{}
	wg.Add(len(m.MensaIds))
	results := make(map[int]zuerichmensa.MensaMenuResponse, len(m.MensaIds))
	errs := make(map[int]error, len(m.MensaIds))
	lock := &sync.Mutex{}
	for _, id := range m.MensaIds {
		go func(id int) {
			res, err := zuerichmensa.FetchMenuEth(
				m.Ctx,
				id,
				m.Lang,
				m.Date,
				m.Daytime,
			)
			lock.Lock()
			if err != nil {
				errs[id] = err
			} else {
				results[id] = res
			}
			lock.Unlock()

			wg.Done()
		}(id)
	}

	wg.Wait()

	if len(errs) != 0 {
		return errMsg{errs}
	}
	return resultMsg{results}
}
