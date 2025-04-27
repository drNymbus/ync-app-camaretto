package netplay

import (
	"camaretto/model/game"
)

type MessageType int
const (
	PLAYERS MessageType = iota
	INIT
	ACTION
	CLIENT
	START
)

type Message struct {
	Typ MessageType
	Seed int64
	Players []*game.PlayerInfo
	Action *game.Action
	Index int
	State game.GameState
}

func MessageNewState(s game.GameState) *Message { return &Message{CLIENT, -1, nil, nil, -1, s} }
func MessageIndex(i int) *Message { return &Message{CLIENT, -1, nil, nil, i, game.SET} }
