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

	// To normalize this task.
	//
	// TODO(flowlo): Clear specification.
	Logic logic
}
