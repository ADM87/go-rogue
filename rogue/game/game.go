package game

import (
	"fmt"
	"rogue/engine"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	_renderer strings.Builder
	_camera   engine.ICamera
}

func NewModel() *Model {
	return &Model{
		_renderer: strings.Builder{},
		_camera:   engine.NewCamera(0, 0, 75, 24),
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
			m._camera.Move(0, -1)

		case "down":
			m._camera.Move(0, 1)

		case "left":
			m._camera.Move(-1, 0)

		case "right":
			m._camera.Move(1, 0)
		}
	}
	return m, nil
}

func (m *Model) View() string {
	m._renderer.Reset()

	vp := m._camera.Viewport()
	for y := vp.Top(); y < vp.Bottom(); y++ {
		for x := vp.Left(); x < vp.Right(); x++ {
			if x == 0 && y == 0 {
				m._renderer.WriteString("@")
			} else if x == vp.GetCenterX() && y == vp.GetCenterY() {
				m._renderer.WriteString("+")
			} else {
				m._renderer.WriteString(".")
			}
		}
		m._renderer.WriteString("\n")
	}

	m._renderer.WriteString(fmt.Sprintf("Camera: %s\n", m._camera))

	return m._renderer.String()
}
