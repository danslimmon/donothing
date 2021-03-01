package donothing

import (
	"bytes"
	"fmt"
	"path"
	"text/template"
)

// DefaultCLI is a default CLI for do-nothing scripts.
type DefaultCLI struct {
	ExecName    string
	Pcd         *Procedure
	DefaultStep string
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

// NewDefaultCLI returns a DefaultCLI instance initialized with the given executable name.
//
// execName is the name of the executable that has imported donothing. pcd is the procedure to run
// actions against. defaultStep is the step to execute if the user doesn't specify STEP_NAME; if
// defaultStep is "", omission of STEP_NAME from the invocation will trigger an error.
func NewDefaultCLI(execName string, pcd *Procedure, defaultStep string) (*DefaultCLI, error) {
	if _, err := pcd.Check(); err != nil {
		return nil, err
	}
	return &DefaultCLI{
		ExecName:    execName,
		Pcd:         pcd,
		DefaultStep: defaultStep,
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
	if len(args) <= 0 {
		return fmt.Errorf("Failed to determine executable name from arguments: args slice too short")
	}
	cli, err := NewDefaultCLI(path.Base(args[0]), pcd, defaultStep)
	if err != nil {
		return err
	}
	fmt.Println(cli.Usage())
	return nil
}
