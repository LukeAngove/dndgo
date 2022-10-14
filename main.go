package main

import (
  "fmt"
  "io"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/table"
  ftable "github.com/lukeangove/bubbles-interface-fixes/bubbles/table"
	"github.com/charmbracelet/bubbles/list"
  flist "github.com/lukeangove/bubbles-interface-fixes/bubbles/list"

  "github.com/charmbracelet/lipgloss"
)

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4).Width(80).Render
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170")).Render
)

func makeAttributes() tea.Model {
  rows := []table.Row{
    {"STR", "18", "+4"},
    {"DEX", "18", "+4"},
    {"CON", "18", "+4"},
    {"INT", "18", "+4"},
    {"WIS", "18", "+4"},
    {"CHA", "18", "+4"},
  }

	return ftable.New(table.New(
    table.WithColumns([]table.Column{
      {Title: "Attr", Width: 4},
      {Title: "Val", Width: 3},
      {Title: "Mod", Width: 3},
    }),
    table.WithRows(rows),
    table.WithHeight(len(rows)),
  ))
}

func makeSkills() tea.Model {
  rows := []table.Row{
    {"acrobatics", "Y", "2"},
    {"athletics", "N", "0"},
    {"history", "H", "1"},
    {"arcana", "M", "4"},
  }

	return ftable.New(table.New(
    table.WithColumns([]table.Column{
      {Title: "Skill", Width: 12},
      {Title: "Prof", Width: 4},
      {Title: "Val", Width: 3},
    }),
    table.WithRows(rows),
    table.WithHeight(len(rows)),
  ))
}

func makeItems() tea.Model {
  rows := []table.Row{
    {"Y", "sword", "4"},
  }

	return ftable.New(table.New(
    table.WithColumns([]table.Column{
      {Title: "Eq", Width: 2},
      {Title: "Name", Width: 12},
      {Title: "Wgt", Width: 3},
    }),
    table.WithRows(rows),
    table.WithHeight(len(rows)),
  ))
}

func makeAttacks() tea.Model {
  rows := []table.Row{
    {"sword", "+6", "1d6+4"},
  }

	return ftable.New(table.New(
    table.WithColumns([]table.Column{
      {Title: "Name", Width: 12},
      {Title: "Atk", Width: 3},
      {Title: "Damage", Width: 8},
    }),
    table.WithRows(rows),
    table.WithHeight(len(rows)),
  ))
}

type item string

func (i item) FilterValue() string { return "" }

type itemDelegate struct{}

func (d itemDelegate) Height() int                               { return 1 }
func (d itemDelegate) Spacing() int                              { return 0 }
func (d itemDelegate) Update(msg tea.Msg, m *list.Model) tea.Cmd { return nil }
func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := fmt.Sprintf("%d. %s", index+1, i)

	fn := itemStyle
	if index == m.Index() {
		fn = func(s string) string {
			return selectedItemStyle("> " + s)
		}
	}

	fmt.Fprint(w, fn(str))
}

func makeAbilities() tea.Model {
  abilities := []list.Item{
    item("hi"),
    item("hello"),
  }
  return flist.New(list.New(abilities, itemDelegate{}, 20, 10))
}

func main() {
  windows := []NamedModel{
    NamedModel{"a", "attributes", makeAttributes()},
    NamedModel{"s", "skills", makeSkills()},
    NamedModel{"i", "inventory", makeItems()},
    NamedModel{"t", "attack", makeAttacks()},
    NamedModel{"b", "abilities", makeAbilities()},
  }

	p := tea.NewProgram(New(windows, 50))

	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}

