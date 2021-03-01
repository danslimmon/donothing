package donothing

import (
	"errors"
	"fmt"
	"io"
	"strings"
	"text/template"
)

const (
	MarkdownTemplate = `{{define "step" -}}
{{.HeaderPrefix}} {{.Title}}

{{if .Body}}{{.Body}}

{{end -}}

{{template "step_inputs" .InputDefs -}}
{{template "step_outputs" .OutputDefs -}}

{{range .Children}}{{template "step" .}}{{end -}}
{{- /* End {{define}} block */ -}}
{{end}}

{{define "step_inputs" -}}
{{if . -}}
**Inputs**:
{{range .}}    - @@{{.Name}}@@
{{end}}
{{end -}}

{{/* End {{define}} block */ -}}
{{end}}

{{define "step_outputs" -}}
{{if . -}}
**Outputs**:
{{range .}}    - @@{{.Name}}@@ ({{.ValueType}}): {{.Short}}
{{end}}
{{end -}}

{{/* End {{define}} block */ -}}
{{end}}`
)

// StepTemplateData is the thing that gets passed to a step template on evaluation.
type StepTemplateData struct {
	HeaderPrefix string
	Title        string
	Body         string
	InputDefs    []InputDef
	OutputDefs   []OutputDef
	Children     []StepTemplateData
}

// NewStepTemplateData returns a StepTemplateData instance for the given Step.
//
// It is called recursively on children of the Step in order to populate the StepTemplateData's
// Children attribute.
func NewStepTemplateData(step *Step) StepTemplateData {
	td := StepTemplateData{
		HeaderPrefix: strings.Repeat("#", step.Depth()+1),
		Title:        step.GetShort(),
		Body:         step.GetLong(),
		InputDefs:    step.GetInputDefs(),
		OutputDefs:   step.GetOutputDefs(),
		Children:     []StepTemplateData{},
	}

	for _, c := range step.GetChildren() {
		td.Children = append(td.Children, NewStepTemplateData(c))
	}

	return td
}

// A Procedure is a sequence of Steps that can be executed or rendered to markdown.
type Procedure struct {
	// The root step of the procedure, of which all other steps are descendants.
	rootStep *Step
}

// Short provides the procedure with a short description.
//
// The short description will be the title of the rendered markdown document when Render is called,
// so it should be concise and accurate.
func (pcd *Procedure) Short(s string) {
	pcd.rootStep.Short(s)
}

// Long provides the procedure with a long description.
//
// The long description will be shown to the user when they first execute the procedure. It will
// also be included in the opening section of the procedure's Markdown documentation.
//
// It should give an overview of the procedure's purpose and any important assumptions the procedure
// makes about the state of the world at the beginning of execution.
func (pcd *Procedure) Long(s string) {
	pcd.rootStep.Long(s)
}

// AddStep adds a step to the procedure.
//
// A new Step will be instantiated and passed to fn to be defined.
func (pcd *Procedure) AddStep(fn func(*Step)) {
	pcd.rootStep.AddStep(fn)
}

// Check validates that the procedure makes sense.
//
// If problems are found, it returns the list of problems along with an error.
//
// It checks the procedure against the following expectations:
//
//   1. Every step has a unique absolute name with no empty parts.
//   2. Every step has a short description
//   3. Every input has a name that matches the name of an output from a previous step.
func (pcd *Procedure) Check() ([]string, error) {
	steps := make(map[string]*Step)
	outputs := make(map[string]OutputDef)
	problems := make([]string, 0)

	err := pcd.rootStep.Walk(func(step *Step) error {
		absName := step.AbsoluteName()
		if absName[len(absName)-1:] == "." {
			if step.parent == nil {
				// I really hope this never happens. The root step should get its name from
				// donothing, not from the calling code.
				problems = append(problems, "Root step does not have name")
			} else {
				problems = append(problems, fmt.Sprintf("Child step of '%s' does not have name", step.parent.AbsoluteName()))
			}
		}

		if steps[absName] != nil {
			problems = append(problems, fmt.Sprintf("More than one step with name '%s'", absName))
		}
		steps[step.AbsoluteName()] = step

		if step.GetShort() == "" {
			problems = append(problems, fmt.Sprintf("Step '%s' has no Short value", absName))
		}

		for _, inputDef := range step.GetInputDefs() {
			matchingOutputDef, ok := outputs[inputDef.Name]
			if !ok {
				problems = append(problems, fmt.Sprintf(
					"Input '%s' of step '%s' does not refer to an output from any previous step",
					inputDef.Name,
					absName,
				))
				continue
			}
			if matchingOutputDef.ValueType != inputDef.ValueType {
				problems = append(problems, fmt.Sprintf(
					"Input '%s' of step '%s' has type '%s', but output '%s' has type '%s'",
					inputDef.Name,
					absName,
					inputDef.ValueType,
					matchingOutputDef.Name,
					matchingOutputDef.ValueType,
				))
			}
		}

		for _, outputDef := range step.GetOutputDefs() {
			outputs[outputDef.Name] = outputDef
		}

		return nil
	})
	if err != nil {
		return []string{}, fmt.Errorf("Error while checking procedure: %w", err)
	}

	if len(problems) > 0 {
		return problems, errors.New("Problems were found in the procedure")
	}
	return []string{}, nil
}

// Render prints the procedure's Markdown representation to f.
//
// Any occurrence of the string "@@" in the executed template output will be replaced with a
// backtick.
func (pcd *Procedure) Render(f io.Writer) error {
	if _, err := pcd.Check(); err != nil {
		return err
	}

	tpl, err := template.New("step").Parse(MarkdownTemplate)
	if err != nil {
		return err
	}

	tplData := NewStepTemplateData(pcd.rootStep)

	var b strings.Builder
	err = tpl.Execute(&b, tplData)
	if err != nil {
		return err
	}

	fmt.Fprintf(f, "%s", strings.Replace(b.String(), "@@", "`", -1))
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
