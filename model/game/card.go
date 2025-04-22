package game

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"camaretto/view"
)

type Card struct {
	Name string
	Value int
	Hidden bool
	
	Trigger func()

	revealedImg *ebiten.Image
	hiddenImg *ebiten.Image
	SSprite *view.Sprite
}

// @desc: Init a new Card struct then returns it
func NewCard(name string, value int, revealedImg *ebiten.Image, hiddenImg *ebiten.Image) *Card {
	var c *Card = &Card{}

	c.Name = name
	c.Value = value
	c.Hidden = false

	c.Trigger = nil

	c.revealedImg = revealedImg
	c.hiddenImg = hiddenImg
	c.SSprite = view.NewSprite(revealedImg, nil)

	return c
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

func (c *Card) Update() error {
	if c.Trigger != nil && inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		var x, y int = ebiten.CursorPosition()
		if c.SSprite.In(float64(x), float64(y)) {
			c.Trigger()
		}
	}

	c.SSprite.Update()
	return nil
}

func (c *Card) Draw(screen *ebiten.Image) {
	c.SSprite.Draw(screen)
}
