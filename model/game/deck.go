package game

import (
	"log"
	"math/rand"

	"strconv"

	"github.com/hajimehoshi/ebiten/v2"

	"camaretto/view"
)

type Deck struct {
	x, y float64

	Seed int64

	DrawPile []*Card
	LenDrawPile int

	DiscardPile []*Card
	LenDiscardPile int
}

// @desc: Initialize the deck object with 52 cards, 2 Jokers and 1 Non-value card (a.k.a "the rule card")
func (d *Deck) Init(s int64, x, y float64) {
	d.x, d.y = x, y

	d.DrawPile = make([]*Card, 55); d.LenDrawPile = 55
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

	d.Seed = s
	d.shuffleDrawPile()

	for i, card := range d.DrawPile {
		card.SSprite.Move(d.x - float64(view.CardWidth)/2, d.y - float64(i)*0.2, 1)
	}
}

// @desc: Randomize order of cards in the draw pile
func (d *Deck) shuffleDrawPile() {
	rand.Seed(d.Seed)
	rand.Shuffle(d.LenDrawPile, func(i, j int) {
		d.DrawPile[i], d.DrawPile[j] = d.DrawPile[j], d.DrawPile[i]
	})
}

// @desc: Puts all cards in the discard pile on top of the draw pile
func (d *Deck) resetDrawPile() {
	for ;d.LenDiscardPile > 0; {
		var card *Card = d.DiscardPile[d.LenDiscardPile-1]
		card.Hide()

		d.DrawPile[d.LenDrawPile] = card
		d.LenDrawPile++; d.LenDiscardPile--

		card.SSprite.Move(d.x - float64(view.CardWidth)/2, d.y - float64(d.LenDrawPile)*0.2, 1)
		card.SSprite.Rotate(0, 1)
		card.SSprite.MoveOffset(0, 0, 1)
		card.SSprite.RotateOffset(0, 1)
	}
}

// @desc: Draws one card on top of the draw pile then returns the Card object
//        If the draw pile is empty, calls resetDrawPile before drawing one
func (d *Deck) DrawCard() *Card {
	if d.LenDrawPile < 1 {
		d.resetDrawPile()
	}

	var card *Card = d.DrawPile[d.LenDrawPile-1]
	d.LenDrawPile--
	return card
}

// @desc: Discard the card passed as a parameter on top of the discard pile
func (d *Deck) DiscardCard(card *Card) {
	if d.LenDiscardPile > 54 {
		log.Fatal("[Deck.DiscardCard] How come you're discarding this one lil fella :/")
	}

	if card != nil {
		d.DiscardPile[d.LenDiscardPile] = card
		d.LenDiscardPile++

		card.SSprite.Move(d.x + float64(view.CardWidth)/2, d.y - float64(d.LenDiscardPile)*0.2, 1)
		card.SSprite.Rotate(0, 1)
		card.SSprite.MoveOffset(0, 0, 1)
		card.SSprite.RotateOffset(0, 1)
	}
}

// @desc: Find the first card that has the value asked for in the draw pile then return its reference
func (d *Deck) FindInDrawPile(val int) *Card {
	for i := 0; i < d.LenDrawPile; i++ {
		var card *Card = d.DrawPile[i]
		if card.Value == val {
			d.DrawPile = append(d.DrawPile[:i], d.DrawPile[i+1:]...)
			d.LenDrawPile--
			return card
		}
	}
	return nil
}

// @desc: Find the first card that has the value asked for in the discard pile then return its reference
func (d *Deck) FindInDiscardPile(val int) *Card {
	for i := d.LenDiscardPile-1; i >= 0; i-- {
		var card *Card = d.DiscardPile[i]
		if card.Value == val {
			d.DiscardPile = append(d.DiscardPile[:i], d.DiscardPile[i+1:]...)
			d.LenDiscardPile--
			return card
		}
	}
	return nil
}

func (d *Deck) Update() error {
	for _, card := range d.DrawPile[:d.LenDrawPile] {
		card.Update()
	}

	for _, card := range d.DiscardPile[:d.LenDiscardPile] {
		card.Update()
	}

	return nil
}

func (d *Deck) Draw(screen *ebiten.Image) {
	for _, card := range d.DrawPile[:d.LenDrawPile] {
		card.Draw(screen)
	}

	for _, card := range d.DiscardPile[:d.LenDiscardPile] {
		card.Draw(screen)
	}
}
