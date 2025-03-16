package dialogue

func LoadIntroPresidentDialogue() *ActionDialogue {
	var intro *ActionDialogue{}
	intro.common = []string {"", ""}
	intro.uncommon = []string {"", ""}
	intro.rare = []string {"", ""}
	intro.epic = []string {"", ""}
	intro.legendary = []string {"", ""}
	return intro
}

func LoadAttackPresidentDialogue() *ConditionalActionDialogue {
	var atk *ConditionalActionDialogue = &ConditionalActionDialogue{}

	atk.minus = &ActionDialogue{}
	atk.minus.common = []string {"", ""}
	atk.minus.uncommon = []string {"", ""}
	atk.minus.rare = []string {"", ""}
	atk.minus.epic = []string {"", ""}
	atk.minus.legendary = []string {"", ""}

	atk.equal = &ActionDialogue{}
	atk.equal.common = []string {"", ""}
	atk.equal.uncommon = []string {"", ""}
	atk.equal.rare = []string {"", ""}
	atk.equal.epic = []string {"", ""}
	atk.equal.legendary = []string {"", ""}

	atk.plus = &ActionDialogue{}
	atk.plus.common = []string {"", ""}
	atk.plus.uncommon = []string {"", ""}
	atk.plus.rare = []string {"", ""}
	atk.plus.epic = []string {"", ""}
	atk.plus.legendary = []string {"", ""}

	return atk
}

func LoadPresidentDialogue() *PersonaDialogue {
	var dialogue *PersonaDialogue = &PersonaDialogue{}
	dialogue.name = "PRESIDENT"

	dialogue.intro = LoadIntroDialogue()
	dialogue.attack = LoadAttackPresidentDialogue()

	return dialogue
}