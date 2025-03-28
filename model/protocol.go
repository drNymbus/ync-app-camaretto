package model

import (
	"time"

	"strconv"

	"camaretto/model/game"
)

type PlayerInfo struct {
	Index int
	Name string
}

type CamarettoState struct {
	Seed int64
	Players []*PlayerInfo
	Game game.GameState // Camaretto.state
	Focus game.FocusState // Camaretto.focus
	Turn int // Camaretto.playerTurn
	Player int // Camaretto.playerFocus
	Card int // Camaretto.cardFocus
	Reveal []bool // Camaretto.toReveal
}

func NewCamarettoState() *CamarettoState {
	var state *CamarettoState = &CamarettoState{}

	state.Seed = time.Now().UnixNano()
	state.Players = []*PlayerInfo{}
	
	state.Game = game.SET
	state.Focus = game.NONE

	state.Turn = -1
	state.Player = -1
	state.Card = -1

	state.Reveal = []bool{}

	return state
}

func (state *CamarettoState) toString() string {
	var s string = "STATE:\n"
	s = s + "\tPLAYERS:\n"

	for _, info := range state.Players {
		s = s + "\t\t" + strconv.Itoa(info.Index) + "," + info.Name + "\n"
	}

	return s
}

type MessageType int
const (
	PLAYERS MessageType = iota
	STATE
	START
)

type Message struct {
	Typ MessageType
	Players []*PlayerInfo
	Game *CamarettoState
}
