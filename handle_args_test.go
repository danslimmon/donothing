package donothing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// DefaultCLI.Usage should return the usage string.
func TestDefaultCLI_Usage(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	pcd := NewProcedure()
	pcd.Short("Procedure's short description")

	type testCase struct {
		DefaultStep string
		Exp         string
	}

	testCases := []testCase{
		// With default step specified
		testCase{
			DefaultStep: "root/stepName",
			Exp: `USAGE: foo [options] [STEP_NAME]

Procedure's short description

OPTIONS: 
    --markdown    Instead of executing the procedure, print its Markdown documentation to stdout
    --help        Print usage message`,
		},
		// Without default step
		testCase{
			DefaultStep: "",
			Exp: `USAGE: foo [options] STEP_NAME

Procedure's short description

OPTIONS: 
    --markdown    Instead of executing the procedure, print its Markdown documentation to stdout
    --help        Print usage message`,
		},
	}

	for i, tc := range testCases {
		t.Logf("test case %d", i)
		cli, err := NewDefaultCLI("foo", pcd, tc.DefaultStep)
		assert.Nil(err)
		assert.Equal(tc.Exp, cli.Usage())
	}
}
