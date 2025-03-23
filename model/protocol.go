package model

import (
	"camaretto/model/game"
)

type PlayerInfo struct {
	Index int
	Name string
}

type CamarettoState struct {
	Game game.GameState // Camaretto.state
	Focus game.FocusState // Camaretto.focus
	Turn int // Camaretto.playerTurn
	Player int // Camaretto.playerFocus
	Card int // Camaretto.cardFocus
	Reveal []bool // Camaretto.toReveal
}

type CamarettoInit struct {
	Seed int64
	NbPlayers int
	Names []string
}

type MessageType int
const (
	HANDSHAKE MessageType = 0
	STATE MessageType = 1
	INIT MessageType = 2
)

type Message struct {
	Typ MessageType
	Info *PlayerInfo
	State *CamarettoState
	Init *CamarettoInit
}
