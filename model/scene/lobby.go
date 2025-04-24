package scene

import (
	"log"
	"math"

	"strconv"

	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"

	"camaretto/model/ui"
	"camaretto/view"
)

const (
	MaxNbPlayers int = 6
)

type Lobby struct {
	width, height float64
	online, hosting bool

	Names []*ui.TextCapture

	focus int
	cursor *view.Sprite

	NbPlayers int
	minusButton, plusButton *ui.Button

	start *ui.Button
}

func (lobby *Lobby) Init(w, h int, online, host bool, startGame func()) {
	lobby.width, lobby.height = float64(w), float64(h)
	lobby.online, lobby.hosting = online, host

	lobby.NbPlayers = 2

	lobby.Names = make([]*ui.TextCapture, MaxNbPlayers)
	var tcWidth, tcHeight float64 = lobby.width*3/4, lobby.height/10
	for i := 0; i < MaxNbPlayers; i++ {
		lobby.Names[i] = ui.NewTextCapture(55, int(tcWidth), int(tcHeight), 2)
		var diffY float64 = float64(i - MaxNbPlayers/2) * tcHeight + float64(i*10)
		lobby.Names[i].SSprite.SetCenter(lobby.width/2, lobby.height/2 + 50 + diffY, 0)
		lobby.Names[i].Disable()
	}

	if !lobby.online { lobby.Names[0].Enable() }
	lobby.focus = 0
	lobby.cursor = view.NewSprite(view.LoadCursorImage(), nil)

	var x, y float64 = lobby.width/2, lobby.height/8
	lobby.minusButton = ui.NewButton("-", color.RGBA{0, 0, 0, 255}, "RED", lobby.removePlayer)
	lobby.minusButton.SSprite.SetCenter(x - float64(view.ButtonWidth)/2 - 5, y, 0)

	lobby.plusButton = ui.NewButton("+", color.RGBA{0, 0, 0, 255}, "RED", lobby.addPlayer)
	lobby.plusButton.SSprite.SetCenter(x + float64(view.ButtonWidth)/2 + 5, y, 0)

	lobby.start = ui.NewButton("START", color.RGBA{0, 0, 0, 255}, "GREEN", startGame)
	lobby.start.SSprite.SetCenter(lobby.width/2, lobby.height - float64(view.ButtonHeight), 0)
}

func (lobby *Lobby) addPlayer() {
	if lobby.NbPlayers < MaxNbPlayers {
		lobby.NbPlayers++
	}
}

func (lobby *Lobby) removePlayer() {
	if lobby.NbPlayers > 2 {
		lobby.NbPlayers--
	}
}

func (lobby *Lobby) handleError(err error, from string, action string) error {
	var msg string = "[Lobby." + from + "] " + action + ":"
	log.Println(msg, err)
	return err
}

func (lobby *Lobby) Update() error {
	var err error

	if !lobby.online {
		lobby.minusButton.Update(nil)
		if err != nil { return lobby.handleError(err, "Update", "Button minusButton.Update") }
		lobby.plusButton.Update(nil)
		if err != nil { return lobby.handleError(err, "Update", "Button plusButton.Update") }
	}

	if !lobby.online || lobby.hosting {
		lobby.start.Update(nil)
		if err != nil { return lobby.handleError(err, "Update", "Button start.Update") }
	}

	if !lobby.online && inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		var x, y int = ebiten.CursorPosition()

		for i, tc := range lobby.Names {
			if tc.SSprite.In(float64(x), float64(y)) {
				lobby.Names[lobby.focus].Disable()
				lobby.Names[i].Enable()
				lobby.focus = i
			}
		}
	}

	var tc *ui.TextCapture = lobby.Names[lobby.focus]
	var x, y float64
	x = lobby.width/2 - tc.SSprite.Width/2
	y = lobby.height/2 + float64(lobby.focus - MaxNbPlayers/2) * tc.SSprite.Height
	y = y + float64(lobby.focus * 10) + 50
	lobby.cursor.Move(x, y, 3)
	lobby.cursor.Rotate(math.Pi/2, 1)

	lobby.cursor.Update()

	for i := 0; i < lobby.NbPlayers; i++ {
		lobby.Names[i].Update(nil)
	}
	// lobby.Names[lobby.focus].Update()

	return nil
}

func (lobby *Lobby) Draw(screen *ebiten.Image) {
	if !lobby.online {
		lobby.minusButton.Draw(screen)
		lobby.plusButton.Draw(screen)
	}

	if !lobby.online || lobby.hosting {
		lobby.start.Draw(screen)
	}

	var x, y float64 = lobby.width/2, lobby.height/8 - float64(view.ButtonHeight)/2

	var textImg *ebiten.Image
	var tw, th float64
	textImg, tw, th = view.TextToImage(strconv.Itoa(lobby.NbPlayers), color.RGBA{0, 0, 0, 255})
	op := &ebiten.DrawImageOptions{}; op.GeoM.Translate(x - tw/2, y - th)
	screen.DrawImage(textImg, op)

	for i := 0; i < lobby.NbPlayers; i++ {
		lobby.Names[i].Draw(screen)
	}

	lobby.cursor.Draw(screen)
}
