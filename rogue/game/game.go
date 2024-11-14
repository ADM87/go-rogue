package game

import (
	"fmt"
	"math/rand"
	"rogue/core"
	"rogue/data"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	camera   core.ICamera
	quadTree core.IQuadNode
	player   core.IEntity
	testMap  core.IMap
	renderer strings.Builder
}

func NewModel() *Model {
	mdl := &Model{}
	mdl.camera = core.NewCamera(0, 0, 75, 25)
	mdl.testMap = core.NewMap(data.NewMapConfig(20, 30, 10, 15, 5, 10))
	mdl.quadTree = core.NewQuadTree(
		mdl.testMap.GetX(),
		mdl.testMap.GetY(),
		mdl.testMap.GetWidth(),
		mdl.testMap.GetHeight(),
		4, 4,
	)
	mdl.player = core.NewEntity(0, 0, mdl.moveEntity)
	mdl.renderer = strings.Builder{}
	return mdl
}

func (m *Model) Init() tea.Cmd {
	spawnX, spawnY := m.testMap.GetStart()
	totalObjs := (m.quadTree.GetWidth() * m.quadTree.GetHeight() / 10) - 1 // Minus 1 for the player

	m.player.SetXY(spawnX, spawnY)
	m.quadTree.Insert(m.player)
	m.updateCamera()

	positions := make(map[string]bool)
	positions[fmt.Sprintf("%d,%d", m.player.GetX(), m.player.GetY())] = true
	for i := 0; i < totalObjs; i++ {
		for {
			x := m.testMap.GetX() + rand.Intn(m.testMap.GetWidth())
			y := m.testMap.GetY() + rand.Intn(m.testMap.GetHeight())
			key := fmt.Sprintf("%d,%d", x, y)
			if _, ok := positions[key]; !ok && !m.quadTree.IsBorder(x, y) {
				positions[key] = true
				m.quadTree.Insert(core.NewEntity(x, y, m.moveEntity))
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
			m.player.MoveBy(0, -1)

		case "down":
			m.player.MoveBy(0, 1)

		case "left":
			m.player.MoveBy(-1, 0)

		case "right":
			m.player.MoveBy(1, 0)
		}
	}
	return m, nil
}

func (m *Model) moveEntity(entity core.IEntity, x, y int) {
	if x < m.quadTree.Left() || x >= m.quadTree.Right() || y < m.quadTree.Top() || y >= m.quadTree.Bottom() {
		return
	}
	if others := m.quadTree.Query(core.NewRectangle(x, y, 1, 1), false); len(others) > 0 {
		for _, other := range others {
			if other != entity && other.GetX() == x && other.GetY() == y {
				entity.OnCollisionStart(other)
				if entity == m.player {
					dx, dy := x-entity.GetX(), y-entity.GetY()
					other.MoveBy(dx, dy)
					entity.OnCollisionEnd()
				}
				return
			}
		}
		if entity.IsColliding() {
			entity.OnCollisionEnd()
		}
	}
	m.quadTree.Move(entity, x, y)
}

func (m *Model) updateCamera() {
	m.camera.MoveTo(m.player.GetX(), m.player.GetY())
	// m.camera.ClampToBounds(m.quadTree)
}

func (m *Model) View() string {
	m.updateCamera()

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
	m.renderer.WriteString(fmt.Sprintf("Player: %s\n", m.player.String()))
	m.renderer.WriteString(fmt.Sprintf("Player Colliding: %t\n", m.player.IsColliding()))

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
			// if m.quadTree.IsBorder(x, y) {
			// 	m.renderer.WriteRune('▒')
			// 	continue
			// }
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
