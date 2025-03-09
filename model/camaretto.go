package model

import (
	"log"
	"math"

	"github.com/hajimehoshi/ebiten/v2"

	// "camaretto/view"
)

type GameState int
const (
	SET GameState = 0
	ATTACK GameState = 1
	SHIELD GameState = 2
	CHARGE GameState = 3
	HEAL GameState = 4
	REVEAL GameState = 5
	END GameState = 6
)

type FocusState int
const (
	NONE FocusState = 0
	PLAYER FocusState = 1
	CARD FocusState = 2
	COMPLETE FocusState = 3
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
	DrawnCard *Card
}

// @desc: Initialize a new Camaretto instance of the game, then returns a reference to the Camaretto object
// func NewCamaretto(n int, sheet *ebiten.Image, tileWidth int, tileHeight int) *Camaretto {
func NewCamaretto(n int) *Camaretto {
	var c *Camaretto = &Camaretto{}
	// c.Init(n, sheet, tileWidth, tileHeight)
	c.Init(n)
	return c
}

// @desc: Initialize attributes of a Camaretto instance, given the number of players: n
// func (c *Camaretto) Init(n int, sheet *ebiten.Image, tileWidth int, tileHeight int) {
func (c *Camaretto) Init(n int) {
	c.state = SET
	c.focus = NONE

	c.playerTurn = 0
	c.playerFocus = -1
	c.cardFocus = -1

	c.DeckPile = &Deck{}
	c.DeckPile.Init()
	c.DeckPile.ShuffleDrawPile()

	c.DrawnCard = nil

	c.nbPlayers = n
	c.Players = make([]*Player, n)

	var names []string = []string{"Alfred", "Robin", "Parker", "Bruce", "Loïs", "Logan"}
	for i, _ := range make([]int, n) { // Init players
		c.Players[i] = NewPlayer(names[i%len(names)])
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
}

/************ ***************************************************************************** ************/
/************ ********************************** GET/SET ********************************** ************/
/************ ***************************************************************************** ************/

func (c *Camaretto) SetState(s GameState) (int, string) {
	if s == CHARGE && c.Players[c.playerTurn].ChargeCard != nil {
		return 1, "Already a card in charge !"
	}

	if s == HEAL && c.Players[c.playerTurn].ChargeCard == nil {
		return 1, "Cannot heal without a card in charge"
	}

	c.state = s
	return 0, ""
}
func (c *Camaretto) GetState() GameState { return c.state }

func (c *Camaretto) SetFocus(f FocusState) { c.focus = f }
func (c *Camaretto) GetFocus() FocusState { return c.focus }

func (c *Camaretto) SetPlayerFocus(i int)  { c.playerFocus = i }
func (c *Camaretto) GetPlayerFocus() int { return c.playerFocus }

func (c *Camaretto) SetCardFocus(i int)  { 
	c.cardFocus = i
}
func (c *Camaretto) GetCardFocus() int { return c.cardFocus }

func (c *Camaretto) EndTurn() {
	c.state = SET
	c.playerFocus = -1
	c.cardFocus = -1

	c.playerTurn = (c.playerTurn+1) % c.nbPlayers
	for ;c.Players[c.playerTurn].Dead; { c.playerTurn = (c.playerTurn+1) % c.nbPlayers }
}
func (c *Camaretto) GetPlayerTurn() int { return c.playerTurn }

/************ ***************************************************************************** ************/
/************ ********************************** ACTIONS ********************************** ************/
/************ ***************************************************************************** ************/

// @desc: Returns true if one player is left, false otherwise
func (c *Camaretto) IsGameOver() bool {
	var count int = 0
	for _, p := range c.Players {
		if p.Dead { count++ }
	}

	if count > 1 { return true }
	return false
}

// @desc: Player at index src attacks the player at index dst
func (c *Camaretto) Attack() (int, string) {
	var src int = c.playerTurn
	var dst int = c.playerFocus
	var at int = c.cardFocus

	var atkCard *Card = c.DeckPile.DrawCard()
	atkCard.Reveal()

	var atkValue int
	var charge *Card
	atkValue, charge = c.Players[src].Attack(atkCard)

	c.DeckPile.DiscardCard(atkCard)
	if charge != nil { c.DeckPile.DiscardCard(charge) }

	var newHealthValue int
	var joker, health1, health2 *Card
	// health1 is the health card at index "at", health2 is the other health card
	newHealthValue, joker, health1, health2 = c.Players[dst].LoseHealth(atkValue, at)

	var jokerSlot, health1Slot, health2Slot bool = false, false, false
	if joker != nil {
		c.DeckPile.DiscardCard(joker)
		jokerSlot = true
	}

	if health1 != nil {
		c.DeckPile.DiscardCard(health1)
		jokerSlot = false
		health1Slot = true
	}

	if health2 != nil {
		c.DeckPile.DiscardCard(health2)
		health1Slot = false
		health2Slot = true
	}

	// There's health to be recovered
	if newHealthValue > 0 {
		var newHealthCard *Card = nil
		newHealthCard = c.DeckPile.FindInDiscardPile(newHealthValue)
		if newHealthCard == nil {
			newHealthCard = c.DeckPile.FindInDrawPile(newHealthValue)
		}
		if newHealthCard == nil { log.Fatal("[Camaretto.Attack] Could not find a card with health points left") }

		newHealthCard.Reveal()

		if jokerSlot {
			c.Players[dst].JokerHealth = newHealthCard
		} else if health1Slot {
			c.Players[dst].HealthCard[at] = newHealthCard
		} else if health2Slot {
			c.Players[dst].HealthCard[0] = newHealthCard
		}
	} else { // Every card the player had in hand are put back in the DiscardPile
		var p *Player = c.Players[dst]
		if p.HealthCard[0] != nil { c.DeckPile.DiscardCard(p.HealthCard[0]); p.HealthCard[0] = nil }
		if p.HealthCard[1] != nil { c.DeckPile.DiscardCard(p.HealthCard[1]); p.HealthCard[1] = nil }
		if p.JokerHealth != nil { c.DeckPile.DiscardCard(p.JokerHealth); p.JokerHealth = nil }
		if p.ShieldCard != nil { c.DeckPile.DiscardCard(p.ShieldCard); p.ShieldCard = nil }
		if p.JokerShield != nil { c.DeckPile.DiscardCard(p.JokerShield); p.JokerShield = nil }
	}

	return 0, "Attack was great, might do it again ! 5/5"
}

// @desc: Player at index player gets assigned a new shield
func (c *Camaretto) Shield() (int, string) {
	var oldCard *Card = c.Players[c.playerFocus].ShieldCard

	var newCard *Card = c.DeckPile.DrawCard()
	newCard.Reveal()

	c.Players[c.playerFocus].ShieldCard = newCard
	c.DeckPile.DiscardCard(oldCard)

	return 0, "It's like getting under a blanket on a rainy day !"
}

// @desc: Player at index player puts the next card into his charge slot
func (c *Camaretto) Charge() (int, string) {
	var p *Player = c.Players[c.playerFocus]
	if p.ChargeCard == nil {
		var card *Card = c.DeckPile.DrawCard()
		p.Charge(card)
	}

	return 0, "Loading up !"
}

// @desc: Player at index player heals himself
func (c *Camaretto) Heal() (int, string) {
	var oldCard *Card = c.Players[c.playerFocus].Heal(c.cardFocus)
	c.DeckPile.DiscardCard(oldCard)

	return 0, "I feel a lil' bit tired, anyone has a vitamin ?"
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
	var centerX float64 = width/2
	var centerY float64 = (height * 6/8)/2

	for i, player := range c.Players {
		var x, y, theta float64 = c.getPlayerGeoM(i)
		player.Render(dst, centerX + x, centerY + y, theta)
	}

	c.DeckPile.Render(dst, centerX, centerY)
}