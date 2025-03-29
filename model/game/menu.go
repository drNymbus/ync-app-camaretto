package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"camaretto/model/component"
	"camaretto/view"
	"camaretto/event"
)

type MenuState int
const (
	HOME MenuState = iota
	SCAN
	JOIN
)

type Menu struct {
	state MenuState

	width, height float64

	NbPlayers int
	minusButton, plusButton *component.Button

	local, join, host *component.Button

	Name *component.TextCapture

	Online bool
	Hosting bool
}

func (menu *Menu) Init(w, h int) {
	menu.state = HOME

	menu.NbPlayers = 2

	menu.width, menu.height = float64(w), float64(h)

	var x, y float64 = menu.width/2, menu.height/8
	menu.minusButton = component.NewButton("-", color.RGBA{0, 0, 0, 255}, "RED")
	menu.minusButton.SSprite.SetCenter(x - float64(view.ButtonWidth)/2 - 5, y, 0)

	menu.plusButton = component.NewButton("+", color.RGBA{0, 0, 0, 255}, "RED")
	menu.plusButton.SSprite.SetCenter(x + float64(view.ButtonWidth)/2 + 5, y, 0)

	y = menu.height/2
	menu.local = component.NewButton("Local", color.RGBA{0, 0, 0, 255}, "GREEN")
	menu.local.SSprite.SetCenter(x, y - float64(view.ButtonHeight) - 5, 0)

	menu.host = component.NewButton("Host", color.RGBA{0, 0, 0, 255}, "BLUE")
	menu.host.SSprite.SetCenter(x, y, 0)

	menu.join = component.NewButton("Join", color.RGBA{0, 0, 0, 255}, "RED")
	menu.join.SSprite.SetCenter(x, y + float64(view.ButtonHeight) + 5, 0)

	menu.Name = component.NewTextCapture(55, int(menu.width*3/4), int(menu.height/10), 2)
	menu.Name.SSprite.SetCenter(x, y, 0)

	menu.Online = false
	menu.Hosting = false
}

func (menu *Menu) MousePress(x, y float64) component.PageSignal {
	if menu.state == HOME {
		if menu.local.SSprite.In(x, y) {
			menu.local.Pressed()
		} else if menu.join.SSprite.In(x, y) {
			menu.join.Pressed()
		} else if menu.host.SSprite.In(x, y) {
			menu.host.Pressed()
		}
	} else if menu.state == SCAN {
	} else if menu.state == JOIN {
		if menu.Hosting && menu.host.SSprite.In(x, y) {
			menu.host.Pressed()
		} else if menu.join.SSprite.In(x, y) {
			menu.join.Pressed()
		}
	}
	return component.UPDATE
}

func (menu *Menu) MouseRelease(x, y float64) component.PageSignal {
	if menu.state == HOME {
		menu.local.Released()
		menu.host.Released()
		menu.join.Released()

		if menu.local.SSprite.In(x, y) {
			menu.Online = false
			// GO TO LOBBY
			return component.NEXT

		} else if menu.host.SSprite.In(x, y) {
			menu.Online = true
			menu.Hosting = true

			var x, y float64 = menu.width/2, menu.height/8
			menu.Name.SSprite.SetCenter(x, menu.height/2 - y/2, 0)
			menu.host.SSprite.Move(x, menu.height/2 + float64(view.ButtonHeight)*2, 0.3)

			menu.state = JOIN

		} else if menu.join.SSprite.In(x, y) {
			menu.Online = true
			menu.Hosting = false

			var x, y float64 = menu.width/2, menu.height/8
			menu.Name.SSprite.SetCenter(x, menu.height/2 - y/2, 0)
			menu.join.SSprite.Move(x, menu.height/2 + float64(view.ButtonHeight)*2, 0.3)

			// GO TO SCAN
			menu.state = JOIN
		}
	} else if menu.state == SCAN {
	} else if menu.state == JOIN {
		if menu.Hosting {
			menu.host.Released()
			if menu.host.SSprite.In(x, y) { return component.NEXT }
		} else {
			menu.join.Released()
			if menu.join.SSprite.In(x, y) { return component.NEXT }
		}
	}

	return component.UPDATE
}

func (menu *Menu) HandleKeyEvent(e *event.KeyEvent) component.PageSignal {
	if e.Event == event.PRESSED && menu.state == JOIN {
		menu.Name.HandleEvent(e, nil)
	}
	return component.UPDATE
}

func (menu *Menu) Display(dst *ebiten.Image) {
	if menu.state == HOME {
		menu.local.SSprite.Display(dst)
		menu.join.SSprite.Display(dst)
		menu.host.SSprite.Display(dst)
	} else if menu.state == SCAN {
	} else if menu.state == JOIN {
		menu.Name.SSprite.Display(dst)
		if menu.Hosting {
			menu.host.SSprite.Display(dst)
		} else {
			menu.join.SSprite.Display(dst)
		}
	}
}
