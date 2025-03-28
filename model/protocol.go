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
	var s string = "\n"
	s = s + "\tPLAYERS:\n"
	for _, info := range state.Players {
		s = s + "\t\t" + strconv.Itoa(info.Index) + "," + info.Name + "\n"
	}

	s = s + "\tSTATE:\n"
	s = s + "\t\tGame=" + state.Game.String() + "\n"
	s = s + "\t\tFocus=" + state.Focus.String() + "\n"
	s = s + "\t\tTurn=" + strconv.Itoa(state.Turn) + "\n"
	s = s + "\t\tPlayer=" + strconv.Itoa(state.Player) + "\n"
	s = s + "\t\tCard=" + strconv.Itoa(state.Card) + "\n"
	s = s + "\t\tReveal=["
	for i, v := range state.Reveal {
		s = s + "(" + strconv.Itoa(i) + ","
		if v {
			s = s + "true"
		} else {
			s = s + "false"
		}
		s = s + ")"
		if i != len(state.Reveal)-1 { s = s + ")," }
	}
	s = s + "]\n"

	return s
}

type MessageType int
const (
	PLAYERS MessageType = iota
	STATE
	START
)

func (m MessageType) String() string {
	var name []string = []string {"PLAYERS", "STATE", "START"}
	return name[int(m)]
}

type Message struct {
	Typ MessageType
	Players []*PlayerInfo
	Game *CamarettoState
}
