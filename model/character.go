package model

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	// "github.com/hajimehoshi/ebiten/v2/text/v2"

	"camaretto/view"
)

type Character struct {
	bodyWidth, bodyHeight float64
	body *ebiten.Image
	isMouthOpen bool
	count int

	mouthWidth, mouthHeight float64
	openMouth, closedMouth *ebiten.Image

	Speech *TextBox

	SSprite *view.Sprite
}

func NewCharacter(tb *TextBox, name string) *Character {
	var c *Character = &Character{}

	var pi *view.PersonaImage = view.LoadPersonaImage("")
	// var scale float64 = 0.25
	// var op *ebiten.DrawImageOptions = &ebiten.DrawImageOptions{}
	// op.GeoM.Scale(scale, scale)

	// var charBody *ebiten.Image = view.GetImage("assets/characters/char.png")
	// c.bodyWidth, c.bodyHeight = float64(view.CharacterWidth)*scale, float64(view.CharacterHeight)*scale
	c.bodyWidth, c.bodyHeight = float64(view.PersonaBodyWidth), float64(view.PersonaBodyHeight)
	// c.body = ebiten.NewImage(int(c.bodyWidth), int(c.bodyHeight))
	// c.body.DrawImage(charBody, op)
	c.body = pi.Body

	// var tOp *text.DrawOptions = &text.DrawOptions{}
	// tOp.ColorScale.ScaleWithColor(color.Black)
	// text.Draw(c.body, name, &text.GoTextFace{Source: view.FaceSource, Size: view.FontSize}, tOp)

	c.isMouthOpen = false

	// c.mouthWidth, c.mouthHeight = float64(view.MouthWidth)*scale, float64(view.MouthHeight)*scale
	c.mouthWidth, c.mouthHeight = float64(view.PersonaMouthWidth), float64(view.PersonaMouthHeight)

	// var charOpenMouth *ebiten.Image = view.GetImage("assets/characters/mouth_open.png")
	// c.openMouth = ebiten.NewImage(int(c.mouthWidth), int(c.mouthHeight))
	// c.openMouth.DrawImage(charOpenMouth, op)
	c.openMouth = pi.OpenMouth

	// var charClosedMouth *ebiten.Image = view.GetImage("assets/characters/mouth_closed.png")
	// c.closedMouth = ebiten.NewImage(int(c.mouthWidth), int(c.mouthHeight))
	// c.closedMouth.DrawImage(charClosedMouth, op)
	c.closedMouth = pi.ClosedMouth

	c.Speech = tb
	// c.Speech.SSprite.SetCenter(-c.bodyWidth, -c.bodyHeight, 0)


	c.SSprite = view.NewSprite(c.body, false, color.RGBA{0,0,0,0}, nil)
	// c.SSprite.SetCenter(-c.bodyWidth, -c.bodyHeight, 0)

	return c
}

func (c *Character) Talk(state GameState) {
	var msg string = ""
	if state == SET {
		msg = "Choisis une action, ego player que tu es ! Tu crois j'tai pas vu ? va jouer à la dinette plutot"
	} else if state == ATTACK {
		msg = "Je vais attaquer de toute ma puissance !"
	} else if state == SHIELD {
		msg = "Changement ! Zinedine rentre sur le terrain afin de donner du sang neuf, on voyait bien que Zidane commençait à fatiguer."
	} else if state == CHARGE {
		msg = "Meditation mode"
	} else if state == HEAL {
		msg = "Regenaration de mes pouvoirs"
	}

	c.Speech.SetMessage(msg)
}

// @desc: Set new image to sprite depending on character's state
func (c *Character) RenderBody() {
	// var scale float64 = 0.25
	// var bodyW, bodyH float64 = float64(view.CharacterWidth)*scale, float64(view.CharacterHeight)*scale

	var img *ebiten.Image = ebiten.NewImage(int(c.bodyWidth), int(c.bodyHeight))
	img.DrawImage(c.body, nil)
	
	var op *ebiten.DrawImageOptions = &ebiten.DrawImageOptions{}
	// var mouthW, mouthH = float64(view.MouthWidth)*scale, float64(view.MouthHeight)*scale
	op.GeoM.Translate(-float64(c.mouthWidth)/2, -float64(c.mouthHeight)/2)
	op.GeoM.Translate(float64(c.bodyWidth)/2, float64(c.bodyHeight)*2/5)

	if c.isMouthOpen {
		img.DrawImage(c.openMouth, op)
	} else {
		img.DrawImage(c.closedMouth, op)
	}

	c.SSprite.SetImage(img)
}

// @desc: Render body and textbox
func (c *Character) Render(x, y float64) {
	c.Speech.Render()

	if c.Speech.Finished() {
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

	c.Speech.SSprite.SetCenter(x, y, 0)
	var bodyX float64 = (x - c.Speech.SSprite.Width/2) + c.SSprite.Width/2
	var bodyY float64 = (y + c.Speech.SSprite.Height/2) - c.SSprite.Height/2
	c.SSprite.SetCenter(bodyX, bodyY, 0)
}