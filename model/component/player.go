package game

import (
	"log"
	"math"

	"strconv"

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

	p.DeadSprite = view.NewSprite(view.LoadDeathImage(), false, color.RGBA{0,0,0,0}, nil)
	p.Dead = false

	p.Persona = char

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

func (p *Player) setCard(c *Card, x, y, r, xOff, yOff, rOff float64) {
	if c != nil {
		var speed, rSpeed float64 = 0.5, 0.2
		var s *view.Sprite = c.SSprite
		s.Rotate(r, rSpeed)
		s.MoveOffset(xOff, yOff, speed)
		s.RotateOffset(rOff, rSpeed)
		s.Move(x, y, speed)
	} else { log.Fatal("[player.setCard] Cannot set nil card") }
}

func (p *Player) setShield(x, y, rOff float64) {
	if p.ShieldCard == nil { log.Fatal("[player.setShield] ShieldCard is nil") }
	var xOff, yOff, r float64 = 0, -float64(view.CardWidth)/2, math.Pi/2
	p.setCard(p.ShieldCard, x, y, r, xOff, yOff, rOff)
}

func (p *Player) setJokerShield(x, y, rOff float64) {
	if p.JokerShield == nil { log.Fatal("[player.setJokerShield] JokerShield is nil") }
	var xOff, yOff, r float64 = 0, -float64(view.CardWidth)/2 - 15, math.Pi/2
	p.setCard(p.JokerShield, x, y, r, xOff, yOff, rOff)
}

func (p *Player) setJokerHealth(x, y, rOff float64) {
	if p.JokerHealth == nil { log.Fatal("[player.setJokerHealth] JokerHealth is nil") }
	var xOff, yOff, r float64 = - float64(view.CardWidth) - float64(view.CardWidth)/2, float64(view.CardHeight)/2, 0
	p.setCard(p.JokerHealth, x, y, r, xOff, yOff, rOff)
}

func (p *Player) setHealth(i int, x, y, rOff float64) {
	if p.HealthCard[i] == nil { log.Fatal("[player.setHealth] HealthCard" + strconv.Itoa(i) + " is nil") }
	var xOff float64 = float64((i-1) * view.CardWidth) + float64(view.CardWidth)/2
	var yOff, r float64 = float64(view.CardHeight)/2, 0
	p.setCard(p.HealthCard[i], x, y, r, xOff, yOff, rOff)
}

func (p *Player) setCharge(x, y, rOff float64) {
	if p.ChargeCard == nil { log.Fatal("[player.setCharge] ChargeCard is nil") }
	var xOff, yOff, r float64 = float64(view.CardWidth) + float64(view.CardWidth)/2, float64(view.CardHeight)/2, 0
	p.setCard(p.ChargeCard, x, y, r, xOff, yOff, rOff)
}

func (p *Player) RenderCards(dst *ebiten.Image, x, y, r float64) {
	if p.Dead {
		p.DeadSprite.Rotate(r, 0.2)
		p.DeadSprite.Move(x, y, 0.5)
		p.DeadSprite.Display(dst)
	} else {
		if p.ShieldCard != nil {
			p.setShield(x, y, r)
			p.ShieldCard.SSprite.Display(dst)
		}

		if p.JokerShield != nil {
			p.setJokerShield(x, y, r)
			p.JokerShield.SSprite.Display(dst)
		}
	
		if p.JokerHealth != nil {
			p.setJokerHealth(x, y, r)
			p.JokerHealth.SSprite.Display(dst)
		}

		if p.HealthCard[0] != nil {
			p.setHealth(0, x, y, r)
			p.HealthCard[0].SSprite.Display(dst)
		}

		if p.HealthCard[1] != nil {
			p.setHealth(1, x, y, r)
			p.HealthCard[1].SSprite.Display(dst)
		}

		if p.ChargeCard != nil {
			p.setCharge(x, y, r)
			p.ChargeCard.SSprite.Display(dst)
		}
	}

	p.NameSprite.MoveOffset(0, float64(view.CardHeight) * 3/2, 0.5)
	p.NameSprite.RotateOffset(r, 0.2)
	p.NameSprite.Move(x, y, 0.5)
	p.NameSprite.Display(dst)
}
