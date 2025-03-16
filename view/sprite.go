package view

import (
	"strconv"
	"math"

	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
)

type Sprite struct {
	Width, Height float64
	xCenter, yCenter float64
	rCenter float64
	targetXCenter, targetYCenter float64
	targetRCenter float64
	speedCenter, rSpeedCenter float64

	xOffset, yOffset float64
	rOffset float64
	targetXOffset, targetYOffset float64
	targetROffset float64
	speedOffset, rSpeedOffset float64

	scaleX, scaleY float64

	// bg bool
	// bgColor color.RGBA

	options *ebiten.DrawImageOptions
	// background *ebiten.Image
	image *ebiten.Image
}

// @desc: Init a sprite struct then returns it
func NewSprite(img *ebiten.Image, bgEnabled bool, c color.RGBA, op *ebiten.DrawImageOptions) *Sprite {
	var s *Sprite = &Sprite{}

	var w, h int = img.Size()
	s.image = ebiten.NewImage(w, h)
	s.image.DrawImage(img, nil)

	s.Width, s.Height = float64(w), float64(h)
	// s.bg, s.bgColor = bgEnabled, c

	if op == nil { op = &ebiten.DrawImageOptions{} }
	s.options = op

	s.xCenter, s.yCenter, s.rCenter = 0, 0, 0
	s.targetXCenter, s.targetYCenter, s.targetRCenter = 0, 0, 0
	s.speedCenter, s.rSpeedCenter = 0, 0

	s.xOffset, s.yOffset, s.rOffset = 0, 0, 0
	s.targetXOffset, s.targetYOffset, s.targetROffset = 0, 0, 0
	s.speedOffset, s.rSpeedOffset = 0, 0

	s.scaleX, s.scaleY = 1, 1

	return s
}

func (s *Sprite) SetImage(img *ebiten.Image) {
	s.image.Clear()
	s.image.DrawImage(img, nil)
}

func (s *Sprite) SetCenter(x, y, r float64) {
	s.xCenter, s.yCenter, s.rCenter = x, y, r
	s.targetXCenter, s.targetYCenter, s.targetRCenter = x, y, r
	s.options.GeoM.Reset()
}
func (s *Sprite) GetCenter() (float64, float64, float64) { return s.xCenter, s.yCenter, s.rCenter }

func (s *Sprite) SetOffset(x, y, r float64) {
	s.xOffset, s.yOffset, s.rOffset = x, y, r
	s.targetXOffset, s.targetYOffset, s.targetROffset = x, y, r
	s.options.GeoM.Reset()
}
func (s *Sprite) GetOffset() (float64, float64, float64) { return s.xOffset, s.yOffset, s.rOffset }

// @desc: Returns true if the coordinates (x,y) are within the sprite, false otherwise
func (s *Sprite) In(x, y float64) bool {
	var inv ebiten.GeoM = s.options.GeoM
	inv.Invert()
	x, y = inv.Apply(x, y)

	if x < 0 || x > s.Width { return false }
	if y < 0 || y > s.Height { return false }
	return true
}

func (s *Sprite) Scale(x, y float64) {
	s.scaleX, s.scaleY = x, y
	s.Width, s.Height = s.Width * x, s.Height * y
}

func (s *Sprite) Move(x, y, sp float64) { s.targetXCenter, s.targetYCenter, s.speedCenter = x, y, sp }
func (s *Sprite) Rotate(r, sp float64) { s.targetRCenter, s.rSpeedCenter = r, sp }

func (s *Sprite) MoveOffset(x, y, sp float64) { s.targetXOffset, s.targetYOffset, s.speedOffset = x, y, sp }
func (s *Sprite) RotateOffset(r, sp float64) { s.targetROffset, s.rSpeedOffset = math.Mod(r, math.Pi*2), sp }

func (s *Sprite) tickTranslateCenter() {
	var dx, dy float64 = (s.targetXCenter - s.xCenter), (s.targetYCenter - s.yCenter)

	var vx, vy float64 = dx * s.speedCenter/50, dy * s.speedCenter/50

	if math.Abs(dx) < 1 {
		s.xCenter = s.targetXCenter
	} else {
		s.xCenter = s.xCenter + vx
	}

	if math.Abs(dy) < 1 {
		s.yCenter = s.targetYCenter
	} else {
		s.yCenter = s.yCenter + vy
	}
}

func (s *Sprite) tickRotateCenter() {
	var vr float64 = s.targetRCenter * s.rSpeedCenter/50
	if s.targetRCenter - s.rCenter < math.Pi/90 {
		s.rCenter = s.targetRCenter
	} else {
		s.rCenter = s.rCenter + vr
	}
}

func (s *Sprite) tickTranslateOffset() {
	var dx, dy float64 = (s.targetXOffset - s.xOffset), (s.targetYOffset - s.yOffset)

	var vx, vy float64 = dx * s.speedOffset/50, dy * s.speedOffset/50

	if math.Abs(dx) < 1 {
		s.xOffset = s.targetXOffset
	} else {
		s.xOffset = s.xOffset + vx
	}

	if math.Abs(dy) < 1 {
		s.yOffset = s.targetYOffset
	} else {
		s.yOffset = s.yOffset + vy
	}
}

func (s *Sprite) tickRotateOffset() {
	// var vr float64 = s.targetROffset * s.rSpeedOffset/50
	// if s.targetROffset - s.rOffset < math.Pi/90 {
	// 	s.rOffset = s.targetROffset
	// } else {
	// 	s.rOffset = s.rOffset + vr
	// }
	s.rOffset = s.targetROffset
}

func (s *Sprite) tick() {
	if s.xCenter != s.targetXCenter || s.yCenter != s.targetYCenter { s.tickTranslateCenter() }
	if s.rCenter != s.targetRCenter { s.tickRotateCenter() }

	if s.xOffset != s.targetXOffset || s.yOffset != s.targetYOffset { s.tickTranslateOffset() }
	if s.rOffset != s.targetROffset { s.tickRotateOffset() }
}

func (s *Sprite) Display(dst *ebiten.Image) {
	s.tick()

	s.options.GeoM.Reset()
	s.options.GeoM.Scale(s.scaleX, s.scaleY)
	s.options.GeoM.Translate(-s.Width/2, -s.Height/2) // Center img
	s.options.GeoM.Rotate(s.rCenter) // Rotate in place
	s.options.GeoM.Translate(s.xOffset, s.yOffset) // Offset img
	s.options.GeoM.Rotate(s.rOffset) // Apply offset rotation
	s.options.GeoM.Translate(s.xCenter, s.yCenter) // Put img in place

	dst.DrawImage(s.image, s.options)
}

func (s *Sprite) ToString() string {
	msg := "DX" + strconv.FormatFloat((s.targetXCenter - s.xCenter), 'f', 3, 64)
	msg = msg + ", DY" + strconv.FormatFloat((s.targetYCenter - s.yCenter), 'f', 3, 64)
	return msg
}