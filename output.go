package donothing

// An InputDef specifies a value that a step can receive.
type InputDef struct {
	// The type for values of the input. Either "string" or "int"
	valueType string

	// The input's name.
	//
	// If name matches the name of an output from a previous step, then the input will automatically
	// take the value of that output. Otherwise, the user will be prompted for a value.
	name string

	// A short description of the input.
	//
	// This will be used in the procedure's rendered documentation, and also as part of the prompt
	// during Procedure.Execute() if the input doesn't have a value yet.
	short string

	// Whether the input is required by the step
	required bool
}

func NewInputDef(valueType string, name, short string, required bool) InputDef {
	return InputDef{
		valueType: valueType,
		name:      name,
		short:     short,
		required:  required,
	}
}
