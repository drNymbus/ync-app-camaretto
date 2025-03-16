package dialogue

import ()

type ActionDialogue struct {
	common []string
	uncommon []string
	rare []string
	epic []string
	legendary []string
}

type ConditionalActionDialogue struct {
	minus *ActionDialogue
	equal *ActionDialogue
	plus *ActionDialogue
}

type PersonaDialogue struct {
	name string
	intro *ActionDialogue

	victory *ActionDialogue
	defeat *ActionDialogue

	attack *ConditionalActionDialogue
	attackCharged *ConditionalActionDialogue

	charge *ActionDialogeue
	heal *ConditionalActionDialogue
}