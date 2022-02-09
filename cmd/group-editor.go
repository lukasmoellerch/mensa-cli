package cmd

import (
	"context"
	"fmt"
	"io"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/lukasmoellerch/mensa-cli/internal/base"
)

const listHeight = 24

var (
	titleStyle              = lipgloss.NewStyle().MarginLeft(2)
	itemStyle               = lipgloss.NewStyle().PaddingLeft(2)
	selectedItemStyle       = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
	activeItemStyle         = lipgloss.NewStyle().PaddingLeft(2).Bold(true)
	selectedActiveItemStyle = lipgloss.NewStyle().PaddingLeft(2).Bold(true).Foreground(lipgloss.Color("170"))
	paginationStyle         = list.DefaultStyles().PaginationStyle.PaddingLeft(4)
	helpStyle               = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)
	quitTextStyle           = lipgloss.NewStyle().Margin(1, 0, 2, 4)
)

type item struct {
	index    int
	selected bool
	label    string
	provider string
	id       string
}

func (i item) FilterValue() string { return i.label }

type itemDelegate struct{}

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	var fn func(s string) string
	if index == m.Index() {
		if i.selected {
			fn = func(s string) string {
				return selectedActiveItemStyle.Render("> [X] " + s)
			}
		} else {
			fn = func(s string) string {
				return selectedItemStyle.Render("> [ ] " + s)
			}
		}
	} else {
		if i.selected {
			fn = func(s string) string {
				return activeItemStyle.Render("  [X] " + s)
			}
		} else {
			fn = func(s string) string {
				return itemStyle.Render("  [ ] " + s)
			}
		}
	}
	fmt.Fprint(w, fn("("+i.provider+") "+i.label))
}

type groupEditorModel struct {
	list     list.Model
	selected map[int]struct{}
	aborting *bool
	quitting *bool
}

func (m groupEditorModel) Init() tea.Cmd {
	return nil
}

func (m groupEditorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil

	case tea.KeyMsg:
		keypress := msg.String()
		if keypress == "ctrl+c" || keypress == "q" {
			*m.aborting = true
			return m, tea.Quit
		}
		if keypress == "enter" && !m.list.SettingFilter() {
			*m.quitting = true
			return m, tea.Quit
		}
		if keypress == " " && !m.list.SettingFilter() {
			i := m.list.SelectedItem().(item)
			i.selected = !i.selected
			cmd := m.list.SetItem(i.index, i)
			if i.selected {
				m.selected[i.index] = struct{}{}
			} else {
				delete(m.selected, i.index)
			}
			m.list.Title = title(len(m.selected))
			return m, cmd
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m groupEditorModel) View() string {
	if *m.aborting {
		return quitTextStyle.Render("Edit aborted.")
	}
	if *m.quitting {
		return quitTextStyle.Render("Saved.")
	}
	return "\n" + m.list.View()
}

func title(num int) string {
	return fmt.Sprintf("%d canteens selected", num)
}

func startGroupEditor(ctx context.Context, store *base.Store, initialRefs []canteenRef) ([]canteenRef, error) {
	idMap := make(map[string]map[string]struct{}, len(providers))
	for _, c := range initialRefs {
		if _, ok := idMap[c.Provider]; !ok {
			idMap[c.Provider] = make(map[string]struct{})
		}
		idMap[c.Provider][c.Id] = struct{}{}
	}

	selectedMap := make(map[int]struct{})
	itemList := make([]list.Item, len(store.Data.Canteens))
	for i, c := range store.Data.Canteens {
		selected := false
		if prov, ok := idMap[c.Provider]; ok {
			if _, ok := prov[c.Id]; ok {
				selected = true
			}
		}
		if selected {
			selectedMap[i] = struct{}{}
		}
		itemList[i] = item{
			index:    i,
			selected: selected,
			label:    c.Label[langFlag],
			provider: c.Provider,
			id:       c.Id,
		}
	}

	const defaultWidth = 20

	l := list.New(itemList, itemDelegate{}, defaultWidth, listHeight)
	l.Title = title(len(initialRefs))
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)
	l.Styles.Title = titleStyle
	l.Styles.PaginationStyle = paginationStyle
	l.Styles.HelpStyle = helpStyle

	aborted := false
	quit := false
	m := groupEditorModel{
		list:     l,
		selected: selectedMap,
		aborting: &aborted,
		quitting: &quit,
	}

	if err := tea.NewProgram(m).Start(); err != nil {
		return nil, err
	}

	if aborted {
		return nil, fmt.Errorf("aborted")
	}

	refs := make([]canteenRef, 0, len(m.selected))
	for i := range m.selected {
		item := itemList[i].(item)
		refs = append(refs, canteenRef{
			Provider: item.provider,
			Id:       item.id,
		})
	}
	return refs, nil
}
