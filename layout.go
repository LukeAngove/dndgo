package main

import (
  //"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var baseStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240")).
	Render

type NamedModel struct {
  key   string
  label string
  model tea.Model
}

type Model struct {
	focused         int
  minColumnWidth  int
  width           int
  height          int
	windows         []NamedModel
	keys            map[string]int
}

func New(windows []NamedModel, minColumnWidth int) Model {
  keys := make(map[string]int)
  keys["self"] = len(windows)

  for i, w := range windows {
    keys[w.key] = i
  }

	return Model{
		focused: keys["self"],
    minColumnWidth: minColumnWidth,
    width: 0,
    height: 0,
		windows: windows,
    keys: keys,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) isSelfFocused() bool {
  return m.focused == m.keys["self"]
}

func (m Model) updateChild(msg tea.Msg) tea.Cmd {
  updatedModel, cmd := m.windows[m.focused].model.Update(msg)
  m.windows[m.focused] = NamedModel{m.windows[m.focused].key, m.windows[m.focused].label, updatedModel}
  return cmd
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
    switch msg.String() {
    case "q", "esc", "ctrl+c":
      if m.isSelfFocused() {
        return m, tea.Quit
      } else {
        m.focused = m.keys["self"]
        return m, nil
      }
    default:
      if m.isSelfFocused() {
        if idx, ok := m.keys[msg.String()]; ok {
          m.focused = idx
        }
        return m, nil
      } else {
        cmd = m.updateChild(msg)
        return m, cmd
      }
    }
  case tea.WindowSizeMsg:
    m.width = msg.Width
    m.height = msg.Height
  default:
    cmd = m.updateChild(msg)
  }

	return m, cmd
}

func (m Model) View() string {
  ncolumns := m.width / m.minColumnWidth
  if ncolumns == 0 {
    return "Screen too narrow, not rendering"
  }

  columnWidth := m.width / ncolumns

  //return fmt.Sprintf("Width: %d, Height: %d, columnWidth: %d, columns: %d\n", m.width, m.height, columnWidth, ncolumns)

  var normalStyle = lipgloss.NewStyle().
        Padding(1, 2, 0, 2).
        Border(lipgloss.RoundedBorder()).
        Width(columnWidth).
        Align(lipgloss.Center).
        Render


  var strings []string
  var heights []int

  for _, w := range m.windows {
    data := normalStyle(w.model.View())
    strings = append(strings, data)
    heights = append(heights, lipgloss.Height(data))
  }

  layout := balance(heights, m.height, ncolumns)

  var columns []string
  for _, column := range layout {
    var column_string []string

    for _, cell := range column {
      column_string = append(column_string, strings[cell])
    }
    columns = append(columns, lipgloss.JoinVertical(lipgloss.Left, column_string...))
  }

  return lipgloss.JoinHorizontal(lipgloss.Top, columns...)
}

func sum(vals []int) int {
  sum := 0

  for _, v := range vals {
    sum += v
  }

  return sum
}

func abs_diff(a int, b int) int {
  if a > b {
    return a - b
  } else {
    return b - a
  }
}

func balance(heights []int, max_height int, columns int) [][]int {
  column_sizes := make([]int, columns)
  for i, _ := range column_sizes {
    column_sizes[i] = 0
  }

  column_idx := make([][]int, columns)
  for i, _ := range column_idx {
    column_idx[i] = make([]int, 0)
  }

  ideal := sum(heights) / columns

  current_col := 0
  for i, v := range heights {
    if column_sizes[current_col] + v > ideal &&
       current_col < (columns-1) &&
       len(column_idx[current_col]) > 0 {

      columns_left := columns-(current_col+1)
      sum_left := sum(heights[i:])
      higher_expected := sum_left / columns_left
      lower_expected := (sum_left-heights[i]) / columns_left

      if (abs_diff(higher_expected, ideal) < abs_diff(lower_expected, ideal)) ||
          column_sizes[current_col] + v > max_height {
        current_col++
      }
    }
    column_sizes[current_col] += v
    column_idx[current_col] = append(column_idx[current_col], i)
  }

  return column_idx
}

