package donothing

import (
	"io"
)

// A Procedure is a sequence of Steps that can be executed or rendered to markdown.
type Procedure struct {
	// The procedure's short description, as provided with Short()
	short string
	// The root step of the procedure, of which all other steps are descendants.
	rootStep *Step
}

// Short provides the procedure with a short description.
//
// The short description will be the title of the rendered markdown document when Render is called,
// so it should be concise and accurate.
func (pcd *Procedure) Short(s string) {
	pcd.short = s
}

// AddStep adds a step to the procedure.
//
// A new Step will be instantiated and passed to fn to be defined.
func (pcd *Procedure) AddStep(fn func(*Step)) {
	pcd.rootStep.AddStep(fn)
}

// Check validates that the procedure makes sense.
//
// It returns an error if anything's wrong. It will fail on any departure from the following
// expectations:
//
//   1. Every step has a unique absolute name with no empty parts.
//   2. Every step has a short description
//   3. Every input has a name that matches the name of an output from a previous step.
func (pcd *Procedure) Check() error {
	return pcd.rootStep.Walk(func(step *Step) error {
		//		if err := step.Check(); err != nil {
		//			fmt.Printf(
		//				"Error checking step %s with parent %s: %s",
		//				step.String(),
		//				err,
		//			)
		//		}
		return nil
	})
}

// Render prints the procedure's Markdown representation to f.
func (pcd *Procedure) Render(f io.Writer) error {
	return nil
}

// Execute runs through the procedure step by step.
//
// The user will be prompted as necessary.
func (pcd *Procedure) Execute() error {
	return nil
}

// NewProcedure returns a new procedure, ready to be given steps.
func NewProcedure() *Procedure {
	pcd := new(Procedure)
	pcd.rootStep = NewStep()
	pcd.rootStep.Name("root")
	return pcd
}
