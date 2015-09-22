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
	computed := Skills{}
	if s.Readability > 0 {
		computed.Readability = s.Readability / skills.Readability
	}
	if s.Security > 0 {
		computed.Security = s.Security / skills.Security
	}
	if s.Algorithmics > 0 {
		computed.Algorithmics = s.Algorithmics / skills.Algorithmics
	}
	if s.CodingSpeed > 0 {
		computed.CodingSpeed = s.CodingSpeed / skills.CodingSpeed
	}
	return computed
}

func (s Skills) DivBy(v float64) Skills {
	if v == 0 {
		v = 1
	}
	return Skills{
		Readability:  s.Readability / v,
		Security:     s.Security / v,
		Algorithmics: s.Algorithmics / v,
		CodingSpeed:  s.CodingSpeed / v,
	}
}
