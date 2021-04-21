package donothing

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"text/template"
)

// AddTemplateDoc adds to the given template the overall Markdown doc template.
func AddTemplateDoc(tpl *template.Template) {
	txt := `{{template "step" .}}
`
	template.Must(tpl.Parse(txt))
}

// AddTemplateStep adds to the given template the Markdown template with which we render a Step.
//
// The input passed as . is an instance of StepTemplateData.
func AddTemplateStep(tpl *template.Template) {
	newTpl := tpl.New("step")
	txt := `{{define "step" -}}
{{.SectionHeader}}{{if .ParentAnchor}}

@@{{.StepName}}@@
•
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
	template.Must(newTpl.Parse(txt))
}

// AddTemplateExecStep adds to the given template the template that represents a Step in Execute()
func AddTemplateExecStep(tpl *template.Template) {
	txt := `{{.SectionHeader}}{{if .Body}}

{{.Body}}{{end -}}`
	template.Must(tpl.Parse(txt))
}

// AddTemplateInputs adds the step inputs template to the given template.
//
// This is the "**Inputs**" section of a step's documentation. It takes as . a slice of InputDef
// instances.
func AddTemplateInputs(tpl *template.Template) {
	newTpl := tpl.New("inputs")
	txt := `{{define "inputs" -}}
{{if . -}}
**Inputs**:
{{range .}}
  - @@{{.Name}}@@{{end -}}
{{else}}{{end -}}
{{end}}`
	template.Must(newTpl.Parse(txt))
}

// AddTemplateOutputs adds the step outputs template to the given template.
//
// This is the "**Outputs**" section of a step's documentation. It takes as . a slice of OutputDef
// instances.
func AddTemplateOutputs(tpl *template.Template) {
	newTpl := tpl.New("outputs")
	txt := `{{define "outputs" -}}
{{if . -}}
**Outputs**:
{{range .}}
  - @@{{.Name}}@@ ({{.ValueType}}): {{.Short}}{{end -}}
{{else -}}{{end -}}
{{end}}`
	template.Must(newTpl.Parse(txt))
}

// AddTemplateTableOfContents adds the table of contents template to the given template.
func AddTemplateTableOfContents(tpl *template.Template) {
	newTpl := tpl.New("table_of_contents")
	newTpl.Funcs(template.FuncMap{
		// We use this to tell whether we're in the last element of a pipeline in order to decide whether to
		// add a newline.
		"plus1": func(i int) int {
			return i + 1
		},
	})
	txt := `{{define "table_of_contents" -}}
{{if . -}}
{{- $n := len . -}}
{{range $i, $e := . -}}
{{.TOCIndent}}- [{{.Title}}]({{.Anchor}}){{if $e.Children}}
{{template "table_of_contents" .Children}}{{end}}{{if lt (plus1 $i) $n}}
{{end}}{{end -}}
{{end -}}
{{end}}`
	template.Must(newTpl.Parse(txt))
}

// DocTemplate returns the template for a Markdown document.
func DocTemplate() (*template.Template, error) {
	tpl := template.New("doc")

	AddTemplateDoc(tpl)
	AddTemplateStep(tpl)
	AddTemplateInputs(tpl)
	AddTemplateOutputs(tpl)
	AddTemplateTableOfContents(tpl)

	return tpl, nil
}

// ExecTemplate returns the template for output during procedure execution.
func ExecTemplate() (*template.Template, error) {
	tpl := template.New("exec")
	AddTemplateExecStep(tpl)
	return tpl, nil
}

// StepTemplateData is the thing that gets passed to a step template on evaluation.
type StepTemplateData struct {
	Depth      int
	Pos        []int
	StepName   string
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
		StepName:   step.AbsoluteName(),
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
