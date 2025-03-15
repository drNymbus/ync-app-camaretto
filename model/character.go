package model

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"camaretto/view"
)

type Character struct {
	bodyWidth, bodyHeight float64
	body *ebiten.Image

	Talking bool
	isMouthOpen bool

	mouthWidth, mouthHeight float64
	openMouth, closedMouth *ebiten.Image

	image *ebiten.Image
	SSprite *view.Sprite

	count int
}

func NewCharacter(name string) *Character {
	var c *Character = &Character{}

	var pi *view.PersonaImage = view.LoadPersonaImage("")
	c.bodyWidth, c.bodyHeight = float64(view.PersonaBodyWidth), float64(view.PersonaBodyHeight)
	c.body = pi.Body

	var tOp *text.DrawOptions = &text.DrawOptions{}
	tOp.ColorScale.ScaleWithColor(color.Black)
	text.Draw(c.body, name, &text.GoTextFace{Source: view.FaceSource, Size: view.FontSize}, tOp)

	c.Talking = false
	c.isMouthOpen = false

	c.mouthWidth, c.mouthHeight = float64(view.PersonaMouthWidth), float64(view.PersonaMouthHeight)
	c.openMouth = pi.OpenMouth
	c.closedMouth = pi.ClosedMouth

	c.image = ebiten.NewImage(view.PersonaBodyWidth, view.PersonaBodyHeight)
	c.SSprite = view.NewSprite(c.image, false, color.RGBA{0,0,0,0}, nil)

	return c
}

func (c *Character) Talk(state GameState) string {
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

	return msg
}

// @desc: Generate sprite image depending on character's state
func (c *Character) Render() {
	var modify bool = false
	if c.Talking {
		c.count++
		if c.count > 7 {
			c.isMouthOpen = !c.isMouthOpen
			c.count = 0
			modify = true
		}
	} else {
		c.isMouthOpen = false
		c.count = 0
		modify = true
	}

	if modify {
		c.image.Clear()
		c.image.DrawImage(c.body, nil)

		var op *ebiten.DrawImageOptions = &ebiten.DrawImageOptions{}
		op.GeoM.Translate(-float64(c.mouthWidth)/2, -float64(c.mouthHeight)/2)
		op.GeoM.Translate(float64(c.bodyWidth)/2, float64(c.bodyHeight)*2/5)

		if c.isMouthOpen {
			c.image.DrawImage(c.openMouth, op)
		} else {
			c.image.DrawImage(c.closedMouth, op)
		}

		c.SSprite.SetImage(c.image)
	}
}