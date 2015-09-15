package model

// Skills we assess.
//
// TODO(flowlo, victorbalan): Add further categories
// of assessment.
type Skills struct {
	Readability,
	Security,
	CodingSpeed,
	Algorithmics float64 `json:",omitempty"`
}

// SkillWeights can be used to express what impact or
// rating a Task has on a set of Skills.
type SkillWeights Skills

func (s Skills) Add(skills Skills) Skills {
	return Skills{
		Readability:  skills.Readability + s.Readability,
		Security:     skills.Security + s.Security,
		Algorithmics: skills.Algorithmics + s.Algorithmics,
		CodingSpeed:  skills.CodingSpeed + s.CodingSpeed,
	}
}

func (s Skills) Mul(skills Skills) Skills {
	return Skills{
		Readability:  skills.Readability * s.Readability,
		Security:     skills.Security * s.Security,
		Algorithmics: skills.Algorithmics * s.Algorithmics,
		CodingSpeed:  skills.CodingSpeed * s.CodingSpeed,
	}
}

func (s Skills) Div(skills Skills) Skills {
	return Skills{
		Readability:  skills.Readability / s.Readability,
		Security:     skills.Security / s.Security,
		Algorithmics: skills.Algorithmics / s.Algorithmics,
		CodingSpeed:  skills.CodingSpeed / s.CodingSpeed,
	}
}
