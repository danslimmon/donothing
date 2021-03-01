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

// DefaultCLI should print usage when --help is passed or the args are wrong.
func TestDefaultCLI_PrintUsage(t *testing.T) {
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
		err = cli.Run(tc.Args)
		if tc.ErrorExp != (err != nil) {
			if err != nil {
				t.Logf("cli.Run returned unexpected error '%s'", err.Error())
				t.Fail()
			} else {
				t.Logf("cli.Run should have returned an error but didn't")
			}
		}
		assert.Contains(buf.String(), "USAGE:")
	}
}

// DefaultCLI should render a step when --markdown is passed
func TestDefaultCLI_Render(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	pcd := NewProcedure()
	pcd.Short("Procedure's short description")
	pcd.AddStep(func(step *Step) {
		step.Name("blahBlah")
		step.Short("the blahBlah step")
	})

	type testCase struct {
		// os.Args
		Args []string
		// The CLI's default step
		DefaultStep string
		// A function that makes assertions about the output of cli.Run
		Match func(s string)
		// Whether an error is expected from cli.Run
		ErrorExp bool
	}

	testCases := []testCase{
		testCase{
			Args:        []string{"foo", "--markdown"},
			DefaultStep: "root",
			Match: func(s string) {
				assert.Contains(s, "Procedure's short description")
				assert.Contains(s, "the blahBlah step")
			},
			ErrorExp: false,
		},
		testCase{
			Args:        []string{"foo", "--markdown"},
			DefaultStep: "",
			Match:       func(s string) {},
			ErrorExp:    true,
		},
		testCase{
			Args:        []string{"foo", "--markdown", "root.blahBlah"},
			DefaultStep: "",
			Match: func(s string) {
				assert.NotContains(s, "Procedure's short description")
				assert.Contains(s, "the blahBlah step")
			},
			ErrorExp: false,
		},
		testCase{
			Args:        []string{"foo", "--markdown"},
			DefaultStep: "root.blahBlah",
			Match: func(s string) {
				assert.NotContains(s, "Procedure's short description")
				assert.Contains(s, "the blahBlah step")
			},
			ErrorExp: false,
		},
		testCase{
			Args:        []string{"foo", "--markdown", "nonexistentStep"},
			DefaultStep: "",
			Match:       func(s string) {},
			ErrorExp:    true,
		},
		testCase{
			Args:        []string{"foo", "--markdown"},
			DefaultStep: "nonexistentStep",
			Match:       func(s string) {},
			ErrorExp:    true,
		},
	}

	for i, tc := range testCases {
		t.Logf("test case %d", i)

		cli, err := NewDefaultCLI("foo", pcd, tc.DefaultStep)
		assert.Nil(err)

		var buf bytes.Buffer
		cli.out = &buf
		err = cli.Run(tc.Args)
		if tc.ErrorExp != (err != nil) {
			if err != nil {
				t.Logf("cli.Run returned unexpected error '%s'", err.Error())
				t.Fail()
			} else {
				t.Logf("cli.Run should have returned an error but didn't")
			}
		}
		tc.Match(buf.String())
	}
}
