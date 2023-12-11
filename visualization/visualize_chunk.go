package visualization

import (
	"errors"
	"strings"

	"git.nicholasnovak.io/nnovak/spatial-db/server"
	"git.nicholasnovak.io/nnovak/spatial-db/storage"
	"git.nicholasnovak.io/nnovak/spatial-db/world"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

const (
	chunkWidth  = 15
	chunkHeight = 5
)

var (
	missingChunkStyle = lipgloss.NewStyle().
				Width(chunkWidth).
				Height(chunkHeight).
				Align(lipgloss.Center, lipgloss.Center).
				BorderStyle(lipgloss.HiddenBorder()).
				BorderForeground(lipgloss.Color("69"))
	presentChunkStyle = lipgloss.NewStyle().
				Width(chunkWidth).
				Height(chunkHeight).
				Align(lipgloss.Center, lipgloss.Center).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("69"))
	selectedChunkStyle = lipgloss.NewStyle().
				Width(chunkWidth).
				Height(chunkHeight).
				Align(lipgloss.Center, lipgloss.Center).
				BorderStyle(lipgloss.DoubleBorder()).
				BorderForeground(lipgloss.Color("202"))

	loadedChunkCache = make(map[world.ChunkPos]bool)
)

type chunkViewerModel struct {
	chunkServer *server.SimpleServer

	visibleChunkRows int
	visibleChunkCols int

	currentPos world.ChunkPos
}

func (m *chunkViewerModel) updateShownChunks(newWidth, newHeight int) {
	horiz, vert := missingChunkStyle.GetFrameSize()
	m.visibleChunkRows = newHeight / (chunkHeight + vert)
	m.visibleChunkCols = newWidth / (chunkWidth + horiz)
}

func (m chunkViewerModel) Init() tea.Cmd {
	return tea.EnterAltScreen
}

func (m chunkViewerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		}
		switch msg.Type {
		case tea.KeyRight:
			m.currentPos.X += 1
		case tea.KeyLeft:
			m.currentPos.X -= 1
		case tea.KeyUp:
			m.currentPos.Z -= 1
		case tea.KeyDown:
			m.currentPos.Z += 1
		}
	case tea.WindowSizeMsg:
		m.updateShownChunks(msg.Width, msg.Height)
	}
	return m, nil
}

func (m chunkViewerModel) View() string {
	var s strings.Builder

	midRow := m.visibleChunkRows / 2
	midCol := m.visibleChunkCols / 2

	for rowIndex := 0; rowIndex < m.visibleChunkRows; rowIndex++ {
		renderedRow := make([]string, m.visibleChunkCols)
		for colIndex := 0; colIndex < m.visibleChunkCols; colIndex++ {
			currentChunkPos := world.ChunkPos{
				X: midCol - colIndex - m.currentPos.X,
				Z: midRow - rowIndex - m.currentPos.Z,
			}

			var fetchChunkErr error
			if isPresent, cached := loadedChunkCache[currentChunkPos]; cached {
				if isPresent {
					fetchChunkErr = nil
				} else {
					fetchChunkErr = storage.ChunkNotFoundError
				}
			} else {
				_, fetchChunkErr = m.chunkServer.FetchChunk(currentChunkPos)
				loadedChunkCache[currentChunkPos] = fetchChunkErr == nil
			}

			chunkDisplay := currentChunkPos.StringCoords()

			if rowIndex == midRow && colIndex == midCol {
				renderedRow[colIndex] = selectedChunkStyle.Render(chunkDisplay)
			} else {
				if fetchChunkErr == nil {
					renderedRow[colIndex] = presentChunkStyle.Render(chunkDisplay)
				} else if errors.Is(fetchChunkErr, storage.ChunkNotFoundError) {
					renderedRow[colIndex] = missingChunkStyle.Render(chunkDisplay)
				} else {
					panic(fetchChunkErr)
				}
			}
		}
		s.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, renderedRow...) + "\n")
	}

	return s.String()
}

func initChunkViewer(chunkServer *server.SimpleServer) chunkViewerModel {
	var model chunkViewerModel

	model.chunkServer = chunkServer

	return model
}

var VisualizeChunkCommand = &cobra.Command{
	Use:   "chunk",
	Short: "Visualize a single chunk",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {

		// Create a new server to read from those files
		var chunkServer server.SimpleServer
		chunkServer.StorageDir = args[0]

		prog := tea.NewProgram(initChunkViewer(&chunkServer), tea.WithAltScreen())

		if _, err := prog.Run(); err != nil {
			return err
		}

		return nil
	},
}
