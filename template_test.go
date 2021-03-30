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
				Depth:      0,
				Pos:        []int{},
				Title:      "root step",
				Body:       "",
				InputDefs:  []InputDef{},
				OutputDefs: []OutputDef{},
				Children:   []StepTemplateData{},
			},
			Out: `# root step`,
		},
		testCase{
			In: StepTemplateData{
				Depth:      1,
				Pos:        []int{3},
				Title:      "empty step",
				Body:       "",
				InputDefs:  []InputDef{},
				OutputDefs: []OutputDef{},
				Children:   []StepTemplateData{},
			},
			Out: `## (3) empty step`,
		},
		testCase{
			In: StepTemplateData{
				Depth:      1,
				Pos:        []int{3},
				Title:      "step with parent",
				Body:       "",
				InputDefs:  []InputDef{},
				OutputDefs: []OutputDef{},
				Parent: &StepTemplateData{
					Depth:      0,
					Pos:        []int{},
					Title:      "root step",
					Body:       "",
					InputDefs:  []InputDef{},
					OutputDefs: []OutputDef{},
					Parent:     nil,
					// If this were real, Children would contain "step with parent", but it doesn't
					// matter for the sake of this test.
					Children: []StepTemplateData{},
				},
				Children: []StepTemplateData{},
			},
			Out: `## (3) step with parent

[Up](#root-step)`,
		},
		testCase{
			In: StepTemplateData{
				Depth:      2,
				Pos:        []int{3, 1},
				Title:      "step with inputs",
				Body:       "",
				InputDefs:  []InputDef{InputDef{}},
				OutputDefs: []OutputDef{},
				Children:   []StepTemplateData{},
			},
			Out: `### (3.1) step with inputs

INPUTS`,
		},
		testCase{
			In: StepTemplateData{
				Depth:      2,
				Pos:        []int{4, 1},
				Title:      "step with outputs",
				Body:       "",
				InputDefs:  []InputDef{},
				OutputDefs: []OutputDef{OutputDef{}},
				Children:   []StepTemplateData{},
			},
			Out: `### (4.1) step with outputs

OUTPUTS`,
		},
		testCase{
			In: StepTemplateData{
				Depth:      2,
				Pos:        []int{5, 9},
				Title:      "step with both inputs and outputs",
				Body:       "",
				InputDefs:  []InputDef{InputDef{}},
				OutputDefs: []OutputDef{OutputDef{}},
				Children:   []StepTemplateData{},
			},
			Out: `### (5.9) step with both inputs and outputs

INPUTS

OUTPUTS`,
		},
		testCase{
			In: StepTemplateData{
				Depth:      2,
				Pos:        []int{9, 2},
				Title:      "step with body and outputs",
				Body:       "body of the step",
				InputDefs:  []InputDef{},
				OutputDefs: []OutputDef{OutputDef{}},
				Children:   []StepTemplateData{},
			},
			Out: `### (9.2) step with body and outputs

body of the step

OUTPUTS`,
		},
		testCase{
			In: StepTemplateData{
				Depth:      2,
				Pos:        []int{2, 6},
				Title:      "step with body and inputs and outputs",
				Body:       "body of the step",
				InputDefs:  []InputDef{InputDef{}},
				OutputDefs: []OutputDef{OutputDef{}},
				Children:   []StepTemplateData{},
			},
			Out: `### (2.6) step with body and inputs and outputs

body of the step

INPUTS

OUTPUTS`,
		},
		testCase{
			In: StepTemplateData{
				Depth:      1,
				Title:      "step with child",
				Pos:        []int{2},
				Body:       "",
				InputDefs:  []InputDef{},
				OutputDefs: []OutputDef{},
				Children: []StepTemplateData{
					StepTemplateData{
						Depth:      2,
						Pos:        []int{2, 6},
						Title:      "child step 0",
						Body:       "",
						InputDefs:  []InputDef{},
						OutputDefs: []OutputDef{},
						Children:   []StepTemplateData{},
					},
				},
			},
			Out: `## (2) step with child

### (2.6) child step 0`,
		},
		testCase{
			In: StepTemplateData{
				Depth:      1,
				Pos:        []int{6},
				Title:      "step with body, outputs, and children with bodies",
				Body:       "",
				InputDefs:  []InputDef{},
				OutputDefs: []OutputDef{OutputDef{}},
				Children: []StepTemplateData{
					StepTemplateData{
						Depth:      2,
						Pos:        []int{6, 0},
						Title:      "child step 0",
						Body:       "body of child 0",
						InputDefs:  []InputDef{},
						OutputDefs: []OutputDef{},
						Children:   []StepTemplateData{},
					},
					StepTemplateData{
						Depth:      2,
						Pos:        []int{6, 1},
						Title:      "child step 1",
						Body:       "body of child 1",
						InputDefs:  []InputDef{},
						OutputDefs: []OutputDef{},
						Children:   []StepTemplateData{},
					},
				},
			},
			Out: `## (6) step with body, outputs, and children with bodies

OUTPUTS

### (6.0) child step 0

body of child 0

### (6.1) child step 1

body of child 1`,
		},
		testCase{
			In: StepTemplateData{
				Depth:      1,
				Pos:        []int{5},
				Title:      "step with grandchildren with bodies",
				Body:       "",
				InputDefs:  []InputDef{},
				OutputDefs: []OutputDef{},
				Children: []StepTemplateData{
					StepTemplateData{
						Depth:      2,
						Pos:        []int{5, 3},
						Title:      "child step 0",
						Body:       "body of child 0",
						InputDefs:  []InputDef{},
						OutputDefs: []OutputDef{},
						Children: []StepTemplateData{
							StepTemplateData{
								Depth:      3,
								Pos:        []int{5, 3, 5},
								Title:      "grandchild step 0",
								Body:       "body of grandchild 0",
								InputDefs:  []InputDef{},
								OutputDefs: []OutputDef{},
								Children:   []StepTemplateData{},
							},
						},
					},
				},
			},
			Out: `## (5) step with grandchildren with bodies

### (5.3) child step 0

body of child 0

#### (5.3.5) grandchild step 0

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

// TemplateExecStep should render correctly with various inputs.
//
// Output from TemplateExecStep should never end with a newline. Spacing will be handled by the
// caller.
func TestTemplateExecStep(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	type testCase struct {
		In  StepTemplateData
		Out string
	}

	testCases := []testCase{
		testCase{
			In: StepTemplateData{
				Depth: 0,
				Title: "blah blah",
				Body:  "this is the description of my step",
			},
			Out: `# blah blah

this is the description of my step`,
		},
		testCase{
			In: StepTemplateData{
				Depth: 2,
				Pos:   []int{0, 3},
				Title: "step title",
				Body:  "this is the description of my step",
			},
			Out: `### (0.3) step title

this is the description of my step`,
		},
	}

	tpl, err := template.New("test").Parse(TemplateExecStep)
	assert.Nil(err)

	for i, tc := range testCases {
		t.Logf("test case %d", i)

		var b bytes.Buffer
		err = tpl.Execute(&b, tc.In)
		assert.Nil(err)
		assert.Equal(tc.Out, b.String())
	}
}

func TestStepTemplateData_SectionHeader(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	{
		templateData := StepTemplateData{
			Depth: 2,
			Pos:   []int{0, 2},
			Title: "Short description of step",
		}
		assert.Equal("### (0.2) Short description of step", templateData.SectionHeader())
	}
}

func TestStepTemplateData_Anchor(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	{
		templateData := StepTemplateData{
			Depth: 2,
			Pos:   []int{0, 2},
			Title: "Short! description of step",
		}
		assert.Equal("#02-short-description-of-step", templateData.Anchor())
	}

	{
		templateData := StepTemplateData{
			Depth: 0,
			Pos:   []int{},
			Title: "Description of root step!",
		}
		assert.Equal("#description-of-root-step", templateData.Anchor())
	}

	{
		templateData := StepTemplateData{
			Depth: 3,
			Pos:   []int{0, 1, 2},
			Title: "!",
		}
		assert.Equal("#012-", templateData.Anchor())
	}
}

// NewStepTemplateData with recursive=true should return a StepTemplateData with descendants.
func TestNewStepTemplateData_Recursive(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	// Step with no children. .Children should be an empty slice.
	{
		step := NewStep()
		step.Name("childlessStep")
		step.Short("fhgwhgads")

		templateData := NewStepTemplateData(step, nil, true)

		assert.Equal(0, templateData.Depth)
		assert.Equal("fhgwhgads", templateData.Title)
		assert.Equal([]int{}, templateData.Pos)
		assert.Exactly([]StepTemplateData{}, templateData.Children)
	}

	// Step with children but no grandchildren.
	{
		step := NewStep()
		step.Name("parentStep")
		step.Short("fhgwhgads")
		step.AddStep(func(step *Step) {
			step.Name("child0")
			step.Short("child 0")
		})
		step.AddStep(func(step *Step) {
			step.Name("child1")
			step.Short("child 1")
		})

		templateData := NewStepTemplateData(step, nil, true)

		assert.Equal(0, templateData.Depth)
		assert.Equal("fhgwhgads", templateData.Title)
		assert.Equal([]int{}, templateData.Pos)
		assert.Equal(2, len(templateData.Children))

		child0 := templateData.Children[0]
		assert.Equal(1, child0.Depth)
		assert.Equal("child 0", child0.Title)
		assert.Equal([]int{0}, child0.Pos)
		assert.Equal(0, len(child0.Children))

		child1 := templateData.Children[1]
		assert.Equal(1, child1.Depth)
		assert.Equal("child 1", child1.Title)
		assert.Equal([]int{1}, child1.Pos)
		assert.Equal(0, len(child1.Children))
	}

	// Step with a grandchild
	{
		step := NewStep()
		step.Name("grandparentStep")
		step.Short("fhgwhgads")
		step.AddStep(func(step *Step) {
			step.Name("child0")
			step.Short("child 0")
			step.AddStep(func(step *Step) {
				step.Name("grandchild0")
				step.Short("grandchild 0")
			})
		})

		templateData := NewStepTemplateData(step, nil, true)

		assert.Equal(0, templateData.Depth)
		assert.Equal("fhgwhgads", templateData.Title)
		assert.Equal([]int{}, templateData.Pos)
		assert.Equal(1, len(templateData.Children))
		assert.Equal(1, len(templateData.Children[0].Children))

		child := templateData.Children[0]
		assert.Equal(1, child.Depth)
		assert.Equal("child 0", child.Title)
		assert.Equal([]int{0}, child.Pos)
		assert.Equal(1, len(child.Children))

		grandchild := child.Children[0]
		assert.Equal(2, grandchild.Depth)
		assert.Equal("grandchild 0", grandchild.Title)
		assert.Equal([]int{0, 0}, grandchild.Pos)
		assert.Equal(0, len(grandchild.Children))
	}
}

// NewStepTemplateData with recursive=false should return a StepTemplateData with .Children = nil
func TestNewStepTemplateData_Nonrecursive(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	step := NewStep()
	step.Name("parentStep")
	step.Short("fhgwhgads")
	step.AddStep(func(step *Step) {
		step.Name("child0")
		step.Short("child 0")
	})

	templateData := NewStepTemplateData(step, nil, false)

	assert.Nil(templateData.Children)
}

func TestStep_numericPathToString(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	assert.Equal("", StepTemplateData{Pos: []int{}}.numericPathToString())
	assert.Equal("0", StepTemplateData{Pos: []int{0}}.numericPathToString())
	assert.Equal("3.1.4.1.5.9", StepTemplateData{Pos: []int{3, 1, 4, 1, 5, 9}}.numericPathToString())
}
