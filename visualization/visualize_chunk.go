package visualization

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	modelStyle = lipgloss.NewStyle().
			Width(15).
			Height(5).
			Align(lipgloss.Center, lipgloss.Center).
			BorderStyle(lipgloss.HiddenBorder())
	focusedModelStyle = lipgloss.NewStyle().
				Width(15).
				Height(5).
				Align(lipgloss.Center, lipgloss.Center).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("69"))
	spinnerStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
	helpStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
)

type chunkViewerModel struct {
	index int
}

func (m chunkViewerModel) Init() tea.Cmd {
	return nil
}

func (m chunkViewerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
		switch msg.Type {
		case tea.KeyRight:
			// pass
		}
	}
	return m, nil
}

func (m chunkViewerModel) View() string {
	var s strings.Builder
	if m.index == 0 {
		s.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, focusedModelStyle.Render(fmt.Sprintf("%4s", "No")), modelStyle.Render("NO")))
	} else {
		s.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, modelStyle.Render(fmt.Sprintf("%4s", "Yes")), focusedModelStyle.Render("YES")))
	}
	return s.String()
}

func initChunkViewer() chunkViewerModel {
	var model chunkViewerModel

	return model
}

var VisualizeChunkCommand = &cobra.Command{
	Use:   "chunk",
	Short: "Visualize a single chunk",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		prog := tea.NewProgram(initChunkViewer())

		if _, err := prog.Run(); err != nil {
			return err
		}

		return nil
	},
}
