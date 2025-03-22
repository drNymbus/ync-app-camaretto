package game

import (
	"log"
	"math"

	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"camaretto/model/component"
	"camaretto/view"
	"camaretto/event"
)

type GameState int
const (
	SET GameState = 0
	ATTACK GameState = 1
	SHIELD GameState = 2
	CHARGE GameState = 3
	HEAL GameState = 4
)

type FocusState int
const (
	NONE FocusState = 0
	PLAYER FocusState = 1
	CARD FocusState = 2
	REVEAL FocusState = 3
	COMPLETE FocusState = 4
)

/************ *************************************************************************** ************/
/************ ******************************** CAMARETTO ******************************** ************/
/************ *************************************************************************** ************/

type Camaretto struct {
	state GameState
	focus FocusState

	playerTurn int
	playerFocus int
	cardFocus int

	nbPlayers int
	Players []*Player
	DeckPile *Deck

	toReveal []*Card

	attackButton *component.Button
	shieldButton *component.Button
	chargeButton *component.Button
	healButton *component.Button

	cursor *view.Sprite

	info *component.TextBox

	count int
}

// @desc: Initialize attributes of a Camaretto instance, given the number of players: n
// func (c *Camaretto) Init(n int, sheet *ebiten.Image, tileWidth int, tileHeight int) {
func (c *Camaretto) Init(n int, names []string, seed int64, width, height float64) {
	if len(names) != n { log.Fatal("[camaretto.Init] You finna start a game like that ?!") }

	c.state = SET
	c.focus = NONE

	c.playerTurn = 0
	c.playerFocus = -1
	c.cardFocus = -1

	c.DeckPile = &Deck{}
	c.DeckPile.Init(seed)

	c.nbPlayers = n
	c.Players = make([]*Player, n)

	c.info = component.NewTextBox(width - 50, height*1/5 + 30, "", color.RGBA{0, 0, 0, 255}, color.RGBA{0, 51, 153, 127})
	var x, y float64 = width/2, height*8/10 + 65
	c.info.SSprite.SetCenter(x, y, 0)

	for i, _ := range make([]int, n) { // Init players
		var name string = names[i%len(names)]
		var char *Character = NewCharacter(name)
		var bodyX float64 = (x - c.info.SSprite.Width/2) + char.SSprite.Width/2
		var bodyY float64 = (y + c.info.SSprite.Height/2) - char.SSprite.Height/2
		char.SSprite.SetCenter(bodyX, bodyY, 0)

		c.Players[i] = NewPlayer(name, char)
	}

	for i, _ := range make([]int, n*2) { // Init Health
		var card *Card = c.DeckPile.DrawCard()
		if card.Name == "Joker" {
			c.DeckPile.DiscardCard(card)
			card = c.DeckPile.DrawCard()
			if card.Name == "Joker" { log.Fatal("[Camaretto.Init] (Health) 2 JOKERS IN A ROW ?! What were the chances anyway...") }

			card.Reveal()
			c.Players[i%n].JokerHealth = card

			card = c.DeckPile.DrawCard()
			if card.Name == "Joker" { log.Fatal("[Camaretto.Init] (Health) 2 JOKERS SPACED BY ONE CARD ?! Ok this one might not be THAT crazy but still...") }
		}

		card.Reveal()
		c.Players[i%n].HealthCard[int(i/n)] = card
	}

	for i, _ := range make([]int, n) { // Init Shield
		var card *Card = c.DeckPile.DrawCard()
		card.Reveal()
		if card.Name == "Joker" {
			c.Players[i].JokerShield = card
		} else {
			c.Players[i].ShieldCard = card
		}
	}

	c.attackButton = component.NewButton("ATTACK", color.RGBA{0, 0, 0, 255}, "RED")
	c.shieldButton = component.NewButton("SHIELD", color.RGBA{0, 0, 0, 255}, "BLUE")
	c.chargeButton = component.NewButton("CHARGE", color.RGBA{0, 0, 0, 255}, "YELLOW")
	c.healButton = component.NewButton("HEAL", color.RGBA{0, 0, 0, 255}, "GREEN")

	c.cursor = view.NewSprite(view.LoadCursorImage(), false, color.RGBA{0, 0, 0, 0}, nil)
	c.cursor.SetCenter(-c.cursor.Width, -c.cursor.Height, 0)
	c.cursor.SetOffset(0, 0, 0)

	c.count = 0
}

/************ ***************************************************************************** ************/
/************ ********************************** ACTIONS ********************************** ************/
/************ ***************************************************************************** ************/

// @desc: Returns true if one player is left, false otherwise
func (c *Camaretto) IsGameOver() bool {
	var count int = c.nbPlayers
	for _, p := range c.Players {
		if p.Dead { count-- }
	}

	if count > 1 { return false }
	return true
}

func (c *Camaretto) endTurn() {
	c.state = SET
	c.focus = NONE
	c.playerFocus = -1
	c.cardFocus = -1

	c.playerTurn = (c.playerTurn+1) % c.nbPlayers
	for ;c.Players[c.playerTurn].Dead; { c.playerTurn = (c.playerTurn+1) % c.nbPlayers }
}

func (c *Camaretto) attackPlayer(dst *Player, amount int) {
	if dst.JokerShield != nil {
		c.DeckPile.DiscardCard(dst.JokerShield)
		dst.JokerShield = nil
	} else {
		// In which health slot should we put a new card ?
		var healthSlot int = -1 // -1: none; 0: health[0]; 1:health[1]; 2: joker

		if dst.ShieldCard != nil { amount = amount - dst.ShieldCard.Value }
		if amount > 0 {

			// Do we have a joker health ? Then it's tanking (wether you like it or not)
			if dst.JokerHealth != nil {
				amount = amount - dst.JokerHealth.Value
				c.DeckPile.DiscardCard(dst.JokerHealth)
				dst.JokerHealth = nil
				// We have to replace your jokerHealth then
				healthSlot = 2
			}
	
			// Is the attack still going ?
			if amount > 0 {
				amount = amount - dst.HealthCard[c.cardFocus].Value
				c.DeckPile.DiscardCard(dst.HealthCard[c.cardFocus])
				dst.HealthCard[c.cardFocus] = nil
				// Joker's gone, we replace the health you focused
				healthSlot = c.cardFocus
			}
	
			// Wow that's a really big hit
			if amount > 0 && dst.HealthCard[1-c.cardFocus] != nil {
				amount = amount - dst.HealthCard[1-c.cardFocus].Value
				c.DeckPile.DiscardCard(dst.HealthCard[1-c.cardFocus])
				dst.HealthCard[1-c.cardFocus] = nil
				// Both of your health cards took a hit ? Guess you don't have an option anymore
				healthSlot = 0
			}

			// R.I.P in Peperonni
			if amount >= 0 {
				// Give all your cards to your little friends pls
				if dst.HealthCard[0] != nil {
					c.DeckPile.DiscardCard(dst.HealthCard[0])
					dst.HealthCard[0] = nil
				}
				if dst.HealthCard[1] != nil {
					c.DeckPile.DiscardCard(dst.HealthCard[1])
					dst.HealthCard[1] = nil
				}
				if dst.JokerHealth != nil {
					c.DeckPile.DiscardCard(dst.JokerHealth)
					dst.JokerHealth = nil
				}
				if dst.ShieldCard != nil {
					c.DeckPile.DiscardCard(dst.ShieldCard)
					dst.ShieldCard = nil
				}
				if dst.JokerShield != nil {
					c.DeckPile.DiscardCard(dst.JokerShield)
					dst.JokerShield = nil
				}
				dst.Dead = true
			} else { // Recovering
				amount = amount * -1

				var newHealthCard *Card = nil
				newHealthCard = c.DeckPile.FindInDiscardPile(amount)
				if newHealthCard == nil {
					newHealthCard = c.DeckPile.FindInDrawPile(amount)
				}
				if newHealthCard == nil { log.Fatal("[Camaretto.Attack] Could not find a card with health points left") }

				newHealthCard.Reveal()

				if healthSlot == 2 {
					dst.JokerHealth = newHealthCard
				} else {
					dst.HealthCard[healthSlot] = newHealthCard
				}
			}
		} else if amount == 0 {
			c.DeckPile.DiscardCard(dst.ShieldCard)
			dst.ShieldCard = nil
		}
	}
}

// @desc: Player at index src attacks the player at index dst
func (c *Camaretto) attack() (int, string) {
	var atkCard *Card = c.toReveal[0]
	var chargeCard *Card = nil
	if len(c.toReveal) == 2 { chargeCard = c.toReveal[1] }

	var atkValue int = atkCard.Value
	if chargeCard != nil { atkValue = atkValue + chargeCard.Value }

	c.attackPlayer(c.Players[c.playerFocus], atkValue)

	c.DeckPile.DiscardCard(atkCard)
	if chargeCard != nil { c.DeckPile.DiscardCard(chargeCard) }
	c.toReveal = []*Card{}


	return 0, "Attack was great, might do it again ! 5/5"
}

// @desc: Player at index player gets assigned a new shield
func (c *Camaretto) shield() (int, string) {
	var oldCard *Card = c.Players[c.playerFocus].ShieldCard

	c.Players[c.playerFocus].ShieldCard = c.toReveal[0]
	c.DeckPile.DiscardCard(oldCard)

	c.toReveal = []*Card{}

	return 0, "It's like getting under a blanket on a rainy day !"
}

// @desc: Player at index player puts the next card into his charge slot
func (c *Camaretto) charge() (int, string) {
	var p *Player = c.Players[c.playerFocus]
	if p.ChargeCard == nil {
		var card *Card = c.DeckPile.DrawCard()
		p.Charge(card)
	}

	return 0, "Loading up !"
}

// @desc: Player at index player heals himself
func (c *Camaretto) heal() (int, string) {
	var oldCard *Card = c.Players[c.playerFocus].Heal(c.cardFocus)
	c.DeckPile.DiscardCard(oldCard)

	return 0, "I feel a lil' bit tired, anyone has a vitamin ?"
}

func (c *Camaretto) reveal() {
	if c.state == ATTACK {
		c.toReveal = append(c.toReveal, c.DeckPile.DrawCard())
		var p *Player = c.Players[c.playerTurn]
		if p.ChargeCard != nil {
			c.toReveal = append(c.toReveal, p.ChargeCard)
			p.ChargeCard = nil
		}
	} else if c.state == SHIELD {
		c.toReveal = append(c.toReveal, c.DeckPile.DrawCard())
	}

	// var pers *Character = c.Players[c.playerTurn].Persona
	// var msg string = pers.Talk(c.state)
	c.info.SetMessage("Display this message please, it could be good for testing purposes")
}

/************ *************************************************************************** ************/
/************ ********************************** UPDATE ********************************* ************/
/************ *************************************************************************** ************/

func (c *Camaretto) onReveal(x, y float64) int {
	for i, card := range c.toReveal {
		if card.SSprite.In(x, y) { return i }
	}
	return -1
}

func (c *Camaretto) onHealth(x, y float64) int {
	var p *Player = c.Players[c.playerFocus]
	if p.HealthCard[0] != nil && p.HealthCard[0].SSprite.In(x, y) {
		return 0
	} else if p.HealthCard[1] != nil && p.HealthCard[1].SSprite.In(x, y) {
		return 1
	}
	return -1
}

func (c *Camaretto) onPlayer(x, y float64) int {
	for i, player := range c.Players {
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

func (c *Camaretto) Hover(x, y float64) {
	var speed float64 = 15

	var s *view.Sprite = nil
	if c.state == SET {
		if c.attackButton.SSprite.In(x, y) {
			s = c.attackButton.SSprite
		} else if c.shieldButton.SSprite.In(x, y) {
			s = c.shieldButton.SSprite
		} else if c.chargeButton.SSprite.In(x, y) {
			s = c.chargeButton.SSprite
		} else if c.healButton.SSprite.In(x, y) {
			s = c.healButton.SSprite
		}

		if s != nil {
			var x, y, _ float64 = s.GetCenter()
			c.cursor.Move(x - (s.Width/2), y, speed)
			c.cursor.Rotate(math.Pi/2, speed)
			c.cursor.MoveOffset(0, 0, speed)
			c.cursor.RotateOffset(0, speed)
		}
	} else {
		if c.focus == PLAYER {
			var i int = c.onPlayer(x, y)
			if i != -1 {
				s = c.Players[i].NameSprite
				var x, y, r float64 = s.GetCenter()
				c.cursor.Move(x, y, speed)
				c.cursor.Rotate(math.Pi, speed)
				x, y, r = s.GetOffset()
				c.cursor.MoveOffset(x, y - float64(view.CardHeight*5/2), speed)
				c.cursor.RotateOffset(r, speed)
			}
		} else if c.focus == CARD {
			var i int = c.onHealth(x, y)
			if i != -1 {
				s = c.Players[c.playerFocus].HealthCard[i].SSprite
				var x, y, r float64 = s.GetCenter()
				c.cursor.Move(x, y, speed)
				c.cursor.Rotate(math.Pi, speed)
				x, y, r = s.GetOffset()
				c.cursor.MoveOffset(x, y - float64(view.CardHeight/2), speed)
				c.cursor.RotateOffset(r, speed)
			}
		} else if c.focus == REVEAL {
			for _, card := range c.toReveal {
				if card.SSprite.In(x, y) { s = card.SSprite }
			}

			if s != nil {
				var x, y, r float64 = s.GetCenter()
				c.cursor.Move(x, y + (s.Height/2), speed)
				c.cursor.Rotate(r, speed)
				c.cursor.MoveOffset(0, 0, speed)
				c.cursor.RotateOffset(0, speed)
			}
		}
	}
}

func (c *Camaretto) mousePress(e *event.MouseEvent) {
	if c.state == SET {
		if c.attackButton.SSprite.In(e.X, e.Y) {
			c.attackButton.Pressed()
		} else if c.shieldButton.SSprite.In(e.X, e.Y) {
			c.shieldButton.Pressed()
		} else if c.chargeButton.SSprite.In(e.X, e.Y) {
			c.chargeButton.Pressed()
		} else if c.healButton.SSprite.In(e.X, e.Y) {
			c.healButton.Pressed()
		}
	}
}

func (c *Camaretto) mouseRelease(e *event.MouseEvent) {
	c.cursor.SetCenter(-c.cursor.Width, -c.cursor.Height, 0)
	c.cursor.SetOffset(0, 0, 0)

	c.attackButton.Released()
	c.shieldButton.Released()
	c.chargeButton.Released()
	c.healButton.Released()

	if c.state == SET {
		if c.attackButton.SSprite.In(e.X, e.Y) {
			c.state = ATTACK
			c.focus = PLAYER
		} else if c.shieldButton.SSprite.In(e.X, e.Y) {
			c.state = SHIELD
			c.focus = PLAYER
		} else if c.chargeButton.SSprite.In(e.X, e.Y) {
			c.state = CHARGE
			c.playerFocus = c.playerTurn
			c.focus = COMPLETE
		} else if c.healButton.SSprite.In(e.X, e.Y) {
			c.state = HEAL
			c.playerFocus = c.playerTurn
			c.focus = CARD
		}
	} else {
		if c.focus == PLAYER {
			var i int = c.onPlayer(e.X, e.Y)
			if i != -1 {
				if c.state == ATTACK {
					c.playerFocus = i
					c.focus = CARD
				} else if c.state == SHIELD {
					c.playerFocus = i
					c.reveal()
					c.focus = REVEAL
				}
			}
		} else if c.focus == CARD {
			var i int = c.onHealth(e.X, e.Y)
			if i != -1 {
				c.cardFocus = i
				c.reveal()
				c.focus = REVEAL
			}
		} else if c.focus == REVEAL {
			var i int = c.onReveal(e.X, e.Y)
			if i != -1 { c.toReveal[i].Reveal() }
		}
	}
}

func (c *Camaretto) EventUpdate(e *event.MouseEvent) {
	if e.Event == event.PRESSED {
		c.mousePress(e)
	} else if e.Event == event.RELEASED {
		c.mouseRelease(e)
	}
}

func (c *Camaretto) Update() {
	if c.state == SET {
	}

	if c.focus == REVEAL {
		if len(c.toReveal) == 0 {
			c.focus = COMPLETE
		} else {
			var done bool = true
			for _, card := range c.toReveal { done = done && (!card.Hidden) }

			if done { c.count++ }
			if done && c.count > 33 {
				c.focus = COMPLETE
				c.count = 0
			}
		}
	} else if c.focus == COMPLETE {
		if c.state != SET {
			if c.state == ATTACK {
				c.attack()
			} else if c.state == SHIELD {
				c.shield()
			} else if c.state == CHARGE {
				c.charge()
			} else if c.state == HEAL {
				c.heal()
			}
			c.endTurn()
		}
	}
}

/************ *************************************************************************** ************/
/************ ********************************** RENDER ********************************* ************/
/************ *************************************************************************** ************/

func (c *Camaretto) getPlayerGeoM(i int) (float64, float64, float64) {
	var nbPlayers int = len(c.Players)
	var angleStep float64 = 2*math.Pi / float64(nbPlayers)
	var radius float64 = 200

	var theta float64 = angleStep * float64(i)
	var x float64 = radius * math.Cos(theta + math.Pi/2)
	var y float64 = radius * math.Sin(theta + math.Pi/2)

	return x, y, theta
}

func (c *Camaretto) Render(dst *ebiten.Image, width, height float64) {
	var persona *Character = c.Players[c.playerTurn].Persona
	if c.info.Finished() {
		persona.Talking = false
	} else {
		persona.Talking = true
		c.info.Render()
	}
	persona.Render()

	persona.SSprite.Display(dst)
	c.info.SSprite.Display(dst)

	var buttonXPos float64 = 0
	var buttonYPos float64 = float64(height)*9/10

	if c.state == SET {
		buttonXPos = (float64(width) * 1/4) + (float64(view.ButtonWidth)/2)
		c.attackButton.SSprite.SetCenter(buttonXPos, buttonYPos, 0)
		c.attackButton.SSprite.Display(dst)

		buttonXPos = (float64(width) * 2/4) + (float64(view.ButtonWidth)/2)
		c.shieldButton.SSprite.SetCenter(buttonXPos, buttonYPos, 0)
		c.shieldButton.SSprite.Display(dst)

		buttonXPos = (float64(width) * 3/4) + (float64(view.ButtonWidth)/2)

		if c.Players[c.playerTurn].ChargeCard == nil {
			c.healButton.SSprite.SetCenter(0, 0, 0)
			c.healButton.SSprite.SetOffset(0, 0, 0)

			c.chargeButton.SSprite.SetCenter(buttonXPos, buttonYPos, 0)
			c.chargeButton.SSprite.Display(dst)
		} else {
			c.chargeButton.SSprite.SetCenter(0, 0, 0)
			c.chargeButton.SSprite.SetOffset(0, 0, 0)

			c.healButton.SSprite.SetCenter(buttonXPos, buttonYPos, 0)
			c.healButton.SSprite.Display(dst)
		}
	}

	var centerX float64 = width/2
	var centerY float64 = (height * 6/8)/2

	for i, player := range c.Players {
		var x, y, theta float64 = c.getPlayerGeoM(i)
		player.RenderCards(dst, centerX + x, centerY + y, theta)
	}

	if len(c.toReveal) > 0 {
		for i, card := range c.toReveal {
			var s *view.Sprite = card.SSprite
			var x float64 = (float64(i) - float64(len(c.toReveal)-1)/2) * float64(view.CardWidth)
			s.Move(centerX + x, centerY, 0.5)
			card.SSprite.Rotate(0, 0.2)
			s.MoveOffset(0, 0, 0.2)
			card.SSprite.RotateOffset(0, 0.2)
			card.SSprite.Display(dst)
		}
	} else {
		c.DeckPile.Render(dst, centerX, centerY)
	}

	c.cursor.Display(dst)
}
