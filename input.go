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

	// Whether the input is required by the step
	required bool
}

// NewInputDef returns an InputDef struct describing a step input.
func NewInputDef(valueType string, name string, required bool) InputDef {
	return InputDef{
		valueType: valueType,
		name:      name,
		required:  required,
	}
}
