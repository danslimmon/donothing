package donothing

import (
	"bytes"
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

// DefaultCLI should print usage when --help is passed
func TestDefaultCLI_Help(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	pcd := NewProcedure()
	pcd.Short("Procedure's short description")

	type testCase struct {
		// os.Args
		Args []string
		// The CLI's default step
		DefaultStep string
		// Whether an error is expected
		ErrorExp bool
	}

	testCases := []testCase{
		testCase{
			Args:     []string{"foo", "--help"},
			ErrorExp: false,
		},
		testCase{
			Args:     []string{"foo", "-h"},
			ErrorExp: false,
		},
		testCase{
			Args:        []string{"foo"},
			DefaultStep: "",
			ErrorExp:    true,
		},
		testCase{
			Args:        []string{"foo", "--markdown"},
			DefaultStep: "",
			ErrorExp:    true,
		},
		testCase{
			Args:     []string{"foo", "--nonexistent-flag"},
			ErrorExp: true,
		},
		testCase{
			Args:     []string{"foo", "too", "many", "args"},
			ErrorExp: true,
		},
	}

	for i, tc := range testCases {
		t.Logf("test case %d", i)

		cli, err := NewDefaultCLI("foo", pcd, tc.DefaultStep)
		assert.Nil(err)

		var buf bytes.Buffer
		cli.out = &buf
		assert.Equal(tc.ErrorExp, (cli.Run(tc.Args) != nil))
		assert.Contains(buf.String(), "USAGE:")
	}
}
