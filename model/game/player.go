package game

import (
	"math"

	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"camaretto/view"
)

type PlayerInfo struct {
	Index int
	Name string
}

type Player struct {
	x, y, r float64

	nameSprite *view.Sprite
	Dead bool
	deadSprite *view.Sprite

	Persona *Character

	Shield *Card
	JokerShield *Card
	Health [2]*Card
	JokerHealth *Card
	Charge *Card

	Trigger func()
}

func NewPlayer(name string, char *Character, x, y, r float64) *Player {
	var p *Player = &Player{}

	p.x, p.y, p.r = x, y, r

	var tWidth, tHeight float64 = text.Measure(name, view.TextFace, 0.0)
	var img *ebiten.Image = ebiten.NewImage(int(tWidth), int(tHeight))
	op := &text.DrawOptions{}; op.ColorScale.ScaleWithColor(color.RGBA{0,0,0,255})
	text.Draw(img, name, &text.GoTextFace{Source: view.FaceSource, Size: view.FontSize}, op)

	p.nameSprite = view.NewSprite(img, nil)
	p.nameSprite.MoveOffset(0, float64(view.CardHeight) * 2, 0.5)
	p.nameSprite.RotateOffset(r, 0.5)
	p.nameSprite.Move(x, y, 1)
	p.nameSprite.Rotate(-r, 0.5)

	p.deadSprite = view.NewSprite(view.LoadDeathImage(), nil)
	p.deadSprite.Rotate(r, 0.2)
	p.deadSprite.Move(x, y, 0.5)

	p.Dead = false

	p.Persona = char

	p.Shield = nil
	p.JokerShield = nil
	p.Health = [2]*Card{nil, nil}
	p.JokerHealth = nil
	p.Charge = nil

	p.Trigger = nil

	return p
}

func (p *Player) GetPosition() (float64, float64, float64) { return p.x, p.y, p.r }

// @desc: Set card at shield position modifying sprite position and all then returning the old card
func (p *Player) SetShield(c *Card) *Card {
	var old *Card = p.Shield

	if c != nil {
		c.SSprite.Move(p.x, p.y, 1)
		c.SSprite.RotateOffset(p.r, 1)
		
		var xOff, yOff, r float64 = 0, -float64(view.CardWidth)/2, math.Pi/2
		c.SSprite.MoveOffset(xOff, yOff, 1)
		c.SSprite.Rotate(r, 1)
	}

	p.Shield = c
	return old
}

// @desc: Set card at joker shield position modifying sprite position and all then returning the old card
func (p *Player) SetJokerShield(c *Card) *Card {
	var old *Card = p.JokerShield

	if c != nil {
		c.SSprite.Move(p.x, p.y, 1)
		c.SSprite.RotateOffset(p.r, 1)

		var xOff, yOff, r float64 = 0, -float64(view.CardWidth)/2 - 15, math.Pi/2
		c.SSprite.MoveOffset(xOff, yOff, 1)
		c.SSprite.Rotate(r, 1)
	}

	p.JokerShield = c
	return old
}


// @desc: Set card at Health[i] position modifying sprite position and all then returning the old card
func (p *Player) SetHealth(c *Card, i int) *Card {
	var old *Card = p.Health[i]

	if c != nil {
		c.SSprite.Move(p.x, p.y, 1)
		c.SSprite.RotateOffset(p.r, 1)

		var xOff float64 = float64((i-1) * view.CardWidth) + float64(view.CardWidth)/2
		var yOff float64 = float64(view.CardHeight)/2
		var r float64 = 0
		c.SSprite.MoveOffset(xOff, yOff, 1)
		c.SSprite.Rotate(r, 1)
	}

	p.Health[i] = c
	return old
}

// @desc: Set card at joker Health position modifying sprite position and all then returning the old card
func (p *Player) SetJokerHealth(c *Card) *Card {
	var old *Card = p.JokerHealth

	if c != nil {
		c.SSprite.Move(p.x, p.y, 1)
		c.SSprite.RotateOffset(p.r, 1)

		var xOff float64 = - float64(view.CardWidth) - float64(view.CardWidth)/2
		var yOff float64 = float64(view.CardHeight)/2
		var r float64 = 0
		c.SSprite.MoveOffset(xOff, yOff, 1)
		c.SSprite.Rotate(r, 1)
	}

	p.JokerHealth = c
	return old
}

// @desc: Return true in case Charge is empty, false otherwise
func (p *Player) IsChargeEmpty() bool { return p.Charge == nil }

// @desc: Set card at Charge position modifying sprite position and all then returning the old card
func (p *Player) SetCharge(c *Card) *Card {
	var old *Card = p.Charge

	if c != nil {
		c.SSprite.Move(p.x, p.y, 1)
		c.SSprite.RotateOffset(p.r, 1)

		var xOff float64 = float64(view.CardWidth) + float64(view.CardWidth)/2
		var yOff float64 = float64(view.CardHeight)/2
		var r float64 = 0
		c.SSprite.MoveOffset(xOff, yOff, 1)
		c.SSprite.Rotate(r, 1)
	}

	p.Charge = c
	return old
}

// @desc:
func (p *Player) Update() error {
	p.nameSprite.Update()
	p.deadSprite.Update()

	var ix, iy int = ebiten.CursorPosition()
	var x, y float64 = float64(ix), float64(iy)

	var cursorIn bool = false

	if p.Shield != nil {
		cursorIn = cursorIn || p.Shield.SSprite.In(x,y)
		p.Shield.Update()
	}

	if p.JokerShield != nil {
		cursorIn = cursorIn || p.JokerShield.In(x, y)
		p.JokerShield.Update()
	}

	if p.Health[0] != nil {
		cursorIn = cursorIn || p.Health[0].In(x, y)
		p.Health[0].Update()
	}

	if p.Health[1] != nil {
		cursorIn = cursorIn || p.Health[1].In(x, y)
		p.Health[1].Update()
	}

	if p.JokerHealth != nil {
		cursorIn = cursorIn || p.JokerHealth.In(x, y)
		p.JokerHealth.Update()
	}

	if p.Charge != nil {
		cursorIn = cursorIn || p.Charge.In(x, y)
		p.Charge.Update()
	}

	if p.Trigger != nil {
		if cursorIn && inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
			p.Trigger()
		}
	}

	return nil
}

func (p *Player) Draw(screen *ebiten.Image) {
	if p.Dead {
		p.deadSprite.Draw(screen)
	} else {
		if p.Shield != nil { p.Shield.Draw(screen) }
		if p.JokerShield != nil { p.JokerShield.Draw(screen) }

		if p.Health[0] != nil { p.Health[0].Draw(screen) }
		if p.Health[1] != nil { p.Health[1].Draw(screen) }
		if p.JokerHealth != nil { p.JokerHealth.Draw(screen) }

		if p.Charge != nil { p.Charge.Draw(screen) }
	}

	p.nameSprite.Draw(screen)
}
