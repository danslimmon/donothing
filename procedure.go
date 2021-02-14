package donothing

import (
	"io"
)

// A Procedure is a sequence of Steps that can be executed or rendered to markdown.
type Procedure struct{}

// Short provides the procedure with a short description.
//
// The short description will be the title of the rendered markdown document when Render is called,
// so it should be concise and accurate.
func (pcd *Procedure) Short(s string) {}

// AddStep adds a step to the procedure.
func (pcd *Procedure) AddStep(func(*Step)) {}

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
	return nil
}
