package main

import (
	"log"
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"camaretto/model"
	"camaretto/event"
	"camaretto/view"
)

var (
	err error
)

type Game struct{
	application *model.Application
	events *event.EventQueue
}

func NewGame(nbPlayers int) *Game {
	var g *Game = &Game{}

	g.application = &model.Application{}
	g.application.Init(3)
	g.events = event.NewEventQueue(20)

	return g
}

func (g *Game) Update() error {
	g.events.Update()

	g.application.Hover(g.events.X, g.events.Y)

	var me *event.MouseEvent = nil
	var ke *event.KeyEvent = nil
	for ;!g.events.IsEmpty(); {
		me = g.events.ReadMouseEvent()
		if me != nil { g.application.MouseEventUpdate(me) }
		ke = g.events.ReadKeyEvent()
		if ke != nil { g.application.KeyEventUpdate(ke) }
	}

	g.application.Update()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)
	var img *ebiten.Image = g.application.Display()
	screen.DrawImage(img, nil)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return model.WinWidth, model.WinHeight
}

func main() {
	// Loading assets
	view.LoadFont()
	// view.InitAssets()

	// Init Game
	var g *Game = NewGame(5)

	// Init Window
	ebiten.SetWindowSize(model.WinWidth, model.WinHeight)
	ebiten.SetWindowTitle("Camaretto")

	var icon image.Image
	icon, err = view.InitIcon("assets/amaretto_icon.png")
	if err != nil {
		log.Fatal("[MAIN] InitIcon failed", err)
	}
	ebiten.SetWindowIcon([]image.Image{icon})

	// Game Loop
	if err = ebiten.RunGame(g); err != nil {
		log.Fatal("[MAIN]", err)
	}

	// Free resources
}
