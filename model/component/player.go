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

	nameSprite *view.Sprite
	Dead bool
	deadSprite *view.Sprite

	Persona *Character

	shield *Card
	jokerShield *Card
	health [2]*Card
	jokerHealth *Card
	charge *Card
}

func NewPlayer(name string, char *Character, x, y, r float64) *Player {
	var p *Player = &Player{}

	p.x, p.y, p.r = x, y, r

	var tWidth, tHeight float64 = text.Measure(name, view.TextFace, 0.0)
	var img *ebiten.Image = ebiten.NewImage(int(tWidth), int(tHeight))
	op := &text.DrawOptions{}; op.ColorScale.ScaleWithColor(color.RGBA{0,0,0,255})
	text.Draw(img, name, &text.GoTextFace{Source: view.FaceSource, Size: view.FontSize}, op)

	p.nameSprite = view.NewSprite(img, nil)
	p.nameSprite.MoveOffset(0, float64(view.CardHeight) * 3/2, 0.5)
	p.nameSprite.RotateOffset(r, 0.5)
	p.nameSprite.Move(x, y, 1)

	p.deadSprite = view.NewSprite(view.LoadDeathImage(), nil)
	p.deadSprite.Rotate(r, 0.2)
	p.deadSprite.Move(x, y, 0.5)

	p.Dead = false

	p.Persona = char

	p.shield = nil
	p.jokerShield = nil
	p.health = [2]*Card{nil, nil}
	p.jokerHealth = nil
	p.charge = nil

	return p
}

func (p *Player) GetPosition() (float64, float64, float64) { return p.x, p.y, p.r }

// @desc: Set card at shield position modifying sprite position and all then returning the old card
func (p *Player) SetShield(c *Card) *Card {
	var old *Card = p.shield

	c.SSprite.Move(p.x, p.y, 1)
	c.SSprite.RotateOffset(p.r, 1)
	
	var xOff, yOff, r float64 = 0, -float64(view.CardWidth)/2, math.Pi/2
	c.SSprite.MoveOffset(xOff, yOff, 1)
	c.SSprite.Rotate(r, 1)

	p.shield = c
	return old
}

// @desc: Set card at joker shield position modifying sprite position and all then returning the old card
func (p *Player) SetJokerShield(c *Card) *Card {
	var old *Card = p.jokerShield

	c.SSprite.Move(p.x, p.y, 1)
	c.SSprite.RotateOffset(p.r, 1)

	var xOff, yOff, r float64 = 0, -float64(view.CardWidth)/2 - 15, math.Pi/2
	c.SSprite.MoveOffset(xOff, yOff, 1)
	c.SSprite.Rotate(r, 1)

	p.jokerShield = c
	return old
}


// @desc: Set card at health[i] position modifying sprite position and all then returning the old card
func (p *Player) SetHealth(c *Card, i int) *Card {
	var old *Card = p.health[i]

	c.SSprite.Move(p.x, p.y, 1)
	c.SSprite.RotateOffset(p.r, 1)

	var xOff float64 = float64((i-1) * view.CardWidth) + float64(view.CardWidth)/2
	var yOff float64 = float64(view.CardHeight)/2
	var r float64 = 0
	c.SSprite.MoveOffset(xOff, yOff, 1)
	c.SSprite.Rotate(r, 1)

	p.health[i] = c
	return old
}

// @desc: Set card at joker health position modifying sprite position and all then returning the old card
func (p *Player) SetJokerHealth(c *Card) *Card {
	var old *Card = p.jokerHealth

	c.SSprite.Move(p.x, p.y, 1)
	c.SSprite.RotateOffset(p.r, 1)

	var xOff float64 = - float64(view.CardWidth) - float64(view.CardWidth)/2
	var yOff float64 = float64(view.CardHeight)/2
	var r float64 = 0
	c.SSprite.MoveOffset(xOff, yOff, 1)
	c.SSprite.Rotate(r, 1)

	p.jokerHealth = c
	return old
}

// @desc: Return true in case charge is empty, false otherwise
func (p *Player) IsChargeEmpty() bool { return p.charge == nil }

// @desc: Set card at charge position modifying sprite position and all then returning the old card
func (p *Player) SetCharge(c *Card) *Card {
	var old *Card = p.charge

	c.SSprite.Move(p.x, p.y, 1)
	c.SSprite.RotateOffset(p.r, 1)

	var xOff float64 = float64(view.CardWidth) + float64(view.CardWidth)/2
	var yOff float64 = float64(view.CardHeight)/2
	var r float64 = 0
	c.SSprite.MoveOffset(xOff, yOff, 1)
	c.SSprite.Rotate(r, 1)

	p.charge = c
	return old
}

func (p *Player) ResetTrigger() {
	if p.shield != nil { p.shield.Trigger = nil }
	if p.jokerShield != nil { p.jokerShield.Trigger = nil }

	if p.health[0] != nil { p.health[0].Trigger = nil }
	if p.health[1] != nil { p.health[1].Trigger = nil }
	if p.jokerHealth != nil { p.jokerHealth.Trigger = nil }

	if p.charge != nil { p.charge.Trigger = nil }
}

// @desc: Set callback for any card clicked
func (p *Player) OnPlayer(trigger func()) {
	if p.shield != nil { p.shield.Trigger = trigger }
	if p.jokerShield != nil { p.jokerShield.Trigger = trigger }

	if p.health[0] != nil { p.health[0].Trigger = trigger }
	if p.health[1] != nil { p.health[1].Trigger = trigger }
	if p.jokerHealth != nil { p.jokerHealth.Trigger = trigger }

	if p.charge != nil { p.charge.Trigger = trigger }
}

// @desc: Set callback for health cards and nil for all other cards
func (p *Player) OnHealth(trigger func(int)) {
	p.ResetTrigger()
	for i, card := range p.health {
		if card != nil {
			card.Trigger = func() { trigger(i) }
		}
	}
}

func (p *Player) Update() error {
	if p.shield != nil { p.shield.Update() }
	if p.jokerShield != nil { p.jokerShield.Update() }

	if p.health[0] != nil { p.health[0].Update() }
	if p.health[1] != nil { p.health[1].Update() }
	if p.jokerHealth != nil { p.jokerHealth.Update() }

	if p.charge != nil { p.charge.Update() }

	return nil
}

func (p *Player) Draw(screen *ebiten.Image) {
	if p.Dead {
		p.deadSprite.Draw(screen)
	} else {
		if p.shield != nil { p.shield.Draw(screen) }
		if p.jokerShield != nil { p.jokerShield.Draw(screen) }

		if p.health[0] != nil { p.health[0].Draw(screen) }
		if p.health[1] != nil { p.health[1].Draw(screen) }
		if p.jokerHealth != nil { p.jokerHealth.Draw(screen) }

		if p.charge != nil { p.charge.Draw(screen) }
	}

	p.nameSprite.Draw(screen)
}
