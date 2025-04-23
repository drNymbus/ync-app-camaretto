package game

import (
	"log"

	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Camaretto struct {
	width, height float64

	Log []*Action
	Current *Action
	AlteredState bool

	NbPlayers int
	Players []*Player

	DeckPile *Deck

	ToReveal []*Card

	tick int
}

// @desc: Initialize attributes of a Camaretto instance, given the number of players: n
func (c *Camaretto) Init(seed int64, names []string, w, h float64) {

	c.Log = []*Action{}
	c.Current = NewAction(0)
	c.AlteredState = false

	c.NbPlayers = len(names)
	c.Players = []*Player{}

	var angleStep float64 = 2*math.Pi / float64(c.NbPlayers)
	var radius float64 = 200

	for i, name := range names { // Init players
		var char *Character = NewCharacter(name)

		var theta float64 = angleStep * float64(i)
		var x float64 = w/2 + radius * math.Cos(theta + math.Pi/2)
		var y float64 = h/2 + radius * math.Sin(theta + math.Pi/2)

		c.Players = append(c.Players, NewPlayer(name, char, x, y, theta))
	}

	c.DeckPile = &Deck{}
	c.DeckPile.Init(seed, w/2, h/2)

	for i, _ := range make([]int, c.NbPlayers*2) { // Init Health
		var card *Card = c.DeckPile.DrawCard()
		if card.Name == "Joker" {
			c.DeckPile.DiscardCard(card)
			card = c.DeckPile.DrawCard()
			if card.Name == "Joker" {
				log.Fatal("[Camaretto.Init] (Health) 2 JOKERS IN A ROW ?! What were the chances anyway...")
			}

			card.Reveal()
			c.Players[i%c.NbPlayers].SetJokerHealth(card)

			card = c.DeckPile.DrawCard()
			if card.Name == "Joker" {
				log.Fatal("[Camaretto.Init] (Health) 2 JOKERS SPACED BY ONE CARD ?! Ok this one might not be THAT crazy but still...")
			}
		}

		card.Reveal()
		c.Players[i%c.NbPlayers].SetHealth(card, int(i/c.NbPlayers))
	}

	for i, _ := range make([]int, c.NbPlayers) { // Init Shield
		var card *Card = c.DeckPile.DrawCard()
		card.Reveal()
		if card.Name == "Joker" {
			c.Players[i].SetJokerShield(card)
		} else {
			c.Players[i].SetShield(card)
		}
	}
}

// @desc: Returns true if one player is left, false otherwise
func (c *Camaretto) IsGameOver() bool {
	var count int = c.NbPlayers
	for _, p := range c.Players {
		if p.Dead { count-- }
	}

	if count > 1 { return false }
	return true
}

// @desc: Finish turn reset game state and pass onto the next player's turn
func (c *Camaretto) endTurn() {
	var newTurn int = (c.Current.PlayerTurn+1) % c.NbPlayers
	for ;c.Players[newTurn].Dead; { newTurn = (newTurn+1) % c.NbPlayers }
	
	for _, p := range c.Players {
		p.Trigger = nil
		for _, card := range p.Health {
			if card != nil { card.Trigger = nil }
		}
	}

	c.Log = append(c.Log, c.Current)
	c.Current = NewAction(newTurn)
}

/************ ***************************************************************************** ************/
/************ ********************************** ACTIONS ********************************** ************/
/************ ***************************************************************************** ************/

// @desc: Compute which cards are lost when a player (dst) is attacked by "amount" of points
func (c *Camaretto) attackPlayer(player *Player, amount int) {
	var jokerShield *Card = player.SetJokerShield(nil)
	if jokerShield != nil { // Joker shield popped
		c.DeckPile.DiscardCard(jokerShield)
	} else {
		// In which health slot should we put a new card ?
		var healthSlot int = -1 // -1: none; 0: health[0]; 1:health[1]; 2: joker

		var shield *Card = player.SetShield(nil)
		if shield != nil {
			amount = amount - shield.Value
		}

		if amount > 0 {
			player.SetShield(shield) // Put back this lil' shield of yours
			c.DeckPile.DiscardCard(player.SetCharge(nil)) // Sorry, but if you had a charge it's gone

			// Do we have a joker health ? Then it's tanking (wether you like it or not)
			var jokerHealth *Card = player.SetJokerHealth(nil)
			if jokerHealth != nil {
				amount = amount - jokerHealth.Value
				c.DeckPile.DiscardCard(jokerHealth)
				// We have to replace your jokerHealth then
				healthSlot = 2
			}
	
			// Is the attack still going ?
			if amount > 0 {
				var health *Card = player.SetHealth(nil, c.Current.CardFocus)
				if health != nil { // Why this test should not pass ? I mean ...
					amount = amount - health.Value
					c.DeckPile.DiscardCard(health)
					// Joker's gone, we replace the health you focused
					healthSlot = c.Current.CardFocus
				}
			}
	
			// Wow that's a really big hit
			if amount > 0 {
				var health *Card = player.SetHealth(nil, 1 - c.Current.CardFocus)
				if health != nil { // You really have nothing going on for you huh ?
					amount = amount - health.Value
					c.DeckPile.DiscardCard(health)
					// Both of your health cards took a hit ? Guess you don't have an option anymore
					healthSlot = 0
				}
			}

			// R.I.P in Peperonni
			if amount >= 0 {
				// Give all your cards to your little friends pls
				c.DeckPile.DiscardCard(player.SetHealth(nil, 0))
				c.DeckPile.DiscardCard(player.SetHealth(nil, 1))
				c.DeckPile.DiscardCard(player.SetJokerHealth(nil))
				c.DeckPile.DiscardCard(player.SetShield(nil))
				c.DeckPile.DiscardCard(player.SetJokerShield(nil))
				player.Dead = true
			} else { // Hang tight, we're gonna fix you sport !
				amount = amount * -1

				var newHealthCard *Card = nil
				newHealthCard = c.DeckPile.FindInDiscardPile(amount)
				if newHealthCard == nil {
					newHealthCard = c.DeckPile.FindInDrawPile(amount)
				}
				if newHealthCard == nil { log.Fatal("[Camaretto.Attack] Could not find a card with health points left") }

				newHealthCard.Reveal()

				if healthSlot == 2 {
					player.SetJokerHealth(newHealthCard)
				} else {
					player.SetHealth(newHealthCard, healthSlot)
				}
			}
		} else if amount == 0 { // Poof ! Shield's gone !
			c.DeckPile.DiscardCard(shield)
		}
	}
}

// @desc: Player at index src attacks the player at index dst
func (c *Camaretto) attack() {
	var atk int = 0
	for _, card := range c.ToReveal {
		atk += card.Value
		c.DeckPile.DiscardCard(card)
	}
	c.ToReveal = []*Card{}

	c.attackPlayer(c.Players[c.Current.PlayerFocus], atk)
}

// @desc: Player at index player gets assigned a new shield
func (c *Camaretto) shield() {
	var player *Player = c.Players[c.Current.PlayerFocus]

	var old *Card
	if c.ToReveal[0].Name == "Joker" {
		old = player.SetJokerShield(c.ToReveal[0])
	} else {
		old = player.SetShield(c.ToReveal[0])
	}

	c.ToReveal = []*Card{}
	c.DeckPile.DiscardCard(old)
}

// @desc: Player at index player puts the next card into his charge slot
func (c *Camaretto) charge() {
	var player *Player = c.Players[c.Current.PlayerFocus]
	if player.IsChargeEmpty() {
		var card *Card = c.DeckPile.DrawCard()
		player.SetCharge(card)
	}
}

// @desc: Player at index player heals himself
func (c *Camaretto) heal() {
	var player *Player = c.Players[c.Current.PlayerFocus]
	if !player.IsChargeEmpty() {
		var charge *Card = player.SetCharge(nil)
		charge.Reveal()
		var old *Card = player.SetHealth(charge, c.Current.CardFocus)
		c.DeckPile.DiscardCard(old)
	}
}

func (c *Camaretto) addCardToReveal(card *Card) {
	c.Current.Reveal = append(c.Current.Reveal, false)
	c.ToReveal = append(c.ToReveal, card)

	card.Trigger = func() {
		c.Current.Reveal[len(c.ToReveal)-1] = true
		card.Reveal()
	}
	
	for i, reveal := range c.ToReveal {
		var x, y, rOff float64 = c.Players[c.Current.PlayerTurn].GetPosition()

		reveal.SSprite.Move(x, y, 1)
		reveal.SSprite.Rotate(0, 1)
		reveal.SSprite.RotateOffset(rOff, 1)

		var iOff int = i - len(c.ToReveal)/2
		var xOff float64 = reveal.SSprite.Width * float64(iOff)
		var yOff float64 = -reveal.SSprite.Height * 3/2
		reveal.SSprite.MoveOffset(xOff, yOff, 1)
	}
}

// @desc: Place cards to be revealed before the action takes place
func (c *Camaretto) reveal() {
	for _, player := range c.Players {
		player.Trigger = nil
		for _, h := range player.Health { h.Trigger = nil }
	}

	if c.Current.State == ATTACK {
		c.addCardToReveal(c.DeckPile.DrawCard())

		var player *Player = c.Players[c.Current.PlayerTurn]
		if !player.IsChargeEmpty() {
			c.addCardToReveal(player.SetCharge(nil))
		}
	} else if c.Current.State == SHIELD {
		c.addCardToReveal(c.DeckPile.DrawCard())
	}
}

func (c *Camaretto) Update() error {
	for _, player := range c.Players {
		player.Update()
	}

	c.DeckPile.Update()

	for _, card := range c.ToReveal {
		card.Update()
	}

	c.AlteredState = false
	if c.Current.State != SET {
		if c.Current.Focus == NONE {
			if c.Current.State == ATTACK {
				c.Current.Focus = PLAYER

				c.AlteredState = true
			} else if c.Current.State == SHIELD {
				c.Current.Focus = PLAYER

				c.AlteredState = true
			} else if c.Current.State == CHARGE {
				c.Current.Focus = COMPLETE
				c.Current.PlayerFocus = c.Current.PlayerTurn

				c.AlteredState = true
			} else if c.Current.State == HEAL {
				c.Current.Focus = CARD
				c.Current.PlayerFocus = c.Current.PlayerTurn

				c.AlteredState = true
			}
		} else if c.Current.Focus == PLAYER {
			if c.Current.PlayerFocus != -1 {
				if c.Current.State == ATTACK {
					c.Current.Focus = CARD
					var player *Player = c.Players[c.Current.PlayerFocus]
					for i, health := range player.Health {
						if health != nil { health.Trigger = func() { c.Current.CardFocus = i } }
					}

					c.AlteredState = true
				} else if c.Current.State == SHIELD {
					c.Current.Focus = REVEAL
					c.reveal()

					c.AlteredState = true
				}
			} else {
				for i, player := range c.Players {
					if c.Current.State != ATTACK || c.Current.PlayerTurn != i {
						player.Trigger = func() { c.Current.PlayerFocus = i }
					}
				}
			}
		} else if c.Current.Focus == CARD {
			if c.Current.CardFocus != -1 {
				c.Current.Focus = REVEAL
				c.reveal()

				c.AlteredState = true
			} else {
				var player *Player = c.Players[c.Current.PlayerFocus]
				for i, health := range player.Health {
					health.Trigger = func() { c.Current.CardFocus = i }
				}
			}
		} else if c.Current.Focus == REVEAL {
			if len(c.ToReveal) > 0 {
				var done bool = true
				for _, revealed := range c.Current.Reveal { done = done && revealed }
				if done { c.Current.Focus = COMPLETE }

				c.AlteredState = true
			}

		} else if c.Current.Focus == COMPLETE {
			switch ;c.Current.State {
				case ATTACK: c.attack()
				case SHIELD: c.shield()
				case CHARGE: c.charge()
				case HEAL: c.heal()
			}
			c.endTurn()

			c.AlteredState = true
		}
	}

	if c.AlteredState { log.Println("ALTER", c.Current) }

	return nil
}

func (c *Camaretto) Draw(screen *ebiten.Image) {
	for _, player := range c.Players {
		player.Draw(screen)
	}

	c.DeckPile.Draw(screen)

	for _, card := range c.ToReveal {
		card.Draw(screen)
	}
}
