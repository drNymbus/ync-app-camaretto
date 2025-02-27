package view

import (
	"image/color"
	"github.com/hajimehoshi/ebiten/v2"
)

type Sprite struct {
	Width, Height float64

	backgroundEnabled bool
	backgroundColor color.RGBA

	Options *ebiten.DrawImageOptions
	Background *ebiten.Image
	Img *ebiten.Image
}

// @desc: Init a sprite struct then returns it
func NewSprite(img *ebiten.Image, bgEnabled bool, c color.RGBA, op *ebiten.DrawImageOptions) *Sprite {
	var w, h int = img.Size()
	if op == nil { op = &ebiten.DrawImageOptions{} }

	var s *Sprite = &Sprite{float64(w), float64(h), bgEnabled, c, op, nil, img}
	if bgEnabled { s.RenderBackground() }
	return s
}

// @desc: Create an image the same size as the sprite then fills it with the registered color
func (s *Sprite) RenderBackground() {
	s.Background = ebiten.NewImage(int(s.Width), int(s.Height))
	s.Background.Fill(s.backgroundColor)
}

// @desc: Switch the flag to true indicating wether the background should be displayed or not
func (s *Sprite) EnableBackground() {
	s.backgroundEnabled = true
	s.RenderBackground()
}
// @desc: Switch the flag to false indicating wether the background should be displayed or not
func (s *Sprite) DisableBackground() { s.backgroundEnabled = false }

// @desc: Change the background color of the background image
func (s *Sprite) SetBackgroundColor(c color.RGBA) {
	s.backgroundColor = c
	s.RenderBackground()
}
func (s *Sprite) GetBackgroundColor() color.RGBA { return s.backgroundColor }

// @desc: Apply/Record a translation to the sprite's image
func (s *Sprite) MoveImg(x, y float64) {
	s.Options.GeoM.Translate(x, y)
}

// @desc: Calls MoveImg and translates it half the width and height
func (s *Sprite) CenterImg() {
	s.MoveImg(-s.Width/2, -s.Height/2)
}

// @desc: Apply/Record a rotation to the sprite's image
func (s *Sprite) RotateImg(r float64) {
	s.Options.GeoM.Rotate(r)
}

// @desc: Resets all modifications (translations & rotations) applied to the sprite's image
func (s *Sprite) ResetGeoM() {
	s.Options.GeoM.Reset()
}

// @desc: Returns true if the coordinates (x,y) are within the sprite, false otherwise
func (s *Sprite) In(x, y float64) bool {
	var inv ebiten.GeoM = s.Options.GeoM
	inv.Invert()
	x, y = inv.Apply(x, y)

	if x < 0 || x > s.Width { return false }
	if y < 0 || y > s.Height { return false }
	return true
}

// @desc: Draws the sprite onto the dst *ebiten.Image given
func (s *Sprite) Display(dst *ebiten.Image) {
	if s.backgroundEnabled { dst.DrawImage(s.Background, s.Options) }
	dst.DrawImage(s.Img, s.Options)
}