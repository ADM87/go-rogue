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

	followPlayer bool
	renderLoS    bool
}

func NewModel() *Model {
	mdl := &Model{}
	mdl.camera = core.NewCamera(0, 0, 65, 23)
	mdl.testMap = core.NewMap(data.NewMapConfig(14, 14, 7, 7, 20, 20))
	mdl.quadTree = core.NewQuadTree(
		mdl.testMap.GetX(),
		mdl.testMap.GetY(),
		mdl.testMap.GetWidth(),
		mdl.testMap.GetHeight(),
		4, 4,
	)
	mdl.player = core.NewEntity(0, 0, mdl.moveEntity)
	mdl.followPlayer = true
	mdl.renderLoS = true
	mdl.renderer = strings.Builder{}
	return mdl
}

func (m *Model) Init() tea.Cmd {
	spawnX, spawnY := m.testMap.GetStart()
	totalObjs := 0 //(m.quadTree.GetWidth() * m.quadTree.GetHeight() / 10) - 1 // Minus 1 for the player

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

		case "f":
			m.followPlayer = !m.followPlayer

		case "l":
			m.renderLoS = !m.renderLoS

		case "up":
			if m.followPlayer {
				m.player.MoveBy(0, -1)
			} else {
				m.camera.MoveBy(0, -1)
			}

		case "down":
			if m.followPlayer {
				m.player.MoveBy(0, 1)
			} else {
				m.camera.MoveBy(0, 1)
			}

		case "left":
			if m.followPlayer {
				m.player.MoveBy(-1, 0)
			} else {
				m.camera.MoveBy(-1, 0)
			}

		case "right":
			if m.followPlayer {
				m.player.MoveBy(1, 0)
			} else {
				m.camera.MoveBy(1, 0)
			}
		}
	}
	return m, nil
}

func (m *Model) moveEntity(entity core.IEntity, x, y int) {
	if x < m.quadTree.Left() || x >= m.quadTree.Right() || y < m.quadTree.Top() || y >= m.quadTree.Bottom() {
		return
	}
	destination := core.NewRectangle(x, y, entity.GetWidth(), entity.GetHeight())
	rooms := m.testMap.GetRooms(destination)
	for _, room := range rooms {
		if room.IsWall(x, y) {
			return
		}
	}
	if others := m.quadTree.Query(destination, false); len(others) > 0 {
		for _, other := range others {
			if other.GetX() == x && other.GetY() == y {
				continue
			}
			if other.IsColliding() {
				return
			}
		}
	}
	m.quadTree.Move(entity, x, y)

	if entity == m.player {
		rooms = m.testMap.GetRooms(m.player)
		for _, room := range rooms {
			room.Visit()
		}
	}
}

func (m *Model) updateCamera() {
	if !m.followPlayer {
		return
	}
	m.camera.MoveTo(m.player.GetX(), m.player.GetY())
}

func (m *Model) View() string {
	m.updateCamera()

	viewport := m.camera.Viewport()
	totalNodes := m.quadTree.TotalNodes()
	totalObjects := m.quadTree.TotalObjects()
	objects := m.quadTree.Query(viewport, true)
	startX, startY := m.testMap.GetStart()
	endX, endY := m.testMap.GetEnd()

	m.renderer.Reset()
	m.renderer.WriteString(fmt.Sprintf("Camera: %s\n", m.camera.String()))
	m.renderer.WriteString(fmt.Sprintf("QuadTree: %s\n", m.quadTree.String()))
	m.renderer.WriteString(fmt.Sprintf("Visible Objects: %d, Total Nodes: %d, Total Objects: %d\n", len(objects), totalNodes, totalObjects))
	m.renderer.WriteString(fmt.Sprintf("Player: %s\n", m.player.String()))
	m.renderer.WriteString(fmt.Sprintf("Following Player: %t\n", m.followPlayer))
	m.renderer.WriteString(fmt.Sprintf("Start Point: (%d, %d), End Point: (%d, %d)\n", startX, startY, endX, endY))

	m.testMap.SetActiveRegion(viewport)
	for y := viewport.Top(); y < viewport.Bottom(); y++ {
		for x := viewport.Left(); x < viewport.Right(); x++ {

			mapResult := m.testMap.Render(x, y, m.player.GetX(), m.player.GetY(), m.renderLoS)
			switch mapResult {
			case data.Wall, data.OutOfBounds:
				m.renderer.WriteRune('█')
				continue
			case data.NotVisible:
				m.renderer.WriteRune('▒')
				continue
			}

			isEmpty := true
			for _, entity := range objects {
				if entity.GetX() == x && entity.GetY() == y {
					if entity == m.player {
						m.renderer.WriteRune('🯅')
					} else {
						m.renderer.WriteRune('O')
					}
					isEmpty = false
					break
				}
			}

			if isEmpty {
				if x == startX && y == startY {
					m.renderer.WriteRune('⛁')
				} else if x == endX && y == endY {
					m.renderer.WriteRune('⛃')
				} else {
					m.renderer.WriteRune(' ')
				}
			}
		}
		m.renderer.WriteRune('\n')
	}
	return m.renderer.String()
}
