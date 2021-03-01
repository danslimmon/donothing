package donothing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// GetStepByName should return the step with the given name
func TestProcedure_GetStepByName(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	pcd := NewProcedure()
	pcd.Short("The stanky leg")
	pcd.AddStep(func(step *Step) {
		step.Name("maximizeStank")
		step.Short("Maximize leg stankiness")
	})
	pcd.AddStep(func(step *Step) {
		step.Name("repeat")
		step.Short("Repeat")
	})

	for _, name := range []string{"root", "root.maximizeStank", "root.repeat"} {
		step, err := pcd.GetStepByName(name)
		assert.Nil(err)
		assert.Equal(name, step.AbsoluteName())
	}
}

// GetStepByName should return an error if the step doesn't exist
func TestProcedure_GetStepByName_Error(t *testing.T) {
}
