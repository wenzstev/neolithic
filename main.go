package main

import (
	"Neolithic/internal/goalengine"
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

	loc1 := &core.Location{
		Name:      "loc1",
		Inventory: core.NewInventory(),
		Coord:     core.Coord{X: 3, Y: 14},
	}

	loc2 := &core.Location{
		Name:      "loc2",
		Inventory: core.NewInventory(),
		Coord:     core.Coord{X: 21, Y: 4},
	}

	loc3 := &core.Location{
		Name:      "loc3",
		Inventory: core.NewInventory(),
		Coord:     core.Coord{X: 27, Y: 30},
	}

	depo := &core.Location{
		Name:      "depo",
		Inventory: core.NewInventory(),
		Coord:     core.Coord{X: 16, Y: 16},
	}

	res1 := &core.Resource{Name: "Berries"}

	loc1.Inventory.AdjustAmount(res1, 20)
	loc2.Inventory.AdjustAmount(res1, 10)
	loc3.Inventory.AdjustAmount(res1, 20)

	goalDepo := depo.DeepCopy()
	goalDepo.Inventory.AdjustAmount(res1, 50)

	testAgent := agent.NewAgent("agent", logger)
	testAgent.Behavior.GoalEngine = &goalengine.GoalEngine{
		Goal: goalengine.Goal{
			Name: "gather berries",
			Logic: goalengine.GoalLogic{
				Chunker:      goalengine.AddToLocation,
				Fallback:     goalengine.FallbackChunkFunc,
				ShouldGiveUp: goalengine.GiveUpIfNoChange,
			},
			Location: goalDepo,
			Resource: res1,
		},
	}

	createDepositAction := func(params world.ActionCreatorParams) planner.Action {
		return &planner.Deposit{
			DepResource:    params.Resource,
			Amount:         2,
			ActionLocation: params.Location,
			ActionCost:     1,
		}
	}

	createGatherAction := func(params world.ActionCreatorParams) planner.Action {
		return &planner.Gather{
			Res:            params.Resource,
			Amount:         2,
			ActionLocation: params.Location,
			ActionCost:     1,
		}
	}

	if err = engine.RegisterAction("deposit", &planner.Deposit{}, createDepositAction); err != nil {
		log.Fatal(err)
	}
	if err = engine.RegisterAction("gather", &planner.Gather{}, createGatherAction); err != nil {
		log.Fatal(err)
	}
	if err = engine.AddLocation(loc1); err != nil {
		log.Fatal(err)
	}
	if err = engine.AddLocation(loc2); err != nil {
		log.Fatal(err)
	}
	if err = engine.AddLocation(loc3); err != nil {
		log.Fatal(err)
	}
	if err = engine.AddLocation(depo); err != nil {
		log.Fatal(err)
	}
	if err = engine.AddResource(res1); err != nil {
		log.Fatal(err)
	}
	if err = engine.AddAgent(testAgent); err != nil {
		log.Fatal(err)
	}

	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
