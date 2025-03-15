package model

import (
	"log"
	"time"
	"math/rand"

	"strconv"

	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"camaretto/view"
)

type Card struct {
	Name string
	Value int
	Hidden bool

	revealedImg *ebiten.Image
	hiddenImg *ebiten.Image
	SSprite *view.Sprite
}

// @desc: Init a new Card struct then returns it
func NewCard(name string, value int, revealedImg *ebiten.Image, hiddenImg *ebiten.Image) *Card {
	var s *view.Sprite = view.NewSprite(revealedImg, false, color.RGBA{127, 0, 100, 100}, nil)
	// var a *view.Animation = view.NewAnimation(s)
	return &Card{name, value, false, revealedImg, hiddenImg, s}
}

// @desc: Replace original sprite image to the back of a card
func (c *Card) Hide() {
	c.Hidden = true
	c.SSprite.SetImage(c.hiddenImg)
}

// @desc: Put back the img field of the card to the SSprite.Img
func (c *Card) Reveal() {
	c.Hidden = false
	c.SSprite.SetImage(c.revealedImg)
}

type Deck struct {
	DrawPile []*Card
	LenDrawPile int
	DrawPileX, DrawPileY float64
	DiscardPile []*Card
	LenDiscardPile int
	DiscardPileX, DiscardPileY float64
}

// @desc: Initialize the deck object with 52 cards, 2 Jokers and 1 Non-value card (a.k.a "the rule card")
func (d *Deck) Init() {
	d.DrawPile = make([]*Card, 55); d.LenDrawPile = 53
	d.DiscardPile = make([]*Card, 55); d.LenDiscardPile = 0

	var ci *view.CardImage = view.LoadCardImage()

	for i := 0; i < 52; i++ { // Insert all cards of a Deck
		var val int = i%13
		d.DrawPile[i] = NewCard("_" + strconv.Itoa(val+1), val+1, ci.Card[val], ci.Hidden)
		d.DrawPile[i].Hide()
	}

	// Add a non-value card
	d.DrawPile[52] = NewCard("Zero", 0, ci.Empty, ci.Hidden)
	d.DrawPile[52].Hide()
	// Add a Joker
	d.DrawPile[53] = NewCard("Joker", 14, ci.Joker, ci.Hidden)
	d.DrawPile[53].Hide()
	// Add a Joker
	d.DrawPile[54] = NewCard("Joker", 14, ci.Joker, ci.Hidden)
	d.DrawPile[54].Hide()
}

// @desc: Puts all cards in the discard pile on top of the draw pile
func (d *Deck) ResetDrawPile() {
	for ;d.LenDiscardPile > 0; {
		var c *Card = d.DiscardPile[d.LenDiscardPile-1]
		c.Hide()
		d.DrawPile[d.LenDrawPile] = c

		d.LenDrawPile++; d.LenDiscardPile--
	}
}

// @desc: Draws one card on top of the draw pile then returns the Card object
//        If the draw pile is empty, calls ResetDrawPile before drawing one
func (d *Deck) DrawCard() *Card {
	if d.LenDrawPile < 1 {
		d.ResetDrawPile()
	}

	var c *Card = d.DrawPile[d.LenDrawPile-1]
	d.LenDrawPile--
	return c
}

// @desc: Discard the card passed as a parameter on top of the discard pile
func (d *Deck) DiscardCard(c *Card) {
	if d.LenDiscardPile > 54 {
		log.Fatal("[Deck.DiscardCard] How come you're discarding this one lil fella :/")
	}
	d.DiscardPile[d.LenDiscardPile] = c
	d.LenDiscardPile++
}

// @desc: Randomize order of cards in the draw pile
func (d *Deck) ShuffleDrawPile() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(d.LenDrawPile, func(i, j int) { d.DrawPile[i], d.DrawPile[j] = d.DrawPile[j], d.DrawPile[i] })
}

// @desc: Randomize order of cards in the discard pile
func (d *Deck) ShuffleDiscardPile() {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(d.LenDiscardPile, func(i, j int) { d.DiscardPile[i], d.DiscardPile[j] = d.DiscardPile[j], d.DiscardPile[i] })
}

func (d *Deck) FindInDrawPile(val int) *Card {
	for i := 0; i < d.LenDrawPile; i++ {
		var c *Card = d.DrawPile[i]
		if c.Value == val {
			d.DrawPile = append(d.DrawPile[:i], d.DrawPile[i+1:]...)
			d.LenDrawPile--
			return c
		}
	}
	return nil
}

func (d *Deck) FindInDiscardPile(val int) *Card {
	for i := d.LenDiscardPile-1; i >= 0; i-- {
		var c *Card = d.DiscardPile[i]
		if c.Value == val {
			d.DiscardPile = append(d.DiscardPile[:i], d.DiscardPile[i+1:]...)
			d.LenDiscardPile--
			return c
		}
	}
	return nil
}

func (d *Deck) Render(dst *ebiten.Image, x, y float64) {
	var speed, rSpeed float64 = 0.5, 0.2
	var cOff float64 = 0.2

	d.DrawPileX, d.DrawPileY = x - float64(view.CardWidth)/2, y

	var posX, posY float64
	var cx, cy, cr float64

	for i, card := range d.DrawPile[:d.LenDrawPile] {
		posX, posY = d.DrawPileX, d.DrawPileY - float64(i)*cOff
		cx, cy, cr = card.SSprite.GetCenter()
		if cx != posX || cy != posY || cr != 0 {
			card.SSprite.Move(posX, posY, speed)
			card.SSprite.Rotate(0, rSpeed)
			card.SSprite.MoveOffset(0, 0, speed)
			card.SSprite.RotateOffset(0, rSpeed)
		}
		card.SSprite.Display(dst)
	}

	d.DiscardPileX, d.DiscardPileY = x + float64(view.CardWidth)/2, y
	for i, card := range d.DiscardPile[:d.LenDiscardPile] {
		posX, posY = d.DiscardPileX, d.DiscardPileY - float64(i)*cOff
		cx, cy, cr = card.SSprite.GetCenter()
		if cx != posX || cy != posY || cr != 0 {
			card.SSprite.Move(posX, posY, speed)
			card.SSprite.Rotate(0, rSpeed)
			card.SSprite.MoveOffset(0, 0, speed)
			card.SSprite.RotateOffset(0, rSpeed)
		}
		card.SSprite.Display(dst)
	}
}