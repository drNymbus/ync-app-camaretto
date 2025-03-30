package game

import (
	"log"
	"math"

	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"camaretto/model/component"
	"camaretto/view"
)

type GameState int
const (
	SET GameState = iota
	ATTACK
	SHIELD
	CHARGE
	HEAL
)

type FocusState int
const (
	NONE FocusState = iota
	PLAYER
	CARD
	REVEAL
	COMPLETE
)

type Camaretto struct {
	width, height float64

	State GameState
	Focus FocusState

	PlayerTurn int
	PlayerFocus int
	CardFocus int

	nbPlayers int
	Players []*Player
	DeckPile *Deck

	ToReveal []*Card

	attackButton *component.Button
	shieldButton *component.Button
	chargeButton *component.Button
	healButton *component.Button

	cursor *view.Sprite

	info *component.TextBox

	count int
}

// @desc: Initialize attributes of a Camaretto instance, given the number of players: n
func (c *Camaretto) Init(n int, names []string, seed int64, w, h int) {
	if len(names) != n { log.Fatal("[Camaretto.Init] You finna start a game like that ?!") }

	c.width, c.height = float64(w), float64(h)

	c.State = SET
	c.Focus = NONE

	c.PlayerTurn = 0
	c.PlayerFocus = -1
	c.CardFocus = -1

	c.info = component.NewTextBox(c.width - 50, c.height*1/5 + 30, "", color.RGBA{0, 0, 0, 255}, color.RGBA{0, 51, 153, 127})
	var x, y float64 = c.width/2, c.height*8/10 + 65
	c.info.SSprite.SetCenter(x, y, 0)

	c.nbPlayers = n
	c.Players = make([]*Player, n)

	for i, _ := range make([]int, n) { // Init players
		var name string = names[i%len(names)]
		var char *Character = NewCharacter(name)
		var bodyX float64 = (x - c.info.SSprite.Width/2) + char.SSprite.Width/2
		var bodyY float64 = (y + c.info.SSprite.Height/2) - char.SSprite.Height/2
		char.SSprite.SetCenter(bodyX, bodyY, 0)

		c.Players[i] = NewPlayer(name, char)
	}

	c.DeckPile = &Deck{}
	c.DeckPile.Init(seed)

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

	var buttonXPos float64 = 0
	var buttonYPos float64 = c.height * 9/10

	c.attackButton = component.NewButton("ATTACK", color.RGBA{0, 0, 0, 255}, "RED")
	buttonXPos = (c.width * 1/4) + (float64(view.ButtonWidth)/2)
	c.attackButton.SSprite.SetCenter(buttonXPos, buttonYPos, 0)

	c.shieldButton = component.NewButton("SHIELD", color.RGBA{0, 0, 0, 255}, "BLUE")
	buttonXPos = (c.width * 2/4) + (float64(view.ButtonWidth)/2)
	c.shieldButton.SSprite.SetCenter(buttonXPos, buttonYPos, 0)

	buttonXPos = (c.width * 3/4) + (float64(view.ButtonWidth)/2)

	c.chargeButton = component.NewButton("CHARGE", color.RGBA{0, 0, 0, 255}, "YELLOW")
	c.chargeButton.SSprite.SetCenter(buttonXPos, buttonYPos, 0)

	c.healButton = component.NewButton("HEAL", color.RGBA{0, 0, 0, 255}, "GREEN")
	c.healButton.SSprite.SetCenter(buttonXPos, buttonYPos, 0)

	c.cursor = view.NewSprite(view.LoadCursorImage(), false, color.RGBA{0, 0, 0, 0}, nil)
	c.cursor.SetCenter(-c.cursor.Width, -c.cursor.Height, 0)
	c.cursor.SetOffset(0, 0, 0)

	c.count = 0
}

/************ ***************************************************************************** ************/
/************ ********************************** ACTIONS ********************************** ************/
/************ ***************************************************************************** ************/

// @desc: Finish turn reset game state and pass onto the next player's turn
func (c *Camaretto) endTurn() {
	c.State = SET
	c.Focus = NONE
	c.PlayerFocus = -1
	c.CardFocus = -1

	c.PlayerTurn = (c.PlayerTurn+1) % c.nbPlayers
	for ;c.Players[c.PlayerTurn].Dead; { c.PlayerTurn = (c.PlayerTurn+1) % c.nbPlayers }
}

// @desc: Compute which cards are lost when a player (dst) is attacked by "amount" of points
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
				amount = amount - dst.HealthCard[c.CardFocus].Value
				c.DeckPile.DiscardCard(dst.HealthCard[c.CardFocus])
				dst.HealthCard[c.CardFocus] = nil
				// Joker's gone, we replace the health you focused
				healthSlot = c.CardFocus
			}
	
			// Wow that's a really big hit
			if amount > 0 && dst.HealthCard[1-c.CardFocus] != nil {
				amount = amount - dst.HealthCard[1-c.CardFocus].Value
				c.DeckPile.DiscardCard(dst.HealthCard[1-c.CardFocus])
				dst.HealthCard[1-c.CardFocus] = nil
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
func (c *Camaretto) attack() {
	var atkCard *Card = c.ToReveal[0]
	var chargeCard *Card = nil
	if len(c.ToReveal) == 2 { chargeCard = c.ToReveal[1] }

	var atkValue int = atkCard.Value
	if chargeCard != nil { atkValue = atkValue + chargeCard.Value }

	c.attackPlayer(c.Players[c.PlayerFocus], atkValue)

	c.DeckPile.DiscardCard(atkCard)
	if chargeCard != nil { c.DeckPile.DiscardCard(chargeCard) }
	c.ToReveal = []*Card{}
}

// @desc: Player at index player gets assigned a new shield
func (c *Camaretto) shield() {
	var oldCard *Card = c.Players[c.PlayerFocus].ShieldCard

	c.Players[c.PlayerFocus].ShieldCard = c.ToReveal[0]
	c.DeckPile.DiscardCard(oldCard)

	c.ToReveal = []*Card{}
}

// @desc: Player at index player puts the next card into his charge slot
func (c *Camaretto) charge() {
	var p *Player = c.Players[c.PlayerFocus]
	if p.ChargeCard == nil {
		var card *Card = c.DeckPile.DrawCard()
		p.Charge(card)
	}
}

// @desc: Player at index player heals himself
func (c *Camaretto) heal() {
	var oldCard *Card = c.Players[c.PlayerFocus].Heal(c.CardFocus)
	c.DeckPile.DiscardCard(oldCard)
}

// @desc: Place cards to be revealed before the action takes place
func (c *Camaretto) reveal() {
	if c.State == ATTACK {
		c.ToReveal = append(c.ToReveal, c.DeckPile.DrawCard())
		var p *Player = c.Players[c.PlayerTurn]
		if p.ChargeCard != nil {
			c.ToReveal = append(c.ToReveal, p.ChargeCard)
			p.ChargeCard = nil
		}
	} else if c.State == SHIELD {
		c.ToReveal = append(c.ToReveal, c.DeckPile.DrawCard())
	}

	// var pers *Character = c.Players[c.PlayerTurn].Persona
	// var msg string = pers.Talk(c.State)
	c.info.SetMessage("Display this message please, it could be good for testing purposes")
}

/************ *************************************************************************** ************/
/************ ********************************** UPDATE ********************************* ************/
/************ *************************************************************************** ************/

// @desc: Returns true if one player is left, false otherwise
func (c *Camaretto) IsGameOver() bool {
	var count int = c.nbPlayers
	for _, p := range c.Players {
		if p.Dead { count-- }
	}

	if count > 1 { return false }
	return true
}

// @desc: If the coordinates (x,y) are on any of the cards to be revealed returns the card's index, returns -1 otherwise
func (c *Camaretto) onReveal(x, y float64) int {
	for i, card := range c.ToReveal {
		if card.SSprite.In(x, y) { return i }
	}
	return -1
}

// @desc: If the coordinates (x,y) are on any of the focused player's health card return the card index, returns -1 otherwise
func (c *Camaretto) onHealth(x, y float64) int {
	var p *Player = c.Players[c.PlayerFocus]
	if p.HealthCard[0] != nil && p.HealthCard[0].SSprite.In(x, y) {
		return 0
	} else if p.HealthCard[1] != nil && p.HealthCard[1].SSprite.In(x, y) {
		return 1
	}
	return -1
}

// @desc: If the coordinates (x,y) are on any player's card the player's index is returned, -1 otherwise
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

// @desc: Update elements upon mouse hover on coordinates (x,y)
func (c *Camaretto) Hover(x, y float64) {
	var speed float64 = 15

	var s *view.Sprite = nil
	if c.State == SET {
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
		if c.Focus == PLAYER {
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
		} else if c.Focus == CARD {
			var i int = c.onHealth(x, y)
			if i != -1 {
				s = c.Players[c.PlayerFocus].HealthCard[i].SSprite
				var x, y, r float64 = s.GetCenter()
				c.cursor.Move(x, y, speed)
				c.cursor.Rotate(math.Pi, speed)
				x, y, r = s.GetOffset()
				c.cursor.MoveOffset(x, y - float64(view.CardHeight/2), speed)
				c.cursor.RotateOffset(r, speed)
			}
		} else if c.Focus == REVEAL {
			for _, card := range c.ToReveal {
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

// @desc: Update elements on mouse button press on coordinates (x,y)
func (c *Camaretto) MousePress(x, y float64) {
	if c.State == SET {
		if c.attackButton.SSprite.In(x, y) {
			c.attackButton.Pressed()
		} else if c.shieldButton.SSprite.In(x, y) {
			c.shieldButton.Pressed()
		} else {
			if c.Players[c.PlayerTurn].ChargeCard == nil && c.chargeButton.SSprite.In(x, y) {
				c.chargeButton.Pressed()
			} else if c.healButton.SSprite.In(x, y) {
				c.healButton.Pressed()
			}
		}
	}
}

// @desc: Update elements on mouse button release on coordinates (x,y)
func (c *Camaretto) MouseRelease(x, y float64) {
	c.cursor.SetCenter(-c.cursor.Width, -c.cursor.Height, 0)
	c.cursor.SetOffset(0, 0, 0)

	c.attackButton.Released()
	c.shieldButton.Released()
	c.chargeButton.Released()
	c.healButton.Released()

	if c.State == SET {
		if c.attackButton.SSprite.In(x, y) {
			c.State = ATTACK
			c.Focus = PLAYER
		} else if c.shieldButton.SSprite.In(x, y) {
			c.State = SHIELD
			c.Focus = PLAYER
		} else {
			if c.Players[c.PlayerTurn].ChargeCard == nil && c.chargeButton.SSprite.In(x, y) {
				c.State = CHARGE
				c.PlayerFocus = c.PlayerTurn
				c.Focus = COMPLETE
			} else if c.healButton.SSprite.In(x, y) {
				c.State = HEAL
				c.PlayerFocus = c.PlayerTurn
				c.Focus = CARD
			}
		}
	} else {
		if c.Focus == PLAYER {
			var i int = c.onPlayer(x, y)
			if i != -1 {
				if c.State == ATTACK {
					c.PlayerFocus = i
					c.Focus = CARD
				} else if c.State == SHIELD {
					c.PlayerFocus = i
					c.reveal()
					c.Focus = REVEAL
				}
			}
		} else if c.Focus == CARD {
			var i int = c.onHealth(x, y)
			if i != -1 {
				c.CardFocus = i
				c.reveal()
				c.Focus = REVEAL
			}
		} else if c.Focus == REVEAL {
			var i int = c.onReveal(x, y)
			if i != -1 { c.ToReveal[i].Reveal() }
		}
	}
}

// @desc: Complete actions that do not require user input
func (c *Camaretto) Update() {
	if c.Focus == REVEAL {
		if len(c.ToReveal) == 0 {
			c.Focus = COMPLETE
		} else {
			var done bool = true
			for _, card := range c.ToReveal { done = done && (!card.Hidden) }

			if done { c.count++ }
			if done && c.count > 33 {
				c.Focus = COMPLETE
				c.count = 0
			}
		}
	} else if c.Focus == COMPLETE {
		if c.State != SET {
			if c.State == ATTACK {
				c.attack()
			} else if c.State == SHIELD {
				c.shield()
			} else if c.State == CHARGE {
				c.charge()
			} else if c.State == HEAL {
				c.heal()
			}
			c.endTurn()
		}
	}
}

/************ *************************************************************************** ************/
/************ ********************************** RENDER ********************************* ************/
/************ *************************************************************************** ************/

// @desc: Turn index into a coordinate (x,y) and an angle (theta)
func (c *Camaretto) getPlayerGeoM(i int) (float64, float64, float64) {
	var nbPlayers int = len(c.Players)
	var angleStep float64 = 2*math.Pi / float64(nbPlayers)
	var radius float64 = 200

	var theta float64 = angleStep * float64(i)
	var x float64 = radius * math.Cos(theta + math.Pi/2)
	var y float64 = radius * math.Sin(theta + math.Pi/2)

	return x, y, theta
}

// @desc: Render all elements on a given image (dst)
func (c *Camaretto) Display(dst *ebiten.Image) {
	var persona *Character = c.Players[c.PlayerTurn].Persona
	if c.info.Finished() {
		persona.Talking = false
	} else {
		persona.Talking = true
		c.info.Render()
	}
	persona.Render()
	persona.SSprite.Display(dst)
	c.info.SSprite.Display(dst)

	if c.State == SET {
		c.attackButton.SSprite.Display(dst)
		c.shieldButton.SSprite.Display(dst)

		if c.Players[c.PlayerTurn].ChargeCard == nil {
			c.chargeButton.SSprite.Display(dst)
		} else {
			c.healButton.SSprite.Display(dst)
		}
	}

	var centerX float64 = c.width/2
	var centerY float64 = (c.height * 6/8)/2

	for i, player := range c.Players {
		var x, y, theta float64 = c.getPlayerGeoM(i)
		player.RenderCards(dst, centerX + x, centerY + y, theta)
	}

	if len(c.ToReveal) > 0 {
		for i, card := range c.ToReveal {
			var s *view.Sprite = card.SSprite
			var x float64 = (float64(i) - float64(len(c.ToReveal)-1)/2) * float64(view.CardWidth)
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
