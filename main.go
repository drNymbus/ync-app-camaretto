package main

import (
	"log"
	"math"
	"strconv"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"camaretto/model"
	"camaretto/event"
	"camaretto/view"
)

var (
	err error
	// WinWidth int = 640
	// WinHeight int = 480
	WinWidth int = 1200
	WinHeight int = 900
)

type AppState int
const (
	MENU AppState = 0
	GAME AppState = 1
	END AppState = 2
)

type Game struct{
	state AppState
	camaretto *model.Camaretto

	bAttack *model.Button
	bShield *model.Button
	bCharge *model.Button
	bHeal *model.Button

	bInfo *model.Button

	mouse *event.Mouse
}

func (g *Game) Init(nbPlayers int) {
	g.state = GAME
	g.mouse = event.NewMouse(10)

	g.camaretto = model.NewCamaretto(nbPlayers)
	for i, player := range g.camaretto.Players {
		log.Println(strconv.Itoa(i), player.HealthCard[0], player.HealthCard[1], player.JokerHealth, player.ShieldCard)
	}

	g.bAttack = model.NewButton("ATTACK", color.RGBA{0, 0, 0, 255}, color.RGBA{163, 3, 9, 127})
	g.bShield = model.NewButton("SHIELD", color.RGBA{0, 0, 0, 255}, color.RGBA{2, 42, 201, 127})
	g.bCharge = model.NewButton("CHARGE", color.RGBA{0, 0, 0, 255}, color.RGBA{224, 144, 4, 127})
	g.bHeal = model.NewButton("HEAL", color.RGBA{0, 0, 0, 255}, color.RGBA{3, 173, 18, 127})

	g.bInfo = model.NewButton("This contains information.", color.RGBA{0, 0, 0, 255}, color.RGBA{0, 0, 0, 0})
	g.bInfo.SetMessage("PLAYER" + strconv.Itoa(g.camaretto.GetPlayerTurn()) + ": Choose an action.")
}

func (g *Game) Update() error {
	g.mouse.Update()
	g.bAttack.Hover(g.mouse.X, g.mouse.Y)
	// g.bShield.Hover(g.mouse.X, g.mouse.Y)
	// g.bCharge.Hover(g.mouse.X, g.mouse.Y)
	// g.bHeal.Hover(g.mouse.X, g.mouse.Y)

	// for _, player := range g.camaretto.Players {
	// 	if player.HealthCard[0] != nil { player.HealthCard[0].Hover(g.mouse.X, g.mouse.Y) }
	// 	if player.HealthCard[1] != nil { player.HealthCard[1].Hover(g.mouse.X, g.mouse.Y) }
	// 	if player.ShieldCard != nil { player.ShieldCard.Hover(g.mouse.X, g.mouse.Y) }
	// 	if player.ChargeCard != nil { player.ChargeCard.Hover(g.mouse.X, g.mouse.Y) }
	// }

	var state model.GameState = g.camaretto.GetState()
	var playerTurn int = g.camaretto.GetPlayerTurn()

	var focus model.FocusState = g.camaretto.GetFocus()
	var playerFocus int = g.camaretto.GetPlayerFocus()
	var cardFocus int = g.camaretto.GetCardFocus()

	var e *event.MouseEvent = nil
	for ;!g.mouse.EmptyEventQueue(); {
		e = g.mouse.ReadEvent()
		if e.MET == event.RELEASED && e.Click == ebiten.MouseButtonLeft {
			log.Println("UPDATE", e.X, e.Y)

			if state == model.SET {

				if g.bAttack.SSprite.In(e.X, e.Y) {
					log.Println("ATTACK")
					state = model.ATTACK
					focus = model.PLAYER
					g.bInfo.SetMessage("PLAYER" + strconv.Itoa(playerTurn) + ": Choose a player to attack.")
				} else if g.bShield.SSprite.In(e.X, e.Y) {
					log.Println("SHIELD")
					state = model.SHIELD
					focus = model.PLAYER
					g.bInfo.SetMessage("PLAYER" + strconv.Itoa(playerTurn) + ": Choose the shield to be switched.")
				} else if g.bCharge.SSprite.In(e.X, e.Y) {
					log.Println("CHARGE")
					state = model.CHARGE
					focus = model.COMPLETE
					g.bInfo.SetMessage("PLAYER" + strconv.Itoa(playerTurn) + ": Let's play ! Draw a card !")
				} else if g.bHeal.SSprite.In(e.X, e.Y) {
					log.Println("HEAL")
					state = model.HEAL
					playerFocus = playerTurn
					focus = model.CARD
					g.bInfo.SetMessage("PLAYER" + strconv.Itoa(playerTurn) + ": Choose a card of your own to switch.")
				}

			} else {

				if focus == model.PLAYER {

					var i int = event.HandleFocusPlayer(g.camaretto.Players, e)
					if i != -1 {
						playerFocus = i
						if state == model.ATTACK {
							focus = model.CARD
							g.bInfo.SetMessage("PLAYER" + strconv.Itoa(playerFocus) + ": Choose a card to defend against the attack.")
						} else if state == model.SHIELD {
							focus = model.COMPLETE
							g.bInfo.SetMessage("PLAYER" + strconv.Itoa(playerTurn) + ": Let's play ! Draw a card !")
						}
					}

				} else if focus == model.CARD {

					var i int = event.HandleFocusCard(g.camaretto.Players[playerFocus], e)
					if i != -1 {
						cardFocus = i
						focus = model.COMPLETE
						if state == model.HEAL {
							g.bInfo.SetMessage("PLAYER" + strconv.Itoa(playerTurn) + ": So ? What was in your charge ?")
						} else {
							g.bInfo.SetMessage("PLAYER" + strconv.Itoa(playerTurn) + ": Let's play ! Draw a card !")
						}
					}

				} else if focus == model.COMPLETE {

					if state == model.HEAL {
						if g.camaretto.Players[playerTurn].ChargeCard.SSprite.In(e.X, e.Y) {
							g.camaretto.Heal(playerTurn, cardFocus)
							state = model.SET; playerFocus = -1; cardFocus = -1
							g.camaretto.EndTurn()
							g.bInfo.SetMessage("PLAYER" + strconv.Itoa(g.camaretto.GetPlayerTurn()) + ": Choose an action.")
						}
					} else {
						if event.HandleFocusComplete(g.camaretto.DeckPile, e) {
							if state == model.ATTACK {
								g.camaretto.Attack(playerTurn, playerFocus, cardFocus)
							
								for i, player := range g.camaretto.Players {
									log.Println(strconv.Itoa(i), player.HealthCard[0], player.HealthCard[1], player.JokerHealth, player.ShieldCard)
								}
							} else if state == model.SHIELD {
								g.camaretto.Shield(playerFocus)
							} else if state == model.CHARGE {
								g.camaretto.Charge(playerTurn)
							}
							state = model.SET; playerFocus = -1; cardFocus = -1
							g.camaretto.EndTurn()
							g.bInfo.SetMessage("PLAYER" + strconv.Itoa(g.camaretto.GetPlayerTurn()) + ": Choose an action.")
						}
					}

				}
			}
		}
	}

	var code int; var message string;
	code, message = g.camaretto.SetState(state)
	if code == 1 { g.bInfo.SetMessage(message) }

	g.camaretto.SetFocus(focus)
	g.camaretto.SetPlayerFocus(playerFocus)
	g.camaretto.SetCardFocus(cardFocus)

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)

	// // Draw players
	// view.DrawPlayers(screen, g.camaretto.Players)
	var nbPlayers int = len(g.camaretto.Players)
	var angleStep float64 = 2*math.Pi / float64(nbPlayers)
	var radius float64 = 200

	var centerX float64 = float64(WinWidth)/2
	var centerY float64 = (float64(WinHeight) * 6/8)/2

	for i, player := range g.camaretto.Players {
		var theta float64 = angleStep * float64(i)
		var x float64 = centerX + (radius * math.Cos(theta + math.Pi/2))
		var y float64 = centerY + (radius * math.Sin(theta + math.Pi/2))
		player.Render(screen, x, y, theta)
	}

	// Draw deck pile
	// view.DrawDeck(screen, g.camaretto.DeckPile)
	var deck *model.Deck = g.camaretto.DeckPile
	for i, card := range deck.DrawPile[:deck.LenDrawPile] {
		card.SSprite.ResetGeoM()
		card.SSprite.MoveImg(centerX - card.SSprite.Width, centerY - float64(i)*0.2)
		card.SSprite.Display(screen)
	}
	for i, card := range deck.DiscardPile[:deck.LenDiscardPile] {
		card.SSprite.ResetGeoM()
		card.SSprite.MoveImg(centerX, centerY - float64(i)*0.2)
		card.SSprite.Display(screen)
	}

	// Draw buttons
	var buttonYPos float64 = float64(WinHeight)*7/8

	g.bAttack.SSprite.ResetGeoM()
	g.bAttack.SSprite.CenterImg()
	g.bAttack.SSprite.MoveImg(float64(WinWidth)*1/5, buttonYPos)
	g.bAttack.SSprite.Display(screen)

	g.bShield.SSprite.ResetGeoM()
	g.bShield.SSprite.CenterImg()
	g.bShield.SSprite.MoveImg(float64(WinWidth)*2/5, buttonYPos)
	g.bShield.SSprite.Display(screen)

	g.bCharge.SSprite.ResetGeoM()
	g.bCharge.SSprite.CenterImg()
	g.bCharge.SSprite.MoveImg(float64(WinWidth)*3/5, buttonYPos)
	g.bCharge.SSprite.Display(screen)

	g.bHeal.SSprite.ResetGeoM()
	g.bHeal.SSprite.CenterImg()
	g.bHeal.SSprite.MoveImg(float64(WinWidth)*4/5, buttonYPos)
	g.bHeal.SSprite.Display(screen)

	g.bInfo.SSprite.ResetGeoM()
	g.bInfo.SSprite.CenterImg()
	g.bInfo.SSprite.MoveImg(float64(WinWidth)/2, buttonYPos - 50)
	g.bInfo.SSprite.Display(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return WinWidth, WinHeight
}

func main() {
	// Loading assets
	view.InitAssets()

	// Init Game
	var g *Game = &Game{}
	g.Init(6)

	// Init Window
	ebiten.SetWindowSize(WinWidth, WinHeight)
	ebiten.SetWindowTitle("Camaretto")

	// Game Loop
	if err = ebiten.RunGame(g); err != nil {
		log.Fatal("[MAIN]", err)
	}
}