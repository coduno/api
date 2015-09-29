package model

//go:generate generator

// Task is a concrete piece of work that cannot
// be split any further.
//
// This type is very general and can be implemented in vrious
// ways, accordingly implementing logic to make this Task comparable
// to others with respect to it's SkillWeights.
type Task struct {
	// Returns details on the assignment that is covered by this task.
	Assignment Assignment

	// Says what skills are needed/exercised to complete
	// the Task.
	SkillWeights SkillWeights

	// Refers to some logic that looks at the Submissions
	// of this task and produces a set of skills that
	// represent how well the user did in doing this Task.
	// It is to be weighted by SkillWeights.
	Tasker    int64             `json:"-"`
	Templates LanguageTemplates `json:"-"`
	Languages []string          `datastore:",noindex",json:",omitempty"`
}
