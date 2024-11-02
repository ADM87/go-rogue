package game

import (
	"rogue/engine"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	_camera *engine.Camera
}

func NewModel() *Model {
	return &Model{
		_camera: engine.NewCamera(0, 0, 60, 20),
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		switch key {
		case "esc":
			return m, tea.Quit

		case "up":
			m._camera.MoveY(-1)

		case "down":
			m._camera.MoveY(1)

		case "left":
			m._camera.MoveX(-1)

		case "right":
			m._camera.MoveX(1)
		}
	}
	return m, nil
}

func (m *Model) View() string {
	output := ""
	output += m._camera.DrawDebug()

	return output
}
