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
	player   core.IEntity
	renderer strings.Builder
}

func NewModel() *Model {
	return &Model{
		camera:   core.NewCamera(0, 0, 75, 25),
		quadTree: core.NewQuadTree(0, 0, 400, 200, 5, 4),
		player:   core.NewEntity(0, 0),
		renderer: strings.Builder{},
	}
}

func (m *Model) Init() tea.Cmd {
	positions := make(map[string]bool)
	positions[fmt.Sprintf("%d,%d", m.player.GetX(), m.player.GetY())] = true
	for i := 0; i < 600; i++ {
		for {
			x := rand.Intn(m.quadTree.GetWidth())
			y := rand.Intn(m.quadTree.GetHeight())
			key := fmt.Sprintf("%d,%d", x, y)
			if _, ok := positions[key]; !ok {
				positions[key] = true
				m.quadTree.Insert(core.NewEntity(x, y))
				break
			}
		}
	}
	m.quadTree.Insert(m.player)
	m.camera.MoveTo(m.player.GetX(), m.player.GetY())
	m.camera.ClampToBounds(m.quadTree)
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
			m.MovePlayer(0, -1)

		case "down":
			m.MovePlayer(0, 1)

		case "left":
			m.MovePlayer(-1, 0)

		case "right":
			m.MovePlayer(1, 0)
		}
	}
	return m, nil
}

func (m *Model) MovePlayer(dx, dy int) {
	x, y := m.player.GetX()+dx, m.player.GetY()+dy
	if x < m.quadTree.Left() || x >= m.quadTree.Right() || y < m.quadTree.Top() || y >= m.quadTree.Bottom() {
		return
	}
	collisionCheck := m.quadTree.Query(core.NewRectangle(x, y, 1, 1), true)
	if len(collisionCheck) > 0 {
		return
	}
	m.quadTree.Move(m.player, x, y)
}

func (m *Model) View() string {
	m.camera.MoveTo(m.player.GetX(), m.player.GetY())
	m.camera.ClampToBounds(m.quadTree)

	viewport := m.camera.Viewport()
	isColliding := m.quadTree.Collides(viewport)
	isOverlapping := m.quadTree.Overlaps(viewport)
	isContaining := m.quadTree.Contains(m.camera.GetX(), m.camera.GetY())
	isBorder := m.quadTree.IsBorder(m.camera.GetX(), m.camera.GetY())
	totalNodes := m.quadTree.TotalNodes()
	totalObjects := m.quadTree.TotalObjects()
	objects := m.quadTree.Query(viewport, true)

	m.renderer.Reset()
	m.renderer.WriteString(fmt.Sprintf("Camera: %s\n", m.camera.String()))
	m.renderer.WriteString(fmt.Sprintf("QuadTree: %s\n", m.quadTree.String()))
	m.renderer.WriteString(fmt.Sprintf("Colliding: %t, Overlapping: %t, Containing: %t, OnBorder: %t\n", isColliding, isOverlapping, isContaining, isBorder))
	m.renderer.WriteString(fmt.Sprintf("Visible Objects: %d, Total Nodes: %d, Total Objects: %d\n", len(objects), totalNodes, totalObjects))
	m.renderer.WriteString(fmt.Sprintf("Player on border: %t\n", m.quadTree.IsBorder(m.player.GetX(), m.player.GetY())))

	for y := viewport.Top(); y < viewport.Bottom(); y++ {
		for x := viewport.Left(); x < viewport.Right(); x++ {
			if len(objects) > 0 {
				found := false
				for _, o := range objects {
					if o.GetX() == x && o.GetY() == y {
						if o.GetX() == m.player.GetX() && o.GetY() == m.player.GetY() {
							m.renderer.WriteRune('☺')
						} else {
							m.renderer.WriteRune('•')
						}
						found = true
						break
					}
				}
				if found {
					continue
				}
			}
			if m.quadTree.IsBorder(x, y) {
				m.renderer.WriteRune('▒')
				continue
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
