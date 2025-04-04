package main

import (
	"log"

	"Neolithic/internal/agent"
	"Neolithic/internal/camera"
	"Neolithic/internal/core"
	"Neolithic/internal/grid"
	"Neolithic/internal/logging"
	"Neolithic/internal/planner"
	"Neolithic/internal/world"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	Engine   *world.Engine
	Camera   *camera.Camera
	Viewport *camera.Viewport
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.Camera.Move(0, -10)
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.Camera.Move(0, 10)
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.Camera.Move(10, 0)
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.Camera.Move(-10, 0)
	}

	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		g.Camera.ZoomAt(1.05, float64(g.Viewport.Width), float64(g.Viewport.Height))
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		g.Camera.ZoomAt(.95, float64(g.Viewport.Width), float64(g.Viewport.Height))
	}

	if err := g.Engine.Tick(1.0 / 60); err != nil {
		return err
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.Engine.Draw(screen, g.Viewport, g.Camera)
}

func (g *Game) Layout(_, _ int) (screenWidth, screenHeight int) {
	return g.Viewport.Width, g.Viewport.Height
}

func main() {

	cam := camera.NewCamera()
	vp := camera.NewViewport(cam, 800, 600)
	width, height := 32, 32

	logger := logging.NewLogger("info")

	worldGrid, err := grid.New(width, height, 16)
	if err != nil {
		log.Fatal(err)
	}
	if err = worldGrid.Initialize(world.MakeTile); err != nil {
		log.Fatal(err)
	}

	engine, err := world.NewEngine(worldGrid, logger)
	if err != nil {
		log.Fatal(err)
	}

	game := &Game{
		Engine:   engine,
		Camera:   cam,
		Viewport: vp,
	}

	loc1 := core.Location{
		Name:      "loc1",
		Inventory: core.NewInventory(),
		Coord:     core.Coord{X: 3, Y: 14},
	}

	loc2 := core.Location{
		Name:      "loc2",
		Inventory: core.NewInventory(),
		Coord:     core.Coord{X: 21, Y: 4},
	}

	res1 := &core.Resource{Name: "Berries"}

	loc1.Inventory.AdjustAmount(res1, 100)

	goalLoc2 := loc2.DeepCopy()
	goalLoc2.Inventory.AdjustAmount(res1, 100)

	depositLoc1 := &planner.Deposit{
		Resource:       res1,
		Amount:         10,
		ActionLocation: &loc1,
		ActionCost:     1,
	}

	depositLoc2 := &planner.Deposit{
		Resource:       res1,
		Amount:         10,
		ActionLocation: &loc2,
		ActionCost:     1,
	}

	gatherLoc1 := &planner.Gather{
		Resource:       res1,
		Amount:         10,
		ActionLocation: &loc1,
		ActionCost:     1,
	}

	gatherLoc2 := &planner.Gather{
		Resource:       res1,
		Amount:         10,
		ActionLocation: &loc2,
		ActionCost:     1,
	}

	goalState := &core.WorldState{
		Locations: map[string]core.Location{
			goalLoc2.Name: *goalLoc2,
		},
	}

	testAgent := agent.NewAgent("agent", logger)
	testAgent.Behavior.Goal = goalState
	testAgent.Behavior.PossibleActions = []planner.Action{
		depositLoc1,
		depositLoc2,
		gatherLoc1,
		gatherLoc2,
	}

	gameWorld := engine.World
	gameWorld.Agents[testAgent.Name()] = testAgent
	gameWorld.Locations[loc1.Name] = loc1
	gameWorld.Locations[loc2.Name] = loc2

	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
