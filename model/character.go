package model

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	// "github.com/hajimehoshi/ebiten/v2/text/v2"

	"camaretto/view"
)

type Character struct {
	body *ebiten.Image
	isMouthOpen bool
	count int
	openMouth, closedMouth *ebiten.Image
	img *ebiten.Image

	speech *TextBox

	SSprite *view.Sprite
}

func NewCharacter(tb *TextBox) *Character {
	var c *Character = &Character{}

	c.body = view.GetImage("assets/characters/char.png")
	c.isMouthOpen = false
	c.openMouth = view.GetImage("assets/characters/mouth_open.png")
	c.closedMouth = view.GetImage("assets/characters/mouth_closed.png")

	c.speech = tb

	c.SSprite = view.NewSprite(c.body, false, color.RGBA{0,0,0,0}, nil)
	c.SSprite.Scale(0.25, 0.25)

	return c
}

func (c *Character) GetImg() *ebiten.Image { return c.img }

// @desc: Set new image to sprite depending on character's state
func (c *Character) RenderBody() {
	// var scale float64 = 0.25
	// var bodyW, bodyH float64 = float64(view.CharacterWidth)*scale, float64(view.CharacterHeight)*scale

	var img *ebiten.Image = ebiten.NewImage(view.CharacterWidth, view.CharacterHeight)
	img.DrawImage(c.body, nil)
	
	// var mouthW, mouthH = float64(view.MouthWidth)*scale, float64(view.MouthHeight)*scale
	var op *ebiten.DrawImageOptions = &ebiten.DrawImageOptions{}
	op.GeoM.Translate(-float64(view.MouthWidth)/2, -float64(view.MouthHeight)/2)
	op.GeoM.Translate(float64(view.CharacterWidth)/2, float64(view.CharacterHeight)*2/5)

	if c.isMouthOpen {
		img.DrawImage(c.openMouth, op)
	} else {
		img.DrawImage(c.closedMouth, op)
	}

	c.SSprite.SetImage(img)
}

// @desc: Render body and textbox
func (c *Character) Render(dst *ebiten.Image, x, y float64) {
	c.speech.Render()

	if c.speech.Finished() {
		c.isMouthOpen = false
		c.count = 0
	} else {
		c.count++
		if c.count > 7 {
			c.isMouthOpen = !c.isMouthOpen
			c.count = 0
		}
	}
	c.RenderBody()

	c.speech.SSprite.SetCenter(x, y, 0)
	var bodyX float64 = (x - c.speech.SSprite.Width/2) + c.SSprite.Width/2
	var bodyY float64 = (y + c.speech.SSprite.Height/2) - c.SSprite.Height/2
	c.SSprite.SetCenter(bodyX, bodyY, 0)

	c.SSprite.Display(dst)
	c.speech.SSprite.Display(dst)
}