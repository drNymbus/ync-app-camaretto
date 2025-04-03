package model

import (
	"math"

	"strconv"

	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"camaretto/model/component"
	"camaretto/view"
	"camaretto/event"
)

const (
	MaxNbPlayers int = 6
)

type Lobby struct {
	width, height float64
	online, hosting bool

	Names []*component.TextCapture

	Focus int
	cursor *view.Sprite

	NbPlayers int
	minusButton, plusButton *component.Button

	start *component.Button
}

func (lobby *Lobby) Init(w, h int, online, host bool) {
	lobby.width, lobby.height = float64(w), float64(h)
	lobby.online, lobby.hosting = online, host

	lobby.NbPlayers = 2
	lobby.Names = make([]*component.TextCapture, MaxNbPlayers)

	var tcWidth, tcHeight float64 = lobby.width*3/4, lobby.height/10
	for i := 0; i < MaxNbPlayers; i++ {
		lobby.Names[i] = component.NewTextCapture(55, int(tcWidth), int(tcHeight), 2)
		var diffY float64 = float64(i - MaxNbPlayers/2) * tcHeight + float64(i*10)
		lobby.Names[i].SSprite.SetCenter(lobby.width/2, lobby.height/2 + 50 + diffY, 0)
	}

	lobby.Focus = 0
	lobby.cursor = view.NewSprite(view.LoadCursorImage(), false, color.RGBA{0, 0, 0, 0}, nil)

	var x, y float64 = lobby.width/2, lobby.height/8
	lobby.minusButton = component.NewButton("-", color.RGBA{0, 0, 0, 255}, "RED")
	lobby.minusButton.SSprite.SetCenter(x - float64(view.ButtonWidth)/2 - 5, y, 0)

	lobby.plusButton = component.NewButton("+", color.RGBA{0, 0, 0, 255}, "RED")
	lobby.plusButton.SSprite.SetCenter(x + float64(view.ButtonWidth)/2 + 5, y, 0)

	lobby.start = component.NewButton("START", color.RGBA{0, 0, 0, 255}, "GREEN")
	lobby.start.SSprite.SetCenter(lobby.width/2, lobby.height - float64(view.ButtonHeight), 0)
}

func (lobby *Lobby) MousePress(x, y float64) component.PageSignal {
	if lobby.plusButton.SSprite.In(x, y) {
		lobby.plusButton.Pressed()
	} else if lobby.minusButton.SSprite.In(x, y) {
		lobby.minusButton.Pressed()
	} else if lobby.start.SSprite.In(x, y) {
		lobby.start.Pressed()
	} else if !lobby.online {
		for i, textInput := range lobby.Names {
			if textInput.SSprite.In(x, y) { lobby.Focus = i }
		}
	}
	return component.UPDATE
}

func (lobby *Lobby) MouseRelease(x, y float64) component.PageSignal {
	lobby.plusButton.Released()
	lobby.minusButton.Released()
	lobby.start.Released()
	if !lobby.online {
		if lobby.plusButton.SSprite.In(x, y) {
			if lobby.NbPlayers < 6 { lobby.NbPlayers++ }
		} else if lobby.minusButton.SSprite.In(x, y) {
			if lobby.NbPlayers > 2 { lobby.NbPlayers-- }
		} else if lobby.start.SSprite.In(x, y) {
			return component.NEXT
		}
	} else if lobby.hosting {
		lobby.start.Released()
		if lobby.start.SSprite.In(x, y) {
			return component.NEXT
		}
	}
	return component.UPDATE
}

func (lobby *Lobby) HandleKeyEvent(e *event.KeyEvent) component.PageSignal {
	if e.Event == event.PRESSED && !lobby.online {
		lobby.Names[lobby.Focus].HandleEvent(e, nil)
	}
	return component.UPDATE
}

func (lobby *Lobby) Display(dst *ebiten.Image) {
	if !lobby.online {
		lobby.minusButton.SSprite.Display(dst)
		lobby.plusButton.SSprite.Display(dst)
	}

	if !lobby.online || lobby.hosting {
		lobby.start.SSprite.Display(dst)
	}

	var x, y float64 = lobby.width/2, lobby.height/8 - float64(view.ButtonHeight)/2

	var textImg *ebiten.Image
	var tw, th float64
	textImg, tw, th = view.TextToImage(strconv.Itoa(lobby.NbPlayers), color.RGBA{0, 0, 0, 255})
	op := &ebiten.DrawImageOptions{}; op.GeoM.Translate(x - tw/2, y - th)
	dst.DrawImage(textImg, op)

	var tcWidth, tcHeight float64 = lobby.width*3/4, lobby.height/10
	for i := 0; i < lobby.NbPlayers; i++ {
		var diffY float64 = float64(i - MaxNbPlayers/2) * tcHeight + float64(i*10)
		lobby.Names[i].SSprite.SetCenter(lobby.width/2, lobby.height/2 + 50 + diffY, 0)
		lobby.Names[i].SSprite.Display(dst)
	}

	x = lobby.width/2 - tcWidth/2
	y = lobby.height/2 + float64(lobby.Focus - MaxNbPlayers/2) * tcHeight + float64(lobby.Focus * 10) + 50
	lobby.cursor.Move(x, y, 1)
	lobby.cursor.Rotate(math.Pi/2, 1)
	lobby.cursor.Display(dst)
}
