package model

import (
	"log"

	"math"

	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"camaretto/view"
)

type Player struct {
	Name string
	NameSprite *view.Sprite
	Dead bool
	DeadSprite *view.Sprite

	Persona *Character

	HealthCard [2]*Card
	JokerHealth *Card
	ShieldCard *Card
	JokerShield *Card
	ChargeCard *Card
}

func NewPlayer(name string, char *Character) *Player {
	var p *Player = &Player{}
	p.Name = name

	var tWidth, tHeight float64 = text.Measure(name, view.TextFace, 0.0)
	var img *ebiten.Image = ebiten.NewImage(int(tWidth), int(tHeight))

	op := &text.DrawOptions{}; op.ColorScale.ScaleWithColor(color.RGBA{0,0,0,255})
	text.Draw(img, name, &text.GoTextFace{Source: view.FaceSource, Size: view.FontSize}, op)
	var nameSprite *view.Sprite = view.NewSprite(img, true, color.RGBA{75,75,75,127}, nil)
	p.NameSprite = nameSprite

	// var deathSprite *view.Sprite = view.NewSprite(view.GraveImage, false, color.RGBA{0,0,0,0}, nil)
	// p.DeadSprite = deathSprite
	p.DeadSprite = view.NewSprite(view.LoadDeathImage(), false, color.RGBA{0,0,0,0}, nil)
	p.Dead = false

	// p.Persona = nil
	p.Persona = char

	// return &Player{name, nameSprite, false, deathSprite, char, [2]*Card{nil, nil}, nil, nil, nil, nil}
	return p
}

// @desc: Swap charge slot's card with the health card at index at then returns the old health card
func (p *Player) Heal(at int) *Card {
	if (p.ChargeCard == nil) { log.Fatal("[Player.Heal] No card in Player.ChargeCard") }

	var c *Card = p.Uncharge()
	p.HealthCard[at], c = c, p.HealthCard[at]
	return c
}

// @desc: Insert card c into the charge slot
func (p *Player) Charge(c *Card) { 
	if (p.ChargeCard != nil) { log.Fatal("[Player.Charge] Already a card in Player.ChargeCard") }
	p.ChargeCard = c
}

func (p *Player) Uncharge() *Card {
	if (p.ChargeCard == nil) { log.Fatal("[Player.Uncharge] No card in Player.ChargeCard") }
	var c *Card = p.ChargeCard
	c.Reveal()
	p.ChargeCard = nil
	return c
}

func (p *Player) getShieldOffset() (float64, float64, float64) {
	return 0, -float64(view.CardWidth)/2, math.Pi/2
}

func (p *Player) getJokerShieldOffset() (float64, float64, float64) {
	return 0, -float64(view.CardWidth)/2 - 15, math.Pi/2
}

func (p *Player) getJokerHealthOffset() (float64, float64, float64) {
	return - float64(view.CardWidth) - float64(view.CardWidth)/2, float64(view.CardHeight)/2, 0
}

func (p *Player) getHealthOffset(i int) (float64, float64, float64) {
	var x float64 = float64((i-1) * view.CardWidth) + float64(view.CardWidth)/2
	return x, float64(view.CardHeight)/2, 0
}

func (p *Player) getChargeOffset() (float64, float64, float64) {
	return float64(view.CardWidth) + float64(view.CardWidth)/2, float64(view.CardHeight)/2, 0
}

func (p *Player) RenderCards(dst *ebiten.Image, x, y, theta float64) {
	var speed, rSpeed float64 = 0.5, 0.2
	var xOff, yOff, rotate float64
	var s *view.Sprite

	if p.Dead {
		s = p.DeadSprite
		s.Rotate(theta, speed)
		s.Move(x, y, speed)
		s.Display(dst)
	} else {
		if p.ShieldCard != nil {
			s = p.ShieldCard.SSprite
			xOff, yOff, rotate = p.getShieldOffset()
			s.Rotate(rotate, rSpeed)
			s.MoveOffset(xOff, yOff, speed)
			s.RotateOffset(theta, rSpeed)
			s.Move(x, y, speed)
			s.Display(dst)
		}
	
		if p.JokerShield != nil {
			s = p.JokerShield.SSprite
			xOff, yOff, rotate = p.getJokerShieldOffset()
			s.Rotate(rotate, rSpeed)
			s.MoveOffset(xOff, yOff, speed)
			s.RotateOffset(theta, rSpeed)
			s.Move(x, y, speed)
			s.Display(dst)
		}
	
		if p.JokerHealth != nil {
			s = p.JokerHealth.SSprite
			xOff, yOff, rotate = p.getJokerHealthOffset()
			s.Rotate(rotate, rSpeed)
			s.MoveOffset(xOff, yOff, speed)
			s.RotateOffset(theta, rSpeed)
			s.Move(x, y, speed)
			s.Display(dst)
		}
	
		if p.HealthCard[0] != nil {
			s = p.HealthCard[0].SSprite
			xOff, yOff, rotate = p.getHealthOffset(0)
			s.Rotate(rotate, rSpeed)
			s.MoveOffset(xOff, yOff, speed)
			s.RotateOffset(theta, rSpeed)
			s.Move(x, y, speed)
			s.Display(dst)
		}
	
		if p.HealthCard[1] != nil {
			s = p.HealthCard[1].SSprite
			xOff, yOff, rotate = p.getHealthOffset(1)
			s.Rotate(rotate, rSpeed)
			s.MoveOffset(xOff, yOff, speed)
			s.RotateOffset(theta, rSpeed)
			s.Move(x, y, speed)
			s.Display(dst)
		}
	
		if p.ChargeCard != nil {
			s = p.ChargeCard.SSprite
			xOff, yOff, rotate = p.getChargeOffset()
			s.Rotate(rotate, rSpeed)
			s.MoveOffset(xOff, yOff, speed)
			s.RotateOffset(theta, rSpeed)
			s.Move(x, y, speed)
			s.Display(dst)
		}
	}

	s = p.NameSprite
	s.MoveOffset(0, float64(view.CardHeight) * 3/2, speed)
	s.RotateOffset(theta, rSpeed)
	s.Move(x, y, speed)
	s.Display(dst)
}