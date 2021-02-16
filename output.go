package donothing

// An OutputDef specifies a value that a step outputs for later consumption by another step.
type OutputDef struct {
	// The type for values of the output. Either "string" or "int"
	ValueType string

	// The output's name, which another step can refer to in an InputDef if it wants to use this
	// output's value as an input.
	Name string

	// A short description of the output.
	//
	// This will be used in the procedure's rendered documentation, and also as part of the prompt
	// during Procedure.Execute() if the output needs to be provided by the user.
	Short string
}

func NewOutputDef(valueType string, name, short string) OutputDef {
	return OutputDef{
		ValueType: valueType,
		Name:      name,
		Short:     short,
	}
}
