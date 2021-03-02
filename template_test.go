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
