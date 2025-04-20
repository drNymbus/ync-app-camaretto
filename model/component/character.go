package component

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"camaretto/view"
)

type Character struct {
	body *ebiten.Image

	Talking bool
	isMouthOpen bool

	image *view.PersonaImage
	SSprite *Sprite

	count int
}

func NewCharacter(name string) *Character {
	var c *Character = &Character{}

	c.image = view.LoadPersonaImage("")

	var tOp *text.DrawOptions = &text.DrawOptions{}
	tOp.ColorScale.ScaleWithColor(color.Black)
	text.Draw(c.image.NeutralClosed, name, &text.GoTextFace{Source: view.FaceSource, Size: view.FontSize}, tOp)
	text.Draw(c.image.NeutralOpen, name, &text.GoTextFace{Source: view.FaceSource, Size: view.FontSize}, tOp)

	c.SSprite = NewSprite(c.image.NeutralClosed, nil)

	return c
}

func (c *Character) Talk(state GameState) string {
	var msg string = ""
	return msg
}

// @desc: Generate sprite image depending on character's state
func (c *Character) Update() error {
	var modify bool = false
	if c.Talking {
		c.count++
		if c.count > 5 {
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
		if c.isMouthOpen {
			c.SSprite.SetImage(c.image.NeutralOpen)
		} else {
			c.SSprite.SetImage(c.image.NeutralClosed)
		}
	}

	c.SSprite.Update()
	return nil
}

func (c *Character) Draw(screen *ebiten.Image) {
	c.SSprite.Draw(screen)
}
