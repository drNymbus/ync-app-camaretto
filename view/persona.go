package view

import (
	// "log"

	// "os"
	// "io"
	// "bytes"
	// "image"
	// "image/color"

	// "golang.org/x/text/language"

	"github.com/hajimehoshi/ebiten/v2"
	// "github.com/hajimehoshi/ebiten/v2/ebitenutil"
	// "github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	OriginalPersonaBodyWidth int = 1640
	OriginalPersonaBodyHeight int = 2360
	OriginalPersonaMouthWidth int = 635
	OriginalPersonaMouthHeight int = 326
	PersonaBodyWidth int = 410
	PersonaBodyHeight int = 590
	PersonaMouthWidth int = 159
	PersonaMouthHeight int = 82
)

type PersonaImage struct {
	Body *ebiten.Image
	OpenMouth *ebiten.Image
	ClosedMouth *ebiten.Image
}

func LoadPersonaImage(name string) *PersonaImage {
	var pi *PersonaImage = &PersonaImage{}

	pi.Body = GetImage("assets/characters/char.png")
	pi.OpenMouth = GetImage("assets/characters/mouth_open.png")
	pi.ClosedMouth = GetImage("assets/characters/mouth_closed.png")

	var scale float64 = 0.25
	var op *ebiten.DrawImageOptions = &ebiten.DrawImageOptions{}
	op.GeoM.Scale(scale, scale)

	var tmp *ebiten.Image = GetImage("assets/characters/char.png")
	// PersonaBodyWidth = int(float64(OriginalPersonaBodyWidth)*scale)
	// PersonaBodyHeight = int(float64(OriginalPersonaBodyHeight)*scale)
	pi.Body = ebiten.NewImage(PersonaBodyWidth, PersonaBodyHeight)
	pi.Body.DrawImage(tmp, op)

	// PersonaMouthWidth = int(float64(PersonaMouthWidth)*scale)
	// PersonaMouthHeight = int(float64(PersonaMouthHeight)*scale)

	tmp = GetImage("assets/characters/mouth_open.png")
	pi.OpenMouth = ebiten.NewImage(PersonaMouthWidth, PersonaBodyHeight)
	pi.OpenMouth.DrawImage(tmp, op)

	tmp = GetImage("assets/characters/mouth_closed.png")
	pi.ClosedMouth = ebiten.NewImage(PersonaMouthWidth, PersonaBodyHeight)
	pi.ClosedMouth.DrawImage(tmp, op)

	return pi
}

func LoadDeathImage() *ebiten.Image {
	// Death Sprite
	var tmp *ebiten.Image = GetImage("assets/characters/jesus.jpg")

	var width, height int = tmp.Size()
	var xScale, yScale float64 = 0.1, 0.1
	width, height = int(float64(width) * xScale), int(float64(height) * yScale)

	var d *ebiten.Image = ebiten.NewImage(width, height)
	op := &ebiten.DrawImageOptions{}; op.GeoM.Scale(xScale, yScale)
	d.DrawImage(tmp, op)

	return d
}