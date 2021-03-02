package donothing

import (
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
{{.HeaderPrefix}} {{.Title}}{{if .Body}}

{{.Body}}{{end -}}
{{if .InputDefs}}

{{template "inputs" .InputDefs}}{{end -}}
{{if .OutputDefs}}

{{template "outputs" .OutputDefs}}{{end -}}
{{range .Children}}

{{template "step" .}}{{end -}}
{{end}}`

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
