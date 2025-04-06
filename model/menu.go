package model

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"camaretto/model/component"
	"camaretto/view"
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

	local, join, host *component.Button

	Name *component.TextCapture

	Online bool
	Hosting bool

	goToLobby func()
	startServer func()
	joinServer func()
	scanServer func()
}

func (menu *Menu) Init(w, h int, lobby, host, join, scan func()) {
	menu.state = HOME

	menu.width, menu.height = float64(w), float64(h)

	var x, y float64 = menu.width/2, menu.height/2
	menu.local = component.NewButton("Local", color.RGBA{0, 0, 0, 255}, "GREEN", lobby)
	menu.local.SSprite.SetCenter(x, y - float64(view.ButtonHeight) - 5, 0)

	menu.host = component.NewButton("Host", color.RGBA{0, 0, 0, 255}, "BLUE", menu.hostGame)
	menu.host.SSprite.SetCenter(x, y, 0)

	menu.join = component.NewButton("Join", color.RGBA{0, 0, 0, 255}, "RED", menu.joinGame)
	menu.join.SSprite.SetCenter(x, y + float64(view.ButtonHeight) + 5, 0)

	menu.Name = component.NewTextCapture(55, int(menu.width*3/4), int(menu.height/10), 2)
	menu.Name.SSprite.SetCenter(x, y, 0)

	menu.Online = false
	menu.Hosting = false

	menu.goToLobby = lobby
	menu.startServer = host
	menu.joinServer = join
	menu.scanServer = scan
}

func (menu *Menu) hostGame() {
	menu.state = JOIN
	go menu.startServer()

	menu.host.SSprite.Move(menu.width/2, menu.height/2 + float64(view.ButtonHeight)*2, 2)
	menu.host.Trigger = menu.goToLobby
}

func (menu *Menu) joinGame() {
	menu.state = JOIN

	menu.join.SSprite.Move(menu.width/2, menu.height/2 + float64(view.ButtonHeight)*2, 2)
	menu.join.Trigger = menu.goToLobby
}

func (menu *Menu) scanGames() {
}

func (menu *Menu) Update() error {
	menu.local.Update()
	menu.host.Update()
	menu.join.Update()

	if menu.state == JOIN {
		menu.Name.Update()
	}

	return nil
}

/*
func (menu *Menu) MousePress(x, y float64)  {
	if menu.state == HOME {
		if menu.local.SSprite.In(x, y) {
			// menu.local.pressed()
		} else if menu.join.SSprite.In(x, y) {
			// menu.join.pressed()
		} else if menu.host.SSprite.In(x, y) {
			// menu.host.pressed()
		}
	} else if menu.state == SCAN {
	} else if menu.state == JOIN {
		if menu.Hosting && menu.host.SSprite.In(x, y) {
			// menu.host.pressed()
		} else if menu.join.SSprite.In(x, y) {
			// menu.join.pressed()
		}
	}
}

func (menu *Menu) MouseRelease(x, y float64)  {
	if menu.state == HOME {
		// menu.local.released()
		// menu.host.released()
		// menu.join.released()

		if menu.local.SSprite.In(x, y) {
			menu.Online = false

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
			// menu.host.released()
			// if menu.host.SSprite.In(x, y) { return NEXT }
		} else {
			// menu.join.released()
			// if menu.join.SSprite.In(x, y) { return NEXT }
		}
	}

}

func (menu *Menu) HandleKeyEvent(e *event.KeyEvent)  {
	if e.Event == event.PRESSED && menu.state == JOIN {
		// menu.Name.HandleEvent(e, nil)
	}
}
*/

func (menu *Menu) Draw(dst *ebiten.Image) {
	if menu.state == HOME {
		menu.local.Draw(dst)
		menu.join.Draw(dst)
		menu.host.Draw(dst)
	} else if menu.state == SCAN {
	} else if menu.state == JOIN {
		menu.Name.SSprite.Draw(dst)
		if menu.Hosting {
			menu.host.Draw(dst)
		} else {
			menu.join.Draw(dst)
		}
	}
}
