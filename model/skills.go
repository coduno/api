package model

// Skills we assess.
//
// TODO(flowlo, victorbalan): Add further categories
// of assessment.
type Skills struct {
	Readability,
	Security,
	Algorithmics float64
}

// SkillWeights can be used to express what impact or
// rating a Task has on a set of Skills.
type SkillWeights Skills
