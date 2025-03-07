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
	mouse *event.Mouse
}

func NewGame(nbPlayers int) *Game {
	var g *Game = &Game{}

	g.application = &model.Application{}
	g.application.Init(nbPlayers)

	g.mouse = event.NewMouse(20)

	return g
}

func (g *Game) Update() error {
	g.mouse.Update()

	event.HandleGameHover(g.application, g.mouse.X, g.mouse.Y)

	var e *event.MouseEvent = nil
	for ;!g.mouse.IsEmpty(); {
		e = g.mouse.ReadEvent()

		if e.MET == event.RELEASED && e.Click == ebiten.MouseButtonLeft {
			event.HandleCamarettoMouseRelease(g.application, float64(e.X), float64(e.Y))
		} else if e.MET == event.PRESSED && e.Click == ebiten.MouseButtonLeft {
			event.HandleCamarettoMousePress(g.application, float64(e.X), float64(e.Y))
		}
	}

	g.application.Update()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)

	g.application.DrawPlayers(screen)

	g.application.DrawDeck(screen)
	g.application.DrawCenterCards(screen)

	g.application.DrawButtons(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return model.WinWidth, model.WinHeight
}

func main() {
	// Loading assets
	view.InitAssets()

	// Init Game
	var g *Game = NewGame(3)

	// Init Window
	ebiten.SetWindowSize(model.WinWidth, model.WinHeight)
	ebiten.SetWindowTitle("Camaretto")

	var icon image.Image
	icon, err = view.InitIcon("assets/amaretto_trans.png")
	if err != nil {
		log.Fatal("[MAIN] InitIcon failed", err)
	}
	ebiten.SetWindowIcon([]image.Image{icon})

	// Game Loop
	if err = ebiten.RunGame(g); err != nil {
		log.Fatal("[MAIN]", err)
	}
}