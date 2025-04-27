package game

const (
	MaxNbPlayers int = 6
)

type GameState int
const (
	SET GameState = iota
	ATTACK
	SHIELD
	CHARGE
	HEAL
)

type FocusState int
const (
	NONE FocusState = iota
	PLAYER
	CARD
	REVEAL
	COMPLETE
)

type Action struct {
	State GameState
	Focus FocusState

	PlayerTurn int
	PlayerFocus int
	CardFocus int

	Reveal []bool
}

func NewAction(i int) *Action {
	var a *Action = &Action{}
	
	a.State = SET
	a.Focus = NONE

	a.PlayerTurn = i
	a.PlayerFocus = -1
	a.CardFocus = -1

	a.Reveal = []bool{}

	return a
}
