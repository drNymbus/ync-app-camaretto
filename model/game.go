package model

import (
	// "log"
	"math"

	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"camaretto/model/component"
	"camaretto/view"
)

type Game struct {
	width, height float64

	camaretto *component.Camaretto

	attack *component.Button
	shield *component.Button
	charge *component.Button
	heal *component.Button

	cursor *component.Sprite

	info *component.TextBox

	count int
	gotoEnd func()
}


// @desc: Initialize attributes of a Camaretto instance, given the number of players: n
func (g *Game) Init(seed int64, names []string, w, h int, endRoutine func()) {
	g.width, g.height = float64(w), float64(h)

	g.camaretto = &component.Camaretto{}
	g.camaretto.Init(seed, names, g.width, g.height * 8/10)

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

	g.cursor = component.NewSprite(view.LoadCursorImage(), nil)
	g.cursor.SetCenter(-g.cursor.Width, -g.cursor.Height, 0)
	g.cursor.SetOffset(0, 0, 0)

	g.count = 0
	g.gotoEnd = endRoutine
}

// @desc: true if the player (Application.PlayerInfo) is required to do an action, false otherwise
func (g *Game) IsMyTurn(index int) bool {
	return true
}

func (g *Game) Update() error {
	g.camaretto.Update()
	if g.camaretto.IsGameOver() { g.gotoEnd() }
	g.info.Update()

	var player *component.Player = g.camaretto.Players[g.camaretto.Current.PlayerTurn]
	player.Persona.Update()

	var ix, iy int = ebiten.CursorPosition()
	var x, y float64 = float64(ix), float64(iy)

	var cx, cy, cr float64
	var cxOff, cyOff, crOff float64
	var cursorSpeed float64 = 25

	if g.camaretto.Current.State == component.SET {
		g.attack.Update()
		g.shield.Update()
		if g.attack.SSprite.In(x, y) {
			cx, cy, cr = g.attack.SSprite.GetCenter()
			cx = cx - g.attack.SSprite.Width/2
			cr = cr + math.Pi/2

			g.cursor.Move(cx, cy, cursorSpeed)
			g.cursor.Rotate(cr, cursorSpeed)
			g.cursor.MoveOffset(cxOff, cyOff, cursorSpeed)
			g.cursor.RotateOffset(crOff, cursorSpeed)
		} else if g.shield.SSprite.In(x, y) {
			cx, cy, cr = g.shield.SSprite.GetCenter()
			cx = cx - g.shield.SSprite.Width/2
			cr = cr + math.Pi/2

			g.cursor.Move(cx, cy, cursorSpeed)
			g.cursor.Rotate(cr, cursorSpeed)
			g.cursor.MoveOffset(cxOff, cyOff, cursorSpeed)
			g.cursor.RotateOffset(crOff, cursorSpeed)
		}

		if player.IsChargeEmpty() {
			g.charge.Update()
			if g.charge.SSprite.In(x, y) {
				cx, cy, cr = g.charge.SSprite.GetCenter()
				cx = cx - g.charge.SSprite.Width/2
				cr = cr + math.Pi/2

				g.cursor.Move(cx, cy, cursorSpeed)
				g.cursor.Rotate(cr, cursorSpeed)
				g.cursor.MoveOffset(cxOff, cyOff, cursorSpeed)
				g.cursor.RotateOffset(crOff, cursorSpeed)
			}
		} else {
			g.heal.Update()
			if g.heal.SSprite.In(x, y) {
				cx, cy, cr = g.heal.SSprite.GetCenter()
				cx = cx - g.heal.SSprite.Width/2
				cr = cr + math.Pi/2

				g.cursor.Move(cx, cy, cursorSpeed)
				g.cursor.Rotate(cr, cursorSpeed)
				g.cursor.MoveOffset(cxOff, cyOff, cursorSpeed)
				g.cursor.RotateOffset(crOff, cursorSpeed)
			}
		}

	} else if g.camaretto.Current.Focus == component.PLAYER {
		for _, player := range g.camaretto.Players {
			if player.HoverPlayer(x, y) {
				cx, cy, crOff = player.GetPosition()
				cyOff = -float64(view.CardHeight)
				cr = math.Pi

				g.cursor.Move(cx, cy, cursorSpeed)
				g.cursor.Rotate(cr, cursorSpeed)
				g.cursor.MoveOffset(cxOff, cyOff, cursorSpeed)
				g.cursor.RotateOffset(crOff, cursorSpeed)
			}
		}
	} else if g.camaretto.Current.Focus == component.CARD {
		var player *component.Player = g.camaretto.Players[g.camaretto.Current.PlayerFocus]
		var i int = player.HoverHealth(x, y)
		if i != -1 {
			cx, cy, crOff = player.GetPosition()
			cxOff, cyOff, _ = player.Health[i].SSprite.GetOffset()
			cyOff = cyOff - float64(view.CardHeight)/2
			cr = math.Pi

			g.cursor.Move(cx, cy, cursorSpeed)
			g.cursor.Rotate(cr, cursorSpeed)
			g.cursor.MoveOffset(cxOff, cyOff, cursorSpeed)
			g.cursor.RotateOffset(crOff, cursorSpeed)
		}
	}


	g.cursor.Update()
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
	g.cursor.Draw(screen)
}
