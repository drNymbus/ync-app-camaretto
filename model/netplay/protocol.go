package netplay

import (
	"camaretto/model/component"
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
	Players []*component.PlayerInfo
	Action *component.Action
	Reveal []bool
}
