package donothing

import (
	"bytes"
	"testing"
	"text/template"

	"github.com/stretchr/testify/assert"
)

// TemplateInputs should render correctly with various inputs.
//
// Output from TemplateInputs should never end with a newline. Spacing between sections will be
// handled by the template that calls it.
func TestTemplateInputs(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	type testCase struct {
		In  []InputDef
		Out string
	}

	testCases := []testCase{
		testCase{
			In:  []InputDef{},
			Out: ``,
		},
		testCase{
			In: []InputDef{
				InputDef{
					ValueType: "string",
					Name:      "foo",
					Required:  true,
				},
			},
			Out: `**Inputs**:

  - @@foo@@`,
		},
		testCase{
			In: []InputDef{
				InputDef{
					ValueType: "string",
					Name:      "foo",
					Required:  true,
				},
				InputDef{
					ValueType: "int",
					Name:      "bar",
					Required:  false,
				},
			},
			Out: `**Inputs**:

  - @@foo@@
  - @@bar@@`,
		},
	}

	tpl, err := template.New("test").Parse(`{{template "inputs" .}}`)
	assert.Nil(err)
	_, err = tpl.Parse(TemplateInputs)
	assert.Nil(err)

	for i, tc := range testCases {
		t.Logf("test case %d", i)

		var b bytes.Buffer
		err = tpl.Execute(&b, tc.In)
		assert.Nil(err)
		assert.Equal(tc.Out, b.String())
	}
}

// TemplateOutputs should render correctly with various inputs.
//
// Output from TemplateOutputs should never end with a newline. Spacing between sections will be
// handled by the template that calls it.
func TestTemplateOutputs(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	type testCase struct {
		In  []OutputDef
		Out string
	}

	testCases := []testCase{
		testCase{
			In:  []OutputDef{},
			Out: ``,
		},
		testCase{
			In: []OutputDef{
				OutputDef{
					ValueType: "string",
					Name:      "foo",
					Short:     "foo's short description",
				},
			},
			Out: `**Outputs**:

  - @@foo@@ (string): foo's short description`,
		},
		testCase{
			In: []OutputDef{
				OutputDef{
					ValueType: "string",
					Name:      "foo",
					Short:     "foo's short description",
				},
				OutputDef{
					ValueType: "int",
					Name:      "bar",
					Short:     "bar's short description",
				},
			},
			Out: `**Outputs**:

  - @@foo@@ (string): foo's short description
  - @@bar@@ (int): bar's short description`,
		},
	}

	tpl, err := template.New("test").Parse(`{{template "outputs" .}}`)
	assert.Nil(err)
	_, err = tpl.Parse(TemplateOutputs)
	assert.Nil(err)

	for i, tc := range testCases {
		t.Logf("test case %d", i)

		var b bytes.Buffer
		err = tpl.Execute(&b, tc.In)
		assert.Nil(err)
		assert.Equal(tc.Out, b.String())
	}
}

// TemplateStep should render correctly with various inputs.
//
// Output from TemplateStep should never end with a newline. Spacing between sections will be
// handled by the template that calls it.
func TestTemplateStep(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	type testCase struct {
		In  StepTemplateData
		Out string
	}

	testCases := []testCase{
		testCase{
			In: StepTemplateData{
				HeaderPrefix: "#",
				Title:        "empty step",
				Body:         "",
				InputDefs:    []InputDef{},
				OutputDefs:   []OutputDef{},
				Children:     []StepTemplateData{},
			},
			Out: `# empty step`,
		},
		testCase{
			In: StepTemplateData{
				HeaderPrefix: "##",
				Title:        "step with inputs",
				Body:         "",
				InputDefs:    []InputDef{InputDef{}},
				OutputDefs:   []OutputDef{},
				Children:     []StepTemplateData{},
			},
			Out: `## step with inputs

INPUTS`,
		},
		testCase{
			In: StepTemplateData{
				HeaderPrefix: "##",
				Title:        "step with outputs",
				Body:         "",
				InputDefs:    []InputDef{},
				OutputDefs:   []OutputDef{OutputDef{}},
				Children:     []StepTemplateData{},
			},
			Out: `## step with outputs

OUTPUTS`,
		},
		testCase{
			In: StepTemplateData{
				HeaderPrefix: "##",
				Title:        "step with both inputs and outputs",
				Body:         "",
				InputDefs:    []InputDef{InputDef{}},
				OutputDefs:   []OutputDef{OutputDef{}},
				Children:     []StepTemplateData{},
			},
			Out: `## step with both inputs and outputs

INPUTS

OUTPUTS`,
		},
		testCase{
			In: StepTemplateData{
				HeaderPrefix: "##",
				Title:        "step with body and outputs",
				Body:         "body of the step",
				InputDefs:    []InputDef{},
				OutputDefs:   []OutputDef{OutputDef{}},
				Children:     []StepTemplateData{},
			},
			Out: `## step with body and outputs

body of the step

OUTPUTS`,
		},
		testCase{
			In: StepTemplateData{
				HeaderPrefix: "##",
				Title:        "step with body and inputs and outputs",
				Body:         "body of the step",
				InputDefs:    []InputDef{InputDef{}},
				OutputDefs:   []OutputDef{OutputDef{}},
				Children:     []StepTemplateData{},
			},
			Out: `## step with body and inputs and outputs

body of the step

INPUTS

OUTPUTS`,
		},
		testCase{
			In: StepTemplateData{
				HeaderPrefix: "#",
				Title:        "step with child",
				Body:         "",
				InputDefs:    []InputDef{},
				OutputDefs:   []OutputDef{},
				Children: []StepTemplateData{
					StepTemplateData{
						HeaderPrefix: "##",
						Title:        "child step 0",
						Body:         "",
						InputDefs:    []InputDef{},
						OutputDefs:   []OutputDef{},
						Children:     []StepTemplateData{},
					},
				},
			},
			Out: `# step with child

## child step 0`,
		},
		testCase{
			In: StepTemplateData{
				HeaderPrefix: "#",
				Title:        "step with body, outputs, and children with bodies",
				Body:         "",
				InputDefs:    []InputDef{},
				OutputDefs:   []OutputDef{OutputDef{}},
				Children: []StepTemplateData{
					StepTemplateData{
						HeaderPrefix: "##",
						Title:        "child step 0",
						Body:         "body of child 0",
						InputDefs:    []InputDef{},
						OutputDefs:   []OutputDef{},
						Children:     []StepTemplateData{},
					},
					StepTemplateData{
						HeaderPrefix: "##",
						Title:        "child step 1",
						Body:         "body of child 1",
						InputDefs:    []InputDef{},
						OutputDefs:   []OutputDef{},
						Children:     []StepTemplateData{},
					},
				},
			},
			Out: `# step with body, outputs, and children with bodies

OUTPUTS

## child step 0

body of child 0

## child step 1

body of child 1`,
		},
		testCase{
			In: StepTemplateData{
				HeaderPrefix: "#",
				Title:        "step with grandchildren with bodies",
				Body:         "",
				InputDefs:    []InputDef{},
				OutputDefs:   []OutputDef{},
				Children: []StepTemplateData{
					StepTemplateData{
						HeaderPrefix: "##",
						Title:        "child step 0",
						Body:         "body of child 0",
						InputDefs:    []InputDef{},
						OutputDefs:   []OutputDef{},
						Children: []StepTemplateData{
							StepTemplateData{
								HeaderPrefix: "###",
								Title:        "grandchild step 0",
								Body:         "body of grandchild 0",
								InputDefs:    []InputDef{},
								OutputDefs:   []OutputDef{},
								Children:     []StepTemplateData{},
							},
						},
					},
				},
			},
			Out: `# step with grandchildren with bodies

## child step 0

body of child 0

### grandchild step 0

body of grandchild 0`,
		},
	}

	tpl, err := template.New("test").Parse(`{{template "step" .}}`)
	assert.Nil(err)
	template.Must(tpl.Parse(`{{define "inputs"}}INPUTS{{end}}`))
	template.Must(tpl.Parse(`{{define "outputs"}}OUTPUTS{{end}}`))
	_, err = tpl.Parse(TemplateStep)
	assert.Nil(err)

	for i, tc := range testCases {
		t.Logf("test case %d", i)

		var b bytes.Buffer
		err = tpl.Execute(&b, tc.In)
		assert.Nil(err)
		assert.Equal(tc.Out, b.String())
	}
}