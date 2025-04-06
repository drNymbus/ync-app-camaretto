package model

import (
	// "log"
	// "math"

	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"camaretto/model/component"
	// "camaretto/model/game"
	// "camaretto/model/netplay"
	"camaretto/view"
)

type Game struct {
	width, height float64

	camaretto *component.Camaretto

	attack *component.Button
	shield *component.Button
	charge *component.Button
	heal *component.Button

	cursor *view.Sprite

	info *component.TextBox

	count int
}


// @desc: Initialize attributes of a Camaretto instance, given the number of players: n
func (g *Game) Init(seed int64, names []string, w, h int) {
	g.width, g.height = float64(w), float64(h)

	g.camaretto = &component.Camaretto{}
	g.camaretto.Init(seed, names, g.width, g.height)

	for _, player := range g.camaretto.Players {
		var bodyX float64 = player.Persona.SSprite.Width/2
		var bodyY float64 = g.height - player.Persona.SSprite.Height/2
		player.Persona.SSprite.SetCenter(bodyX, bodyY, 0)
	}

	g.info = component.NewTextBox(g.width - 50, g.height*1/5 + 30, "", color.RGBA{0, 0, 0, 255}, color.RGBA{0, 51, 153, 127})
	var x, y float64 = g.width/2, g.height*8/10 + 65
	g.info.SSprite.SetCenter(x, y, 0)

	var buttonXPos float64 = 0
	var buttonYPos float64 = g.height * 9/10

	g.attack = component.NewButton("ATTACK", color.RGBA{0, 0, 0, 255}, "RED", g.camaretto.AttackHook)
	buttonXPos = (g.width * 1/4) + (float64(view.ButtonWidth)/2)
	g.attack.SSprite.SetCenter(buttonXPos, buttonYPos, 0)

	g.shield = component.NewButton("SHIELD", color.RGBA{0, 0, 0, 255}, "BLUE", g.camaretto.ShieldHook)
	buttonXPos = (g.width * 2/4) + (float64(view.ButtonWidth)/2)
	g.shield.SSprite.SetCenter(buttonXPos, buttonYPos, 0)

	buttonXPos = (g.width * 3/4) + (float64(view.ButtonWidth)/2)

	g.charge = component.NewButton("CHARGE", color.RGBA{0, 0, 0, 255}, "YELLOW", g.camaretto.ChargeHook)
	g.charge.SSprite.SetCenter(buttonXPos, buttonYPos, 0)

	g.heal = component.NewButton("HEAL", color.RGBA{0, 0, 0, 255}, "GREEN", g.camaretto.HealHook)
	g.heal.SSprite.SetCenter(buttonXPos, buttonYPos, 0)

	g.cursor = view.NewSprite(view.LoadCursorImage(), nil)
	g.cursor.SetCenter(-g.cursor.Width, -g.cursor.Height, 0)
	g.cursor.SetOffset(0, 0, 0)

	g.count = 0
}

// @desc: true if the player (Application.PlayerInfo) is required to do an action, false otherwise
func (g *Game) IsMyTurn(index int) bool {
	return true
}

func (g *Game) Update() error {
	g.camaretto.Update()
	g.info.Update()

	var player *component.Player = g.camaretto.Players[g.camaretto.Current.PlayerTurn]
	player.Persona.Update()

	if g.camaretto.Current.State == component.SET {
		g.attack.Update()
		g.shield.Update()
		if player.IsChargeEmpty() {
			g.charge.Update()
		} else {
			g.heal.Update()
		}
	}

	return nil
}

// @desc: Render all elements on a given image (dst)
func (g *Game) Draw(screen *ebiten.Image) {
	var player *component.Player = g.camaretto.Players[g.camaretto.Current.PlayerTurn]
	player.Persona.Draw(screen)
	g.info.Draw(screen)

	if g.camaretto.Current.State == component.SET {
		g.attack.Draw(screen)
		g.shield.Draw(screen)
		if player.IsChargeEmpty() {
			g.charge.Draw(screen)
		} else {
			g.heal.Draw(screen)
		}
	}

	g.camaretto.Draw(screen)
}
