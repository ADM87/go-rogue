package game

import (
	"fmt"
	"math/rand"
	"rogue/core"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	camera   core.ICamera
	quadTree core.IQuadNode
	renderer strings.Builder
}

func NewModel() *Model {
	return &Model{
		camera:   core.NewCamera(0, 0, 65, 15),
		quadTree: core.NewQuadTree(0, 0, 65, 15, 2, 4),
		renderer: strings.Builder{},
	}
}

func (m *Model) Init() tea.Cmd {
	positions := make(map[string]bool)
	for i := 0; i < 975; i++ {
		for {
			x := rand.Intn(m.quadTree.GetWidth())
			y := rand.Intn(m.quadTree.GetHeight())
			key := fmt.Sprintf("%d,%d", x, y)
			if _, ok := positions[key]; !ok {
				positions[key] = true
				m.quadTree.Insert(core.NewPoint(x, y))
				break
			}
		}
	}
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
			m.camera.MoveBy(0, -1)

		case "down":
			m.camera.MoveBy(0, 1)

		case "left":
			m.camera.MoveBy(-1, 0)

		case "right":
			m.camera.MoveBy(1, 0)
		}
	}
	return m, nil
}

func (m *Model) View() string {
	viewport := m.camera.Viewport()
	isColliding := m.quadTree.Collides(viewport)
	isOverlapping := m.quadTree.Overlaps(viewport)
	isContaining := m.quadTree.Contains(m.camera.GetX(), m.camera.GetY())
	objects := m.quadTree.Query(viewport, true)

	m.renderer.Reset()
	m.renderer.WriteString(fmt.Sprintf("Camera: %s\n", m.camera.String()))
	m.renderer.WriteString(fmt.Sprintf("QuadTree: %s\n", m.quadTree.String()))
	m.renderer.WriteString(fmt.Sprintf("Colliding: %t, Overlapping: %t, Containing: %t\n", isColliding, isOverlapping, isContaining))
	m.renderer.WriteString(fmt.Sprintf("# Objects: %d\n", len(objects)))

	for y := viewport.Top(); y < viewport.Bottom(); y++ {
		for x := viewport.Left(); x < viewport.Right(); x++ {
			if x == m.camera.GetX() && y == m.camera.GetY() {
				m.renderer.WriteRune('+')
				continue
			}
			if len(objects) > 0 {
				found := false
				for _, o := range objects {
					if o.GetX() == x && o.GetY() == y {
						m.renderer.WriteRune('o')
						found = true
						break
					}
				}
				if found {
					continue
				}
			}
			if !m.quadTree.Contains(x, y) {
				m.renderer.WriteRune('█')
				continue
			}
			m.renderer.WriteRune(' ')
		}
		m.renderer.WriteRune('\n')
	}
	return m.renderer.String()
}
