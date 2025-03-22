package model

import (
	"camaretto/model/game"
)

type PlayerInfo struct {
	index int
	name string
}

type CamarettoState struct {
	game game.GameState // Camaretto.state
	focus game.FocusState // Camaretto.focus
	turn int // Camaretto.playerTurn
	player int // Camaretto.playerFocus
	card int // Camaretto.cardFocus
	reveal []bool // Camaretto.toReveal
}

type CamarettoInit struct {
	seed int64
	nbPlayers int
	names []string
}

type MessageType int
const (
	HANDSHAKE MessageType = 0
	STATE MessageType = 1
	INIT MessageType = 2
)

type Message struct {
	type MessageType
	info PlayerInfo
	state CamarettoState
	init CamarettoInit
}
