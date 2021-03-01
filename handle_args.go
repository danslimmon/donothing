package donothing

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
	"text/template"
)

// DefaultCLI is a default CLI for do-nothing scripts.
type DefaultCLI struct {
	ExecName    string
	Pcd         *Procedure
	DefaultStep string

	// The place we'll write output to. Can be swapped out for testing.
	out io.Writer
}

// Usage returns the usage message.
func (cli *DefaultCLI) Usage() string {
	tplStr := `USAGE: {{.ExecName}} [options] {{if .DefaultStep}}[STEP_NAME]{{else}}STEP_NAME{{end}}

{{if .Pcd.GetShort -}}
{{.Pcd.GetShort}}

{{end -}}
OPTIONS: 
    --markdown    Instead of executing the procedure, print its Markdown documentation to stdout
    --help        Print usage message`
	//tpl := template.Must(template.New("usage").Parse(tplStr))
	tpl, err := template.New("usage").Parse(tplStr)
	if err != nil {
		return err.Error()
	}

	var buf bytes.Buffer
	if err := tpl.Execute(&buf, cli); err != nil {
		return err.Error()
	}
	return buf.String()
}

// Run parses arguments and runs the appropriate actions.
//
// args is the content of os.Args.
func (cli *DefaultCLI) Run(args []string) error {
	if len(args) <= 0 {
		return fmt.Errorf("Must have at least 1 argument")
	}

	flags := make([]string, 0)
	nonFlags := make([]string, 0)
	for _, arg := range args[1:] {
		if strings.IndexRune(arg, '-') == 0 {
			flags = append(flags, arg)
		} else {
			nonFlags = append(nonFlags, arg)
		}
	}

	for _, flag := range flags {
		if flag == "-h" || flag == "--help" {
			fmt.Fprintln(cli.out, cli.Usage())
			return nil
		}
	}

	// Keys of opts are valid flags. Any other flag is an error.
	//
	// At the end of this stanza, the value corresponding to each flag will be true iff the flag was
	// passed.
	opts := map[string]bool{
		"--markdown": false,
	}
	for _, flag := range flags {
		if _, ok := opts[flag]; ok {
			opts[flag] = true
		} else {
			fmt.Fprintln(cli.out, cli.Usage())
			return fmt.Errorf("Unknown flag '%s'", flag)
		}
	}

	if len(nonFlags) == 0 && cli.DefaultStep == "" {
		fmt.Fprintln(cli.out, cli.Usage())
		return fmt.Errorf("Must specify STEP_NAME")
	}

	if len(nonFlags) > 1 {
		fmt.Fprintln(cli.out, cli.Usage())
		return fmt.Errorf("Extraneous arguments passed: %v", nonFlags[1:])
	}

	stepName := cli.DefaultStep
	if len(nonFlags) >= 1 {
		stepName = nonFlags[0]
	}

	if opts["--markdown"] {
		return cli.Pcd.RenderStep(cli.out, stepName)
	}
	return cli.Pcd.ExecuteStep(stepName)
}

// NewDefaultCLI returns a DefaultCLI instance initialized with the given executable name.
//
// execName is the name of the executable that has imported donothing. pcd is the procedure to run
// actions against. defaultStep is the step to execute if the user doesn't specify STEP_NAME; if
// defaultStep is "", omission of STEP_NAME from the invocation will trigger an error.
func NewDefaultCLI(execName string, pcd *Procedure, defaultStep string) (*DefaultCLI, error) {
	if pcd == nil {
		return nil, fmt.Errorf("failed to initialize default CLI: procedure must not be nil")
	}
	if _, err := pcd.Check(); err != nil {
		return nil, err
	}
	return &DefaultCLI{
		ExecName:    execName,
		Pcd:         pcd,
		DefaultStep: defaultStep,

		out: os.Stdout,
	}, nil
}

// HandleArgs parses command line arguments and runs the appropriate actions.
//
// In other words, this function implements a default CLI for do-nothing scripts. Packages that
// import donothing may use this default CLI, or may implement their own CLI by directly calling
// Procedure.Execute, Procedure.Render, and so on.
//
// For details of this default CLI, see the documentation of DefaultCLI and its methods.
//
// args is the content of os.Args. pcd is the procedure to run actions against. defaultStep is the
// step to execute if the user doesn't specify STEP_NAME; if defaultStep is "", omission of
// STEP_NAME from the invocation will trigger an error.
func HandleArgs(args []string, pcd *Procedure, defaultStep string) error {
	args = args[:]
	if len(args) <= 0 {
		return fmt.Errorf("Failed to determine executable name from arguments: args slice too short")
	}
	cli, err := NewDefaultCLI(path.Base(args[0]), pcd, defaultStep)
	if err != nil {
		return err
	}

	return cli.Run(args)
}
