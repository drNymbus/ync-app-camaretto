package game

import (
	"log"
)

const (
	MaxNbPlayers int = 6
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

type Action struct {
	State GameState
	Focus FocusState

	PlayerTurn int
	PlayerFocus int
	CardFocus int
}

func NewAction(i int) *Action {
	var a *Action = &Action{}
	
	a.State = SET
	a.Focus = NONE

	a.PlayerTurn = i
	a.PlayerFocus = -1
	a.CardFocus = -1

	return a
}

type Camaretto struct {
	Log []*Action
	Current *Action

	NbPlayers int
	Players []*Player

	DeckPile *Deck

	ToReveal []*Card
}

// @desc: Initialize attributes of a Camaretto instance, given the number of players: n
func (c *Camaretto) Init(seed int64, n int, names []string) {
	if len(names) != n { log.Fatal("[Camaretto.Init] You finna start a game like that ?!") }

	c.Log = []*Action{}
	c.Current = NewAction(0)

	c.NbPlayers = n
	c.Players = make([]*Player, n)

	for i, _ := range make([]int, n) { // Init players
		var name string = names[i%len(names)]
		var char *Character = NewCharacter(name)
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
}

/************ ***************************************************************************** ************/
/************ ********************************** ACTIONS ********************************** ************/
/************ ***************************************************************************** ************/

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
func (c *Camaretto) EndTurn() {
	c.Log = append(c.Log, c.Current)

	var newTurn int = (c.Current.PlayerTurn+1) % c.NbPlayers
	for ;c.Players[newTurn].Dead; { newTurn = (newTurn+1) % c.NbPlayers }
	
	c.Current = NewAction(newTurn)
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
				amount = amount - dst.HealthCard[c.Current.CardFocus].Value
				c.DeckPile.DiscardCard(dst.HealthCard[c.Current.CardFocus])
				dst.HealthCard[c.Current.CardFocus] = nil
				// Joker's gone, we replace the health you focused
				healthSlot = c.Current.CardFocus
			}
	
			// Wow that's a really big hit
			if amount > 0 && dst.HealthCard[1-c.Current.CardFocus] != nil {
				amount = amount - dst.HealthCard[1-c.Current.CardFocus].Value
				c.DeckPile.DiscardCard(dst.HealthCard[1-c.Current.CardFocus])
				dst.HealthCard[1-c.Current.CardFocus] = nil
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
func (c *Camaretto) Attack() {
	var atkCard *Card = c.ToReveal[0]
	var chargeCard *Card = nil
	if len(c.ToReveal) == 2 { chargeCard = c.ToReveal[1] }

	var atkValue int = atkCard.Value
	if chargeCard != nil { atkValue = atkValue + chargeCard.Value }

	c.attackPlayer(c.Players[c.Current.PlayerFocus], atkValue)

	c.DeckPile.DiscardCard(atkCard)
	if chargeCard != nil { c.DeckPile.DiscardCard(chargeCard) }
	c.ToReveal = []*Card{}
}

// @desc: Player at index player gets assigned a new shield
func (c *Camaretto) Shield() {
	var oldCard *Card = c.Players[c.Current.PlayerFocus].ShieldCard

	c.Players[c.Current.PlayerFocus].ShieldCard = c.ToReveal[0]
	c.DeckPile.DiscardCard(oldCard)

	c.ToReveal = []*Card{}
}

// @desc: Player at index player puts the next card into his charge slot
func (c *Camaretto) Charge() {
	var p *Player = c.Players[c.Current.PlayerFocus]
	if p.ChargeCard == nil {
		var card *Card = c.DeckPile.DrawCard()
		p.Charge(card)
	}
}

// @desc: Player at index player heals himself
func (c *Camaretto) Heal() {
	var oldCard *Card = c.Players[c.Current.PlayerFocus].Heal(c.Current.CardFocus)
	c.DeckPile.DiscardCard(oldCard)
}

// @desc: Place cards to be revealed before the action takes place
func (c *Camaretto) Reveal() {
	if c.Current.State == ATTACK {
		c.ToReveal = append(c.ToReveal, c.DeckPile.DrawCard())
		var p *Player = c.Players[c.Current.PlayerTurn]
		if p.ChargeCard != nil {
			c.ToReveal = append(c.ToReveal, p.ChargeCard)
			p.ChargeCard = nil
		}
	} else if c.Current.State == SHIELD {
		c.ToReveal = append(c.ToReveal, c.DeckPile.DrawCard())
	}
}
