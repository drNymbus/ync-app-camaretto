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

	img *ebiten.Image
	SSprite *view.Sprite
}

// @desc: Init a new Card struct then returns it
func NewCard(name string, value int, img *ebiten.Image) *Card {
	return &Card{name, value, false, img, view.NewSprite(img, false, color.RGBA{0, 0, 0, 0})}
}

// @desc: Replace original sprite image to the back of a card
func (c *Card) Hide() {
	c.Hidden = false
	c.SSprite.Img = view.HiddenCardImage
}

// @desc: Put back the img field of the card to the SSprite.Img
func (c *Card) Reveal() {
	c.Hidden = true
	c.SSprite.Img = c.img
}

type Deck struct {
	DrawPile []*Card
	LenDrawPile int
	DiscardPile []*Card
	LenDiscardPile int
}

// @desc: Initialize the deck object with 52 cards, 2 Jokers and 1 Non-value card (a.k.a "the rule card")
func (d *Deck) Init() {
	d.DrawPile = make([]*Card, 55); d.LenDrawPile = 53
	d.DiscardPile = make([]*Card, 55); d.LenDiscardPile = 0

	for i := 0; i < 52; i++ { // Insert all cards of a Deck
		var val int = i%13
		d.DrawPile[i] = NewCard("_" + strconv.Itoa(val+1), val+1, view.CardImage[val])
		d.DrawPile[i].Hide()
	}

	d.DrawPile[52] = NewCard("Zero", 0, view.EmptyCardImage) // Add a non-value card
	d.DrawPile[53] = NewCard("Joker", 14, view.JokerImage) // Add a Joker
	d.DrawPile[54] = NewCard("Joker", 14, view.JokerImage) // Add a Joker
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
