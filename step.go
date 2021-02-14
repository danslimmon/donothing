package donothing

// A Step is an individual action to be performed as part of a procedure.
//
// Steps must have a name (specified with Name()) and may have any number of substeps (provided with
// AddStep()).
//
// Steps may also have outputs and inputs, defined with Output* and Input*. An output defined by one
// step can be referenced as an input by any subsequent step in the procedure.
type Step struct{}

// Name gives the step a name, which must be unique among siblings.
//
// By convention, step names should be in camelCase.
//
// Step names are used by donothing to refer unambiguously to a step. Each step has an "absolute
// name", which is composed of a dot-separated sequence of the names of all its parent steps,
// followed by the step's own name. For example, the name "restoreBackup.loadData" refers to the
// "loadData" step, which is a child of the "restoreBackup" step.
func (step *Step) Name(s string) {}

// Short gives the step a short description.
//
// The short description will be the name of the step's corresponding section in the rendered
// markdown document.
func (step *Step) Short(s string) {}

// AddStep adds a child step to the Step.
//
// A new Step will be instantiated and passed to fn, which is responsible for defining the new child
// step.
func (step *Step) AddStep(fn func(*Step)) {}

// OutputString specifies a string output to be produced by the step.
//
// name is the output's name, which must be unique within the procedure. If any two outputs have the
// same name – even if the outputs are associated with steps with different parents – the procedure
// will fail to execute or render. Procedure.Check() will also return an error.
//
// desc should be a concise description of the output. This will be used to prompt the user for
// an output value if the Step is manual, and it will also be mentioned in the procedure's Markdown
// documentation.
func (step *Step) OutputString(name string, desc string) {}

// InputInt specifies an integer input taken by the step.
//
// If name matches the name of an integer output produced by a previous step, then the input's value
// will be automatically set to the value of that output. Otherwise, the user will be prompted for
// the input's value.
func (step *Step) InputInt(name string, required bool) {}

// NewStep returns a new step.
//
// Generally, donothing scripts shouldn't call NewStep directly. Instead, they should use
// *Procedure.AddStep or *Step.AddStep.
func NewStep() *Step {
	return nil
}
