package model

import (
	"log"
	"math"

	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"camaretto/model/component"
	"camaretto/model/game"
	"camaretto/model/netplay"
	"camaretto/view"
)

type Game struct {
	width, height float64

	camaretto *game.Camaretto

	attackButton *component.Button
	shieldButton *component.Button
	chargeButton *component.Button
	healButton *component.Button

	cursor *view.Sprite

	info *component.TextBox

	count int
}

// @desc: Initialize attributes of a Camaretto instance, given the number of players: n
func (g *Game) Init(seed int64, n int, names []string, w, h int) {
	if len(names) != n { log.Fatal("[Camaretto.Init] You finna start a game like that ?!") }

	g.width, g.height = float64(w), float64(h)

	g.camaretto = &game.Camaretto{}
	g.camaretto.Init(seed, n, names)

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

	g.attackButton = component.NewButton("ATTACK", color.RGBA{0, 0, 0, 255}, "RED")
	buttonXPos = (g.width * 1/4) + (float64(view.ButtonWidth)/2)
	g.attackButton.SSprite.SetCenter(buttonXPos, buttonYPos, 0)

	g.shieldButton = component.NewButton("SHIELD", color.RGBA{0, 0, 0, 255}, "BLUE")
	buttonXPos = (g.width * 2/4) + (float64(view.ButtonWidth)/2)
	g.shieldButton.SSprite.SetCenter(buttonXPos, buttonYPos, 0)

	buttonXPos = (g.width * 3/4) + (float64(view.ButtonWidth)/2)

	g.chargeButton = component.NewButton("CHARGE", color.RGBA{0, 0, 0, 255}, "YELLOW")
	g.chargeButton.SSprite.SetCenter(buttonXPos, buttonYPos, 0)

	g.healButton = component.NewButton("HEAL", color.RGBA{0, 0, 0, 255}, "GREEN")
	g.healButton.SSprite.SetCenter(buttonXPos, buttonYPos, 0)

	g.cursor = view.NewSprite(view.LoadCursorImage(), false, color.RGBA{0, 0, 0, 0}, nil)
	g.cursor.SetCenter(-g.cursor.Width, -g.cursor.Height, 0)
	g.cursor.SetOffset(0, 0, 0)

	g.count = 0
}

// @desc: true if the player (Application.PlayerInfo) is required to do an action, false otherwise
func (g *Game) IsMyTurn(index int) bool {
	var flag bool = false

	var action *game.Action = g.camaretto.Current
	if action.State == game.SET {
		flag = (action.PlayerTurn == index)
	} else {
		if action.Focus == game.PLAYER || action.Focus == game.REVEAL {
			flag = (action.PlayerTurn == index)
		} else if action.Focus == game.CARD {
			flag = (action.PlayerFocus == index)
		}
	}

	return flag
}

func (g *Game) SerializeCamaretto() *netplay.Message {
	var msg *netplay.Message = &netplay.Message{}

	msg.Typ = netplay.ACTION
	msg.Seed = g.camaretto.DeckPile.Seed

	msg.Players = []*game.PlayerInfo{}
	for i, player := range g.camaretto.Players {
		var info *game.PlayerInfo = &game.PlayerInfo{i, player.Name}
		msg.Players = append(msg.Players, info)
	}

	msg.Action = g.camaretto.Current

	for _, card := range g.camaretto.ToReveal {
		if card.Hidden { msg.Reveal = append(msg.Reveal, false) }
	}

	return msg
}

func (g *Game) DeserializeCamaretto(msg *netplay.Message) {
	g.camaretto.Current = msg.Action

	for i, reveal := range msg.Reveal {
		if reveal { g.camaretto.ToReveal[i].Reveal() }
	}
}

/************ *************************************************************************** ************/
/************ ********************************** UPDATE ********************************* ************/
/************ *************************************************************************** ************/

// @desc: If the coordinates (x,y) are on any of the cards to be revealed returns the card's index, returns -1 otherwise
func (g *Game) onReveal(x, y float64) int {
	for i, card := range g.camaretto.ToReveal {
		if card.SSprite.In(x, y) { return i }
	}
	return -1
}

// @desc: If the coordinates (x,y) are on any of the focused player's health card return the card index, returns -1 otherwise
func (g *Game) onHealth(x, y float64) int {
	var c *game.Camaretto = g.camaretto
	var p *game.Player = c.Players[c.Current.PlayerFocus]
	if p.HealthCard[0] != nil && p.HealthCard[0].SSprite.In(x, y) {
		return 0
	} else if p.HealthCard[1] != nil && p.HealthCard[1].SSprite.In(x, y) {
		return 1
	}
	return -1
}

// @desc: If the coordinates (x,y) are on any player's card the player's index is returned, -1 otherwise
func (g *Game) onPlayer(x, y float64) int {
	for i, player := range g.camaretto.Players {
		if !player.Dead {
			var onPlayer bool = false
			if player.HealthCard[0] != nil { onPlayer = onPlayer || player.HealthCard[0].SSprite.In(x, y) }
			if player.HealthCard[1] != nil { onPlayer = onPlayer || player.HealthCard[1].SSprite.In(x, y) }
			if player.JokerHealth != nil { onPlayer = onPlayer || player.JokerHealth.SSprite.In(x, y) }
			if player.ShieldCard != nil { onPlayer = onPlayer || player.ShieldCard.SSprite.In(x, y) }
			if player.ChargeCard != nil { onPlayer = onPlayer || player.ChargeCard.SSprite.In(x, y) }

			if onPlayer { return i }
		}
	}
	return -1
}

// @desc: Update elements upon mouse hover on coordinates (x,y)
func (g *Game) Hover(x, y float64) {
	var speed float64 = 15

	var s *view.Sprite = nil
	if g.camaretto.Current.State == game.SET {
		if g.attackButton.SSprite.In(x, y) {
			s = g.attackButton.SSprite
		} else if g.shieldButton.SSprite.In(x, y) {
			s = g.shieldButton.SSprite
		} else if g.chargeButton.SSprite.In(x, y) {
			s = g.chargeButton.SSprite
		} else if g.healButton.SSprite.In(x, y) {
			s = g.healButton.SSprite
		}

		if s != nil {
			var x, y, _ float64 = s.GetCenter()
			g.cursor.Move(x - (s.Width/2), y, speed)
			g.cursor.Rotate(math.Pi/2, speed)
			g.cursor.MoveOffset(0, 0, speed)
			g.cursor.RotateOffset(0, speed)
		}
	} else {
		if g.camaretto.Current.Focus == game.PLAYER {
			var i int = g.onPlayer(x, y)
			if i != -1 {
				s = g.camaretto.Players[i].NameSprite
				var x, y, r float64 = s.GetCenter()
				g.cursor.Move(x, y, speed)
				g.cursor.Rotate(math.Pi, speed)
				x, y, r = s.GetOffset()
				g.cursor.MoveOffset(x, y - float64(view.CardHeight*5/2), speed)
				g.cursor.RotateOffset(r, speed)
			}
		} else if g.camaretto.Current.Focus == game.CARD {
			var i int = g.onHealth(x, y)
			if i != -1 {
				s = g.camaretto.Players[g.camaretto.Current.PlayerFocus].HealthCard[i].SSprite
				var x, y, r float64 = s.GetCenter()
				g.cursor.Move(x, y, speed)
				g.cursor.Rotate(math.Pi, speed)
				x, y, r = s.GetOffset()
				g.cursor.MoveOffset(x, y - float64(view.CardHeight/2), speed)
				g.cursor.RotateOffset(r, speed)
			}
		} else if g.camaretto.Current.Focus == game.REVEAL {
			for _, card := range g.camaretto.ToReveal {
				if card.SSprite.In(x, y) { s = card.SSprite }
			}

			if s != nil {
				var x, y, r float64 = s.GetCenter()
				g.cursor.Move(x, y + (s.Height/2), speed)
				g.cursor.Rotate(r, speed)
				g.cursor.MoveOffset(0, 0, speed)
				g.cursor.RotateOffset(0, speed)
			}
		}
	}
}

// @desc: Update elements on mouse button press on coordinates (x,y)
func (g *Game) MousePress(x, y float64) {
	if g.camaretto.Current.State == game.SET {
		if g.attackButton.SSprite.In(x, y) {
			g.attackButton.Pressed()
		} else if g.shieldButton.SSprite.In(x, y) {
			g.shieldButton.Pressed()
		} else {
			if g.camaretto.Players[g.camaretto.Current.PlayerTurn].ChargeCard == nil && g.chargeButton.SSprite.In(x, y) {
				g.chargeButton.Pressed()
			} else if g.healButton.SSprite.In(x, y) {
				g.healButton.Pressed()
			}
		}
	}
}

// @desc: Update elements on mouse button release on coordinates (x,y)
func (g *Game) MouseRelease(x, y float64) {
	g.cursor.SetCenter(-g.cursor.Width, -g.cursor.Height, 0)
	g.cursor.SetOffset(0, 0, 0)

	g.attackButton.Released()
	g.shieldButton.Released()
	g.chargeButton.Released()
	g.healButton.Released()

	var action *game.Action = g.camaretto.Current
	if action.State == game.SET {
		if g.attackButton.SSprite.In(x, y) {
			action.State = game.ATTACK
			action.Focus = game.PLAYER
		} else if g.shieldButton.SSprite.In(x, y) {
			action.State = game.SHIELD
			action.Focus = game.PLAYER
		} else {
			if g.camaretto.Players[action.PlayerTurn].ChargeCard == nil && g.chargeButton.SSprite.In(x, y) {
				action.State = game.CHARGE
				action.Focus = game.COMPLETE
				action.PlayerFocus = action.PlayerTurn
			} else if g.healButton.SSprite.In(x, y) {
				action.State = game.HEAL
				action.Focus = game.CARD
				action.PlayerFocus = action.PlayerTurn
			}
		}
	} else {
		if action.Focus == game.PLAYER {
			var i int = g.onPlayer(x, y)
			if i != -1 {
				if action.State == game.ATTACK {
					action.Focus = game.CARD
					action.PlayerFocus = i
				} else if action.State == game.SHIELD {
					action.Focus = game.REVEAL
					action.PlayerFocus = i
					g.camaretto.Reveal()
				}
			}
		} else if action.Focus == game.CARD {
			var i int = g.onHealth(x, y)
			if i != -1 {
				action.Focus = game.REVEAL
				action.CardFocus = i
				g.camaretto.Reveal()
			}
		} else if action.Focus == game.REVEAL {
			var i int = g.onReveal(x, y)
			if i != -1 { g.camaretto.ToReveal[i].Reveal() }
		}
	}

	g.camaretto.Current = action
}

// @desc: Complete actions that do not require user input
func (g *Game) Update() component.PageSignal {
	if g.camaretto.IsGameOver() { return component.NEXT }

	var c *game.Camaretto = g.camaretto
	if c.Current.Focus == game.REVEAL {
		if len(c.ToReveal) == 0 {
			c.Current.Focus = game.COMPLETE
		} else {
			var done bool = true
			for _, card := range c.ToReveal { done = done && (!card.Hidden) }

			if done { g.count++ }
			if done && g.count > 33 {
				c.Current.Focus = game.COMPLETE
				g.count = 0
			}
		}
	} else if c.Current.Focus == game.COMPLETE {
		if c.Current.State != game.SET {
			if c.Current.State == game.ATTACK {
				c.Attack()
			} else if c.Current.State == game.SHIELD {
				c.Shield()
			} else if c.Current.State == game.CHARGE {
				c.Charge()
			} else if c.Current.State == game.HEAL {
				c.Heal()
			}
			c.EndTurn()
		}
	}

	return component.UPDATE
}

/************ *************************************************************************** ************/
/************ ********************************** RENDER ********************************* ************/
/************ *************************************************************************** ************/

// @desc: Turn index into a coordinate (x,y) and an angle (theta)
func (g *Game) getPlayerGeoM(i int) (float64, float64, float64) {
	var angleStep float64 = 2*math.Pi / float64(g.camaretto.NbPlayers)
	var radius float64 = 200

	var theta float64 = angleStep * float64(i)
	var x float64 = radius * math.Cos(theta + math.Pi/2)
	var y float64 = radius * math.Sin(theta + math.Pi/2)

	return x, y, theta
}

// @desc: Render all elements on a given image (dst)
func (g *Game) Display(dst *ebiten.Image) {
	var persona *game.Character = g.camaretto.Players[g.camaretto.Current.PlayerTurn].Persona
	if g.info.Finished() {
		persona.Talking = false
	} else {
		persona.Talking = true
		g.info.Render()
	}
	persona.Render()
	persona.SSprite.Display(dst)
	g.info.SSprite.Display(dst)

	if g.camaretto.Current.State == game.SET {
		g.attackButton.SSprite.Display(dst)
		g.shieldButton.SSprite.Display(dst)

		if g.camaretto.Players[g.camaretto.Current.PlayerTurn].ChargeCard == nil {
			g.chargeButton.SSprite.Display(dst)
		} else {
			g.healButton.SSprite.Display(dst)
		}
	}

	var centerX float64 = g.width/2
	var centerY float64 = (g.height * 6/8)/2

	for i, player := range g.camaretto.Players {
		var x, y, theta float64 = g.getPlayerGeoM(i)
		player.RenderCards(dst, centerX + x, centerY + y, theta)
	}

	if len(g.camaretto.ToReveal) > 0 {
		for i, card := range g.camaretto.ToReveal {
			var s *view.Sprite = card.SSprite
			var x float64 = (float64(i) - float64(len(g.camaretto.ToReveal)-1)/2) * float64(view.CardWidth)
			s.Move(centerX + x, centerY, 0.5)
			card.SSprite.Rotate(0, 0.2)
			s.MoveOffset(0, 0, 0.2)
			card.SSprite.RotateOffset(0, 0.2)
			card.SSprite.Display(dst)
		}
	} else {
		g.camaretto.DeckPile.Render(dst, centerX, centerY)
	}

	g.cursor.Display(dst)
}
