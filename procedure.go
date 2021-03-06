package donothing

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

// A Procedure is a sequence of Steps that can be executed or rendered to markdown.
type Procedure struct {
	// The root step of the procedure, of which all other steps are descendants.
	rootStep *Step

	stdin  io.Reader
	stdout io.Writer
}

// Short provides the procedure with a short description.
//
// The short description will be the title of the rendered markdown document when Render is called,
// so it should be concise and accurate.
func (pcd *Procedure) Short(s string) {
	pcd.rootStep.Short(s)
}

// GetShort returns the procedure's short description.
func (pcd *Procedure) GetShort() string {
	return pcd.rootStep.GetShort()
}

// Long provides the procedure with a long description.
//
// The long description will be shown to the user when they first execute the procedure. It will
// also be included in the opening section of the procedure's Markdown documentation.
//
// It should give an overview of the procedure's purpose and any important assumptions the procedure
// makes about the state of the world at the beginning of execution.
func (pcd *Procedure) Long(s string) {
	pcd.rootStep.Long(s)
}

// AddStep adds a step to the procedure.
//
// A new Step will be instantiated and passed to fn to be defined.
func (pcd *Procedure) AddStep(fn func(*Step)) {
	pcd.rootStep.AddStep(fn)
}

// GetStepByName returns the step with the given (absolute) name.
func (pcd *Procedure) GetStepByName(stepName string) (*Step, error) {
	var foundStep *Step
	err := pcd.rootStep.Walk(func(step *Step) error {
		absNmae := step.AbsoluteName()
		if absNmae == stepName {
			//if step.AbsoluteName() == stepName {
			foundStep = step
			// Return error to end walk. This error will be ignored since we have set foundStep to
			// something other than nil.
			return fmt.Errorf("")
		}
		return nil
	})

	if foundStep != nil {
		return foundStep, nil
	}
	if err != nil {
		return nil, err
	}
	return nil, fmt.Errorf("No step with name '%s'", stepName)
}

// Check validates that the procedure makes sense.
//
// If problems are found, it returns the list of problems along with an error.
//
// It checks the procedure against the following expectations:
//
//   1. Every step has a unique absolute name with no empty parts.
//   2. Every step has a short description
//   3. Every input has a name that matches the name of an output from a previous step.
func (pcd *Procedure) Check() ([]string, error) {
	steps := make(map[string]*Step)
	outputs := make(map[string]OutputDef)
	problems := make([]string, 0)

	err := pcd.rootStep.Walk(func(step *Step) error {
		absName := step.AbsoluteName()
		if absName[len(absName)-1:] == "." {
			if step.parent == nil {
				// I really hope this never happens. The root step should get its name from
				// donothing, not from the calling code.
				problems = append(problems, "Root step does not have name")
			} else {
				problems = append(problems, fmt.Sprintf("Child step of '%s' does not have name", step.parent.AbsoluteName()))
			}
		}

		if steps[absName] != nil {
			problems = append(problems, fmt.Sprintf("More than one step with name '%s'", absName))
		}
		steps[step.AbsoluteName()] = step

		if step.GetShort() == "" {
			problems = append(problems, fmt.Sprintf("Step '%s' has no Short value", absName))
		}

		for _, inputDef := range step.GetInputDefs() {
			matchingOutputDef, ok := outputs[inputDef.Name]
			if !ok {
				problems = append(problems, fmt.Sprintf(
					"Input '%s' of step '%s' does not refer to an output from any previous step",
					inputDef.Name,
					absName,
				))
				continue
			}
			if matchingOutputDef.ValueType != inputDef.ValueType {
				problems = append(problems, fmt.Sprintf(
					"Input '%s' of step '%s' has type '%s', but output '%s' has type '%s'",
					inputDef.Name,
					absName,
					inputDef.ValueType,
					matchingOutputDef.Name,
					matchingOutputDef.ValueType,
				))
			}
		}

		for _, outputDef := range step.GetOutputDefs() {
			outputs[outputDef.Name] = outputDef
		}

		return nil
	})
	if err != nil {
		return []string{}, fmt.Errorf("Error while checking procedure: %w", err)
	}

	if len(problems) > 0 {
		return problems, errors.New("Problems were found in the procedure")
	}
	return []string{}, nil
}

// Render prints the procedure's Markdown representation to f.
//
// Any occurrence of the string "@@" in the executed template output will be replaced with a
// backtick.
func (pcd *Procedure) Render(f io.Writer) error {
	return pcd.RenderStep(f, "root")
}

// RenderStep prints the given step from the procedure as Markdown to f.
//
// Any occurrence of the string "@@" in the executed template output will be replaced with a
// backtick.
func (pcd *Procedure) RenderStep(f io.Writer, stepName string) error {
	if _, err := pcd.Check(); err != nil {
		return err
	}

	tpl, err := DocTemplate()
	if err != nil {
		return err
	}

	step, err := pcd.GetStepByName(stepName)
	if err != nil {
		return err
	}
	tplData := NewStepTemplateData(step, nil, true)

	var b strings.Builder
	err = tpl.Execute(&b, tplData)
	if err != nil {
		return err
	}

	fmt.Fprintf(f, "%s", strings.Replace(b.String(), "@@", "`", -1))
	return nil
}

// Execute runs through the procedure step by step.
//
// The user will be prompted as necessary.
func (pcd *Procedure) Execute() error {
	return pcd.ExecuteStep("root")
}

// ExecuteStep runs through the given step.
//
// The user will be prompted as necessary.
func (pcd *Procedure) ExecuteStep(stepName string) error {
	if _, err := pcd.Check(); err != nil {
		return err
	}

	step, err := pcd.GetStepByName(stepName)
	if err != nil {
		return err
	}

	tpl, err := ExecTemplate()
	if err != nil {
		return err
	}

	step, err = pcd.GetStepByName(stepName)
	if err != nil {
		return err
	}

	var skipTo string
	step.Walk(func(walkStep *Step) error {
		if skipTo != "" && walkStep.AbsoluteName() != skipTo {
			fmt.Fprintf(pcd.stdout, "Skipping step '%s' on the way to '%s'\n", walkStep.AbsoluteName(), skipTo)
			return nil
		}

		tplData := NewStepTemplateData(walkStep, nil, false)

		var b bytes.Buffer
		err = tpl.Execute(&b, tplData)
		if err != nil {
			return err
		}
		fmt.Fprintf(pcd.stdout, "%s", strings.Replace(b.String(), "@@", "`", -1))

		promptResult := pcd.prompt()
		if promptResult.SkipOne {
			fmt.Fprintf(pcd.stdout, "Skipping step '%s' and its descendants\n", walkStep.AbsoluteName())
			return NoRecurse
		}
		skipTo = promptResult.SkipTo
		return nil
	})

	fmt.Fprintln(pcd.stdout, "Done.")
	return nil
}

// promptResult is the struct returned by Procedure.prompt.
//
// Procedure.Execute uses the contents of a promptResult to decide what to do next.
type promptResult struct {
	// Whether to skip this step and its descendants.
	SkipOne bool
	// The absolute name of the next step that should be executed.
	//
	// If empty, Execute should proceed normally in its walk.
	SkipTo string
}

// prompt prompts the user for the next action to take.
//
// If the user enters an invalid choice, prompt will inform them of this and re-prompt until a valid
// choice is entered.
func (pcd *Procedure) prompt() promptResult {
	// promptOnce prompts the user for input. It returns their input, trimmed of leading and
	// trailing whitespace.
	promptOnce := func() (string, error) {
		fmt.Fprintf(pcd.stdout, "\n\n[Enter] to proceed (or \"help\"): ")
		entry, err := bufio.NewReader(pcd.stdin).ReadBytes('\n')
		fmt.Fprintf(pcd.stdout, "\n")
		return strings.TrimSpace(string(entry)), err
	}

	for {
		entry, err := promptOnce()
		if err != nil {
			fmt.Fprintf(pcd.stdout, "Error reading input: %s\n", err.Error())
			continue
		}

		if entry == "" {
			// Proceed to the next step as normal
			return promptResult{}
		}
		if entry == "help" {
			// Print the help message and prompt again
			pcd.printPromptHelp()
		}
		if entry == "skip" {
			return promptResult{SkipOne: true}
		}
		if strings.HasPrefix(entry, "skipto ") {
			parts := strings.Split(entry, " ")
			if len(parts) != 2 || len(parts[1]) == 0 {
				fmt.Fprintf(pcd.stdout, "Invalid 'skipto' syntax; enter \"help\" for help\n")
			}
			return promptResult{SkipTo: parts[1]}
		}

		fmt.Fprintf(pcd.stdout, "Invalid choice; enter \"help\" for help\n")
	}
}

// printPromptHelp prints the help message for the Execute prompt.
func (pcd *Procedure) printPromptHelp() {
	fmt.Fprintf(pcd.stdout, `Options:

[Enter]			Proceed to the next step
skip			Skip this step and its descendants
skipto STEP 	Skip to the given step by absolute name
help			Print this help message`)
}

// NewProcedure returns a new procedure, ready to be given steps.
func NewProcedure() *Procedure {
	pcd := new(Procedure)
	pcd.rootStep = NewStep()
	pcd.rootStep.Name("root")
	pcd.stdin = os.Stdin
	pcd.stdout = os.Stdout
	return pcd
}
