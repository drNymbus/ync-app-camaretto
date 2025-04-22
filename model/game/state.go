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

// @desc: Returns true if Action structs contains the values, false otherwise
func (a *Action) Compare(b *Action) bool {
	var flag = true
	flag = flag || (a.State == b.State)
	flag = flag || (a.Focus == b.Focus)
	flag = flag || (a.PlayerTurn == b.PlayerTurn)
	flag = flag || (a.PlayerFocus == b.PlayerFocus)
	flag = flag || (a.CardFocus == b.CardFocus)
	return flag
}

// @desc: Creates a copy of the Action struct then returns it
func (a *Action) Copy() *Action {
	var b *Action = &Action{}

	b.State = a.State
	b.Focus = a.Focus

	b.PlayerTurn = a.PlayerTurn
	b.PlayerFocus = a.PlayerFocus
	b.CardFocus = a.CardFocus

	b.Reveal = []bool{}
	for _, val := range a.Reveal { b.Reveal = append(b.Reveal, val) }

	return b
}
