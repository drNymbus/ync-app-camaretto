package component

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

	nameSprite *Sprite
	Dead bool
	deadSprite *Sprite

	Persona *Character

	Shield *Card
	JokerShield *Card
	Health [2]*Card
	JokerHealth *Card
	Charge *Card
}

func NewPlayer(name string, char *Character, x, y, r float64) *Player {
	var p *Player = &Player{}

	p.x, p.y, p.r = x, y, r

	var tWidth, tHeight float64 = text.Measure(name, view.TextFace, 0.0)
	var img *ebiten.Image = ebiten.NewImage(int(tWidth), int(tHeight))
	op := &text.DrawOptions{}; op.ColorScale.ScaleWithColor(color.RGBA{0,0,0,255})
	text.Draw(img, name, &text.GoTextFace{Source: view.FaceSource, Size: view.FontSize}, op)

	p.nameSprite = NewSprite(img, nil)
	p.nameSprite.MoveOffset(0, float64(view.CardHeight) * 2, 0.5)
	p.nameSprite.RotateOffset(r, 0.5)
	p.nameSprite.Move(x, y, 1)
	p.nameSprite.Rotate(-r, 0.5)

	p.deadSprite = NewSprite(view.LoadDeathImage(), nil)
	p.deadSprite.Rotate(r, 0.2)
	p.deadSprite.Move(x, y, 0.5)

	p.Dead = false

	p.Persona = char

	p.Shield = nil
	p.JokerShield = nil
	p.Health = [2]*Card{nil, nil}
	p.JokerHealth = nil
	p.Charge = nil

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

func (p *Player) ResetTrigger() {
	if p.Shield != nil { p.Shield.Trigger = nil }
	if p.JokerShield != nil { p.JokerShield.Trigger = nil }

	if p.Health[0] != nil { p.Health[0].Trigger = nil }
	if p.Health[1] != nil { p.Health[1].Trigger = nil }
	if p.JokerHealth != nil { p.JokerHealth.Trigger = nil }

	if p.Charge != nil { p.Charge.Trigger = nil }
}

// @desc: Set callback for any card clicked
func (p *Player) OnPlayer(trigger func()) {
	if p.Shield != nil { p.Shield.Trigger = trigger }
	if p.JokerShield != nil { p.JokerShield.Trigger = trigger }

	if p.Health[0] != nil { p.Health[0].Trigger = trigger }
	if p.Health[1] != nil { p.Health[1].Trigger = trigger }
	if p.JokerHealth != nil { p.JokerHealth.Trigger = trigger }

	if p.Charge != nil { p.Charge.Trigger = trigger }
}

// @desc: Set callback for Health cards and nil for all other cards
func (p *Player) OnHealth(trigger func(int)) {
	p.ResetTrigger()
	for i, card := range p.Health {
		if card != nil {
			card.Trigger = func() { trigger(i) }
		}
	}
}

func (p *Player) HoverPlayer(x, y float64) bool {
	var flag bool = false

	if p.Shield != nil { flag = flag || p.Shield.SSprite.In(x, y) }
	if p.JokerShield != nil { flag = flag || p.JokerShield.SSprite.In(x, y) }

	if p.Health[0] != nil { flag = flag || p.Health[0].SSprite.In(x, y) }
	if p.Health[1] != nil { flag = flag || p.Health[1].SSprite.In(x, y) }
	if p.JokerHealth != nil { flag = flag || p.JokerHealth.SSprite.In(x, y) }

	if p.Charge != nil { flag = flag || p.Charge.SSprite.In(x, y) }

	return flag
}

func (p *Player) HoverHealth(x, y float64) int {
	if p.Health[0] != nil && p.Health[0].SSprite.In(x, y) { return 0 }
	if p.Health[1] != nil && p.Health[1].SSprite.In(x, y) { return 1 }
	return -1
}

func (p *Player) Update() error {
	p.nameSprite.Update()
	p.deadSprite.Update()

	if p.Shield != nil { p.Shield.Update() }
	if p.JokerShield != nil { p.JokerShield.Update() }

	if p.Health[0] != nil { p.Health[0].Update() }
	if p.Health[1] != nil { p.Health[1].Update() }
	if p.JokerHealth != nil { p.JokerHealth.Update() }

	if p.Charge != nil { p.Charge.Update() }

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
