package model

//go:generate generator

// Task is a concrete piece of work that cannot
// be split any further.
//
// This type is very general and should be embedded in more
// concrete types, accordingly implementing logic to
// make this Task comparable to others with respect to it's
// SkillWeights. For example:
//
//	type CodeTask struct {
//		Task
//		Flags string
//		Languages []string
//	}
//
//	type QuizTask {
//		Task
//		Questions url.URL
//	}
type Task struct {
	Assignment

	// Says what skills are needed/exercised to complete
	// the Task.
	SkillWeights SkillWeights

	// Refers to some logic that looks at the Submissions
	// of this task and produces a set of skills that
	// represent how well the user did in doing this Task.
	// It is to be weighted by Skillweights.
	Tasker int
}
