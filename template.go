package donothing

import (
	"fmt"
	"regexp"
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
{{.SectionHeader}}{{if .ParentAnchor}}

[Up]({{.ParentAnchor}}){{end}}{{if .Body}}

{{.Body}}{{end -}}
{{if .InputDefs}}

{{template "inputs" .InputDefs}}{{end -}}
{{if .OutputDefs}}

{{template "outputs" .OutputDefs}}{{end -}}
{{if eq .Depth 0}}

{{template "table_of_contents" .Children}}{{end -}}
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

	// TemplateTableOfContents is the Markdown template with which we render the table of contents.
	//
	// It is rendered recursively, taking as . the parent step's Children slice.
	TemplateTableOfContents string = `{{define "table_of_contents" -}}
{{if . -}}
{{range .}}
{{.TOCIndent}}- [{{.Title}}]({{.Anchor}}){{template "table_of_contents" .Children}}{{end -}}
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
		TemplateTableOfContents,
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
	Parent     *StepTemplateData
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

// ParentAnchor returns an HTML anchor pointing to the parent section.
//
// If there is no parent section (because this StepTemplateData came from the root step), returns
// the empty string.
func (td StepTemplateData) ParentAnchor() string {
	if td.Parent == nil {
		return ""
	}
	return td.Parent.Anchor()
}

// Anchor returns an href for the step's Markdown section.
//
// For example, if step.SectionHeader() returns "### (2.1) Blah blah! Blah.", Anchor returns
// "#2-1-blah-blah-blah".
//
// According to the internet, this is (er, was in 2015) the code that GitHub uses to convert section
// headers to anchors:
// https://github.com/gjtorikian/html-pipeline/blob/main/lib/html/pipeline/toc_filter.rb
func (td StepTemplateData) Anchor() string {
	// Convert header to lowercase
	s0 := strings.ToLower(td.SectionHeader())
	// Remove header indicators (e.g. ###)
	s1 := strings.TrimLeft(s0, "#")
	// Remove initial space (the space that occurs after the header indicators)
	s2 := strings.TrimLeft(s1, " ")
	// Remove any characters that aren't allowed in an anchor
	s3 := strings.Join(
		strings.FieldsFunc(
			s2,
			func(char rune) bool {
				return ("" == regexp.MustCompile(`[[:alnum:]- ]`).FindString(string(char)))
			},
		),
		"",
	)
	// Replace spaces with hyphens
	s4 := regexp.MustCompile(`\s+`).ReplaceAllLiteralString(s3, "-")
	// Prepend #
	return fmt.Sprintf("#%s", s4)
}

// Returns the indent that should prefix the step's table of contents line.
func (td StepTemplateData) TOCIndent() string {
	return strings.Repeat("    ", td.Depth-1)
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
func NewStepTemplateData(step *Step, parent *StepTemplateData, recursive bool) StepTemplateData {
	td := StepTemplateData{
		Depth:      step.Depth(),
		Pos:        step.Pos(),
		Title:      step.GetShort(),
		Body:       step.GetLong(),
		InputDefs:  step.GetInputDefs(),
		OutputDefs: step.GetOutputDefs(),
		Parent:     parent,
		Children:   nil,
	}

	if recursive {
		td.Children = make([]StepTemplateData, 0)
		for _, c := range step.GetChildren() {
			td.Children = append(td.Children, NewStepTemplateData(c, &td, true))
		}
	}

	return td
}
