package view

import (
	"camaretto/model"
)

func DrawPlayers(dst *ebiten.Image, players []*model.Player) {
	var nbPlayers int = len(players)
	var angleStep float64 = 2*math.Pi / float64(nbPlayers)
	var radius float64 = 200

	var centerX float64 = float64(WinWidth)/2
	var centerY float64 = (float64(WinHeight) * 6/8)/2

	for i, player := range players {
		var theta float64 = angleStep * float64(i)
		var x float64 = centerX + (radius * math.Cos(theta + math.Pi/2))
		var y float64 = centerY + (radius * math.Sin(theta + math.Pi/2))
		player.Render(dst, x, y, theta)
	}
}

func DrawDeck(dst *ebiten.Image, deck *model.Deck) {
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
}