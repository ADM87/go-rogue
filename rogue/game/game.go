package game

import (
	"fmt"
	"rogue/engine"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	renderer strings.Builder
	world    engine.IQuadNode
	camera   engine.ICamera
}

func NewModel() *Model {
	return &Model{
		renderer: strings.Builder{},
		world:    engine.NewQuadTree(100, 100, 2, 4),
		camera:   engine.NewCamera(0, 0, 75, 24),
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
			m.camera.Move(0, -1)

		case "down":
			m.camera.Move(0, 1)

		case "left":
			m.camera.Move(-1, 0)

		case "right":
			m.camera.Move(1, 0)
		}
	}
	return m, nil
}

func (m *Model) View() string {
	m.renderer.Reset()

	vp := m.camera.Viewport()
	w := m.world

	m.renderer.WriteString(fmt.Sprintf("# of nodes: %d\n", w.CountNodes()))

	for y := vp.Top(); y < vp.Bottom(); y++ {
		for x := vp.Left(); x < vp.Right(); x++ {
			if x == 0 && y == 0 {
				m.renderer.WriteString("@")
			} else if x == vp.GetCenterX() && y == vp.GetCenterY() {
				m.renderer.WriteString("+")
			} else if x < w.Left() || x >= w.Right() || y < w.Top() || y >= w.Bottom() {
				m.renderer.WriteString("█")
			} else {
				m.renderer.WriteString(" ")
			}
		}
		m.renderer.WriteString("\n")
	}

	m.renderer.WriteString(fmt.Sprintf("Camera: %s\n", m.camera))

	return m.renderer.String()
}
