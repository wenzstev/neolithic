package main

import (
	"Neolithic/internal/attributes"
	"log"
	"os"
	"runtime/pprof"

	"Neolithic/internal/agent"
	"Neolithic/internal/camera"
	"Neolithic/internal/core"
	"Neolithic/internal/goalengine"
	"Neolithic/internal/grid"
	"Neolithic/internal/logging"
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

	cpuProfileFile, err := os.Create("cpu.pprof")
	if err != nil {
		log.Fatal("could not create CPU profile: ", err)
	}
	defer cpuProfileFile.Close() // Make sure to close the file

	if err := pprof.StartCPUProfile(cpuProfileFile); err != nil {
		log.Fatal("could not start CPU profile: ", err)
	}
	defer pprof.StopCPUProfile() // Make sure to stop profiling when main exits

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

	baseCapacityAttr := &attributes.Capacity{Size: 100}

	loc1 := core.NewLocation("loc1", core.Coord{X: 3, Y: 14}, core.WithAttributes(baseCapacityAttr))
	loc2 := core.NewLocation("loc2", core.Coord{X: 21, Y: 4}, core.WithAttributes(baseCapacityAttr))
	loc3 := core.NewLocation("loc3", core.Coord{X: 27, Y: 30}, core.WithAttributes(baseCapacityAttr))
	depo := core.NewLocation("depo", core.Coord{X: 16, Y: 16}, core.WithAttributes(baseCapacityAttr))

	res1 := core.NewResource("Berries", core.WithResourceAttributes(&attributes.Weight{Amount: 1}))
	res2 := core.NewResource("Wood", core.WithResourceAttributes(&attributes.Weight{Amount: 1}))
	res3 := core.NewResource("Stone", core.WithResourceAttributes(&attributes.Weight{Amount: 1}))

	loc1.Inventory.AdjustAmount(res1, 2000)
	loc2.Inventory.AdjustAmount(res1, 1000)
	loc3.Inventory.AdjustAmount(res1, 2000)

	goalDepo := depo.DeepCopy()

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
	if err = engine.AddResource(res2); err != nil {
		log.Fatal(err)
	}
	if err = engine.AddResource(res3); err != nil {
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
