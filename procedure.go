package donothing

import (
	"io"
)

// A Procedure is a sequence of Steps that can be executed or rendered to markdown.
type Procedure struct {
	// The procedure's short description, as provided with Short()
	short string
	// The steps in the procedure, in the order they were added
	steps []*Step
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
	step := NewStep()
	fn(step)
}

// Check validates that the procedure makes sense.
//
// It returns an error if anything's wrong.
func (pcd *Procedure) Check() error {
	return nil
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
	pcd.steps = make([]*Step, 0)
	return pcd
}
