package donothing

import (
	"fmt"
	"strconv"
	"strings"
	"text/template"
)

var (
	// TemplateDoc is the Markdown template with which we render an entire document.
	//
	// It takes as . an instance of StepTemplateData.
	TemplateDoc string = `{{template "step" .}}
`

	// TemplateStep is the Markdown template with which we render a Step.
	//
	// The input passed as . is an instance of StepTemplateData.
	TemplateStep string = `{{define "step" -}}
{{.SectionHeader}}{{if .Body}}

{{.Body}}{{end -}}
{{if .InputDefs}}

{{template "inputs" .InputDefs}}{{end -}}
{{if .OutputDefs}}

{{template "outputs" .OutputDefs}}{{end -}}
{{range .Children}}

{{template "step" .}}{{end -}}
{{end}}`

	// TemplateExecStep is the template we use to render a step when executing a procedure.
	TemplateExecStep string = `{{.SectionHeader}}{{if .Body}}

{{.Body}}{{end -}}`

	// TemplateInputs is the Markdown template with which we render a Step's InputDefs.
	//
	// It's the "**Inputs**" section of a step's documentation. It takes as . a slice of InputDef
	// instances.
	TemplateInputs string = `{{define "inputs" -}}
{{if . -}}
**Inputs**:
{{range .}}
  - @@{{.Name}}@@{{end -}}
{{else}}{{end -}}
{{end}}`

	// TemplateOutputs is the Markdown template with which we render a Step's OutputDefs.
	//
	// It's the "**Outputs**" section of a step's documentation. It takes as . a slice of OutputDef
	// instances.
	TemplateOutputs string = `{{define "outputs" -}}
{{if . -}}
**Outputs**:
{{range .}}
  - @@{{.Name}}@@ ({{.ValueType}}): {{.Short}}{{end -}}
{{else -}}{{end -}}
{{end}}`
)

// DocTemplate returns the template for a Markdown document.
func DocTemplate() (*template.Template, error) {
	tpl := template.New("doc")

	var err error
	for _, tplDef := range []string{
		TemplateDoc,
		TemplateStep,
		TemplateInputs,
		TemplateOutputs,
	} {
		_, err = tpl.Parse(tplDef)
		if err != nil {
			return nil, err
		}
	}

	return tpl, nil
}

// ExecTemplate returns the template for output during procedure execution.
func ExecTemplate() (*template.Template, error) {
	tpl := template.New("exec")
	_, err := tpl.Parse(TemplateExecStep)
	if err != nil {
		return nil, err
	}
	return tpl, nil
}

// StepTemplateData is the thing that gets passed to a step template on evaluation.
type StepTemplateData struct {
	Depth      int
	Pos        []int
	Title      string
	Body       string
	InputDefs  []InputDef
	OutputDefs []OutputDef
	Children   []StepTemplateData
}

// SectionHeader returns the header line for the step's section.
//
// For example, "## (0.2) Short description of step"
func (td StepTemplateData) SectionHeader() string {
	// Header prefix; e.g. "###"
	parts := []string{strings.Repeat("#", td.Depth+1)}

	// Numeric path part; e.g. "(0.2.1)". Absent if root step.
	if td.Depth > 0 {
		parts = append(parts, fmt.Sprintf("(%s)", td.numericPathToString()))
	}

	// Title part (the step's Short description)
	parts = append(parts, td.Title)

	return strings.Join(parts, " ")
}

// Anchor returns an href for the step's Markdown section.
//
// For example, if step.SectionHeader() returns "### (2.1) Blah blah! Blah.", Anchor returns
// "#2-1-blah-blah-blah".
func (td StepTemplateData) Anchor() string {
	return ""
}

// numericPathToString renders td.Pos to a dot-separated string.
//
// If td.Pos is empty, numericPathToString returns the empty string.
func (td StepTemplateData) numericPathToString() string {
	sPos := make([]string, len(td.Pos))
	for i := range td.Pos {
		sPos[i] = strconv.Itoa(td.Pos[i])
	}
	return strings.Join(sPos, ".")
}

// NewStepTemplateData returns a StepTemplateData instance for the given Step.
//
// If recursive is true, NewStepTemplateData is called recursively on children of the Step in order
// to populate the StepTemplateData's Children attribute. If recursive is false, the returned
// StepTemplateData struct will have Children == nil.
func NewStepTemplateData(step *Step, recursive bool) StepTemplateData {
	td := StepTemplateData{
		Depth:      step.Depth(),
		Pos:        step.Pos(),
		Title:      step.GetShort(),
		Body:       step.GetLong(),
		InputDefs:  step.GetInputDefs(),
		OutputDefs: step.GetOutputDefs(),
		Children:   nil,
	}

	if recursive {
		td.Children = make([]StepTemplateData, 0)
		for _, c := range step.GetChildren() {
			td.Children = append(td.Children, NewStepTemplateData(c, true))
		}
	}

	return td
}
