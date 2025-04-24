package scene

import (
	// "log"
	// "math"

	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"camaretto/model/ui"
	"camaretto/model/game"
	"camaretto/view"
)

type Game struct {
	width, height float64

	online bool
	playerInfo *game.PlayerInfo

	Camaretto *game.Camaretto

	attack *ui.Button
	shield *ui.Button
	charge *ui.Button
	heal *ui.Button

	cursor *view.Sprite

	info *ui.TextBox

	count int
	gotoEnd func()
}


// @desc: Initialize attributes of a Camaretto instance, given the number of players: n
func (g *Game) Init(seed int64, names []string, w, h int, online bool, player *game.PlayerInfo, endRoutine func()) {
	g.width, g.height = float64(w), float64(h)

	g.online = online
	g.playerInfo = player

	g.Camaretto = &game.Camaretto{}
	g.Camaretto.Init(seed, names, g.width, g.height * 8/10)

	for _, player := range g.Camaretto.Players {
		var bodyX float64 = player.Persona.SSprite.Width/2
		var bodyY float64 = g.height - player.Persona.SSprite.Height/2
		player.Persona.SSprite.SetCenter(bodyX, bodyY, 0)
	}

	g.info = ui.NewTextBox(g.width - 50, g.height*1/5 + 30, "", color.RGBA{0, 0, 0, 255}, color.RGBA{0, 51, 153, 127})
	var x, y float64 = g.width/2, g.height*8/10 + 65
	g.info.SSprite.SetCenter(x, y, 0)

	var buttonXPos float64 = 0
	var buttonYPos float64 = g.height * 9/10

	g.attack = ui.NewButton("ATTACK", color.RGBA{0, 0, 0, 255}, "RED", func(){ g.Camaretto.Current.State = game.ATTACK })
	buttonXPos = (g.width * 1/4) + (float64(view.ButtonWidth)/2)
	g.attack.SSprite.SetCenter(buttonXPos, buttonYPos, 0)

	g.shield = ui.NewButton("SHIELD", color.RGBA{0, 0, 0, 255}, "BLUE", func(){ g.Camaretto.Current.State = game.SHIELD })
	buttonXPos = (g.width * 2/4) + (float64(view.ButtonWidth)/2)
	g.shield.SSprite.SetCenter(buttonXPos, buttonYPos, 0)

	buttonXPos = (g.width * 3/4) + (float64(view.ButtonWidth)/2)

	g.charge = ui.NewButton("CHARGE", color.RGBA{0, 0, 0, 255}, "YELLOW", func(){ g.Camaretto.Current.State = game.CHARGE })
	g.charge.SSprite.SetCenter(buttonXPos, buttonYPos, 0)

	g.heal = ui.NewButton("HEAL", color.RGBA{0, 0, 0, 255}, "GREEN", func(){ g.Camaretto.Current.State = game.HEAL })
	g.heal.SSprite.SetCenter(buttonXPos, buttonYPos, 0)

	g.cursor = view.NewSprite(view.LoadCursorImage(), nil)
	g.cursor.SetCenter(-g.cursor.Width, -g.cursor.Height, 0)
	g.cursor.SetOffset(0, 0, 0)

	g.count = 0
	g.gotoEnd = endRoutine
}

// @desc: true if the player (Application.PlayerInfo) is required to do an action, false otherwise
func (g *Game) IsMyTurn() bool {
	if g.Camaretto.Current.Focus == game.CARD {
		return (g.playerInfo.Index == g.Camaretto.Current.PlayerFocus)
	}

	return (g.playerInfo.Index == g.Camaretto.Current.PlayerTurn)
}

func (g *Game) Update() error {
	g.Camaretto.Update(g.cursor)
	if g.Camaretto.IsGameOver() { g.gotoEnd() }

	g.info.Update()

	var player *game.Player = g.Camaretto.Players[g.Camaretto.Current.PlayerTurn]
	player.Persona.Update()

	if g.online && !g.IsMyTurn() { return nil }

	if g.Camaretto.Current.State == game.SET {
		g.attack.Update(g.cursor)
		g.shield.Update(g.cursor)
		if player.IsChargeEmpty() {
			g.charge.Update(g.cursor)
		} else {
			g.heal.Update(g.cursor)
		}
	}

	g.cursor.Update()

	return nil
}

// @desc: Render all elements on a given image (dst)
func (g *Game) Draw(screen *ebiten.Image) {
	var player *game.Player = g.Camaretto.Players[g.Camaretto.Current.PlayerTurn]
	player.Persona.Draw(screen)
	g.info.Draw(screen)

	if g.Camaretto.Current.State == game.SET {
		g.attack.Draw(screen)
		g.shield.Draw(screen)
		if player.IsChargeEmpty() {
			g.charge.Draw(screen)
		} else {
			g.heal.Draw(screen)
		}
	}

	g.Camaretto.Draw(screen)
	g.cursor.Draw(screen)
}
