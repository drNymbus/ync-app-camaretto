package netplay

import (
	"camaretto/model/game"
)

type MessageType int
const (
	PLAYERS MessageType = iota
	INIT
	ACTION
	START
)

type Message struct {
	Typ MessageType
	Seed int64
	Players []*game.PlayerInfo
	Action *game.Action
	Reveal []bool
}
