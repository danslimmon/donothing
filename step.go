package donothing

import (
	"strings"
)

// A Step is an individual action to be performed as part of a procedure.
//
// Steps must have a name (specified with Name()) and may have any number of substeps (provided with
// AddStep()).
//
// Steps may also have outputs and inputs, defined with Output* and Input*. An output defined by one
// step can be referenced as an input by any subsequent step in the procedure.
type Step struct {
	// The Step's name, as set by Name()
	name string
	// The Step's short description, as set by Short()
	short string
	// The Step's long description, as set by Long()
	long string

	// The Step's inputs and outputs, if any
	inputs  []InputDef
	outputs []OutputDef

	// The Step of which this Step is a child. nil if this is the root step.
	parent *Step
	// The Step's substeps, if any
	children []*Step
}

// Name gives the step a name, which must be unique among siblings.
//
// By convention, step names should be in camelCase.
//
// Step names are used by donothing to refer unambiguously to a step. Each step has an "absolute
// name", which is composed of a dot-separated sequence of the names of all its parent steps,
// followed by the step's own name. For example, the name "restoreBackup.loadData" refers to the
// "loadData" step, which is a child of the "restoreBackup" step.
func (step *Step) Name(s string) {
	step.name = s
}

// AbsoluteName returns the step's unique name.
func (step *Step) AbsoluteName() string {
	if step.parent == nil {
		return step.name
	}
	return strings.Join([]string{
		step.parent.AbsoluteName(),
		step.name,
	}, ".")
}

// Depth returns the step's depth in the tree.
//
// The root node's depth is 0, the root node's children are at depth 1, and so on.
func (step *Step) Depth() int {
	if step.parent == nil {
		return 0
	}
	return step.parent.Depth() + 1
}

// Short gives the step a short description.
//
// The short description will be the name of the step's corresponding section in the rendered
// markdown document.
func (step *Step) Short(s string) {
	step.short = s
}

// Long gives the step a long description.
//
// The long description will be the body of the step's corresponding section in the rendered
// markdown document.
//
// Before a step is rendered, any occurrences of the "backtick standin sequence" in the long
// description will be replaced with backtick characters. By default, the backtick standin sequence
// is "@@". This sequence can be reassigned using Procedure.BacktickStandin().
func (step *Step) Long(s string) {
	step.long = s
}

// AddStep adds a child step to the Step.
//
// A new Step will be instantiated and passed to fn, which is responsible for defining the new child
// step.
func (step *Step) AddStep(fn func(*Step)) {
	newStep := NewStep()
	newStep.parent = step
	fn(newStep)
	step.children = append(step.children, newStep)
}

// OutputString specifies a string output to be produced by the step.
//
// name is the output's name, which must be unique within the procedure. If any two outputs have the
// same name – even if the outputs are associated with steps with different parents – the procedure
// will fail to execute or render. Procedure.Check() will also return an error.
//
// desc should be a concise description of the output. This will be used to prompt the user for
// an output value if the Step is manual, and it will also be mentioned in the procedure's Markdown
// documentation.
func (step *Step) OutputString(name string, desc string) {
	output := NewOutputDef("string", name, desc)
	step.outputs = append(step.outputs, output)
}

// InputString specifies a string input taken by the step.
//
// If name matches the name of a string output produced by a previous step, then the input's value
// will be automatically set to the value of that output. Otherwise, the user will be prompted for
// the input's value.
func (step *Step) InputString(name, short string, required bool) {
	input := NewInputDef("string", name, short, required)
	step.inputs = append(step.inputs, input)
}

// Walk visits every step in the tree, calling fn on each.
//
// It's a depth-first walk, starting with step itself, then proceeding in sequence through the
// children of step and their children, recursively. This is the order in which the steps execute
// when Procedure.Execute() is called, as well as the order in which the steps are rendered into
// documentation.
//
// If fn returns an error for any step, Walk immediately exits, returning that error.
func (step *Step) Walk(fn func(*Step) error) error {
	if err := fn(step); err != nil {
		return err
	}
	for _, childStep := range step.children {
		if err := childStep.Walk(fn); err != nil {
			return err
		}
	}
	return nil
}

// NewStep returns a new step.
//
// Generally, donothing scripts shouldn't call NewStep directly. Instead, they should use
// *Procedure.AddStep or *Step.AddStep.
func NewStep() *Step {
	step := new(Step)
	step.children = make([]*Step, 0)
	step.inputs = make([]InputDef, 0)
	step.outputs = make([]OutputDef, 0)
	return step
}
