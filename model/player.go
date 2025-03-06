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
	HealthCard [2]*Card
	JokerHealth *Card
	ShieldCard *Card
	JokerShield *Card
	ChargeCard *Card
}

func NewPlayer(name string) *Player {
	var tWidth, tHeight float64 = text.Measure(name, view.TextFace, 0.0)
	var img *ebiten.Image = ebiten.NewImage(int(tWidth), int(tHeight))

	op := &text.DrawOptions{}; op.ColorScale.ScaleWithColor(color.RGBA{0,0,0,255})
	text.Draw(img, name, &text.GoTextFace{Source: view.FaceSource, Size: view.FontSize}, op)
	var s *view.Sprite = view.NewSprite(img, true, color.RGBA{75,75,75,127}, nil)

	return &Player{name, s, false, [2]*Card{nil, nil}, nil, nil, nil, nil}
}

// @desc: Player attacks an enemy with a given card and the one in charge, the total attack value and the charge card are returned
func (p *Player) Attack(c *Card) (int, *Card) {
	var attack int = c.Value

	var charge *Card = nil
	if p.ChargeCard != nil {
		charge = p.Uncharge()
		attack = attack + charge.Value
	}

	return attack, charge
}

// @desc: Player loose health based on the "attack" value parameter
//        "at" parameter specifies which card to lose health from in priority
//        then returns the new health value to be set and
//        all cards lost in the process of the attack (Maximum 3: joker health and both health cards)
func (p *Player) LoseHealth(attack int, at int) (int, *Card, *Card, *Card) {
	if p.JokerShield != nil {
		var c *Card = p.JokerShield
		p.JokerShield = nil
		return 0, c, nil, nil
	}

	var joker, health1, health2 *Card = nil, nil, nil

	// Attack is bigger than shield, we loose HP
	if p.ShieldCard != nil { attack = attack - p.ShieldCard.Value }
	if attack > 0 {

		// Do we have a joker health ? Then it's tanking (wether you like it or not)
		if p.JokerHealth != nil {
			attack = attack - p.JokerHealth.Value
			joker = p.JokerHealth
			p.JokerHealth = nil
		}

		// Is the attack still going ?
		if attack > 0 {
			attack = attack - p.HealthCard[at].Value
			health1 = p.HealthCard[at]
			p.HealthCard[at] = nil
		}

		// Wow that's a really big hit
		if attack > 0 && p.HealthCard[1-at] != nil {
			attack = attack - p.HealthCard[1-at].Value
			health2 = p.HealthCard[1-at]
			p.HealthCard[1-at] = nil
		}

		// R.I.P in Peperonni
		if attack >= 0 { p.Dead = true }
	}

	return -1*attack, joker, health1, health2
}

// @desc: Swap shield with the given card then returns the old shield
func (p *Player) Shield(c *Card) *Card {
	c, p.ShieldCard = p.ShieldCard, c
	return c
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

func (p *Player) Render(dst *ebiten.Image, x float64, y float64, r float64) {
	var s *view.Sprite = nil

	if p.Dead {
		s = view.NewSprite(view.GraveImage, false, color.RGBA{0,0,0,0}, nil)
		s.CenterImg()
		s.RotateImg(r)
		s.MoveImg(x, y)
		s.Display(dst)
	} else {
		if p.ShieldCard != nil {
			s = p.ShieldCard.SSprite
			s.ResetGeoM()
			s.CenterImg()
			s.RotateImg(math.Pi/2)
			s.MoveImg(0, 0 - s.Width/2)
			s.RotateImg(r)
			s.MoveImg(x, y)
			s.Display(dst)
		}
	
		if p.JokerShield != nil {
			s.ResetGeoM()
			s.CenterImg()
			s.RotateImg(math.Pi/2)
			s.MoveImg(0, 0 - s.Width/2 - 15)
			s.RotateImg(r)
			s.MoveImg(x, y)
			s.Display(dst)
		}
	
		if p.JokerHealth != nil {
			s = p.JokerHealth.SSprite
			s.ResetGeoM()
			s.MoveImg(-s.Width*2, 0)
			s.RotateImg(r)
			s.MoveImg(x, y)
			s.Display(dst)
		}
	
		if p.HealthCard[0] != nil {
			s = p.HealthCard[0].SSprite
			s.ResetGeoM()
			s.MoveImg(-s.Width, 0)
			s.RotateImg(r)
			s.MoveImg(x, y)
			s.Display(dst)
		}
	
		if p.HealthCard[1] != nil {
			s = p.HealthCard[1].SSprite
			s.ResetGeoM()
			s.MoveImg(0, 0)
			s.RotateImg(r)
			s.MoveImg(x, y)
			s.Display(dst)
		}
	
		if p.ChargeCard != nil {
			s = p.ChargeCard.SSprite
			s.ResetGeoM()
			s.MoveImg(-s.Width/2, s.Height)
			s.RotateImg(r)
			s.MoveImg(x, y)
			s.Display(dst)
		}
	}

	s = p.NameSprite
	s.ResetGeoM()
	s.CenterImg()
	s.MoveImg(0, float64(view.TileHeight) * 8/5)
	s.RotateImg(r)
	s.MoveImg(x, y)
	s.Display(dst)
}