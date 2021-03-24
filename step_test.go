package donothing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTrimCommonIndent(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	cases := []struct{ In, Out string }{
		struct{ In, Out string }{"", ""},
		struct{ In, Out string }{"hello", "hello"},
		struct{ In, Out string }{"    hello", "hello"},
		// 3 tabs
		struct{ In, Out string }{"\t\t\tEnter your phone number, without area code. Formatting doesn't matter.\n", "Enter your phone number, without area code. Formatting doesn't matter.\n"},
		struct{ In, Out string }{"    hello\n    goodbye", "hello\ngoodbye"},
		struct{ In, Out string }{"    hello\n\n    goodbye", "hello\n\ngoodbye"},
		// the "middle" line has four spaces followed by a tab.
		struct{ In, Out string }{"    hello\n    \tmiddle\n    goodbye", "hello\n\tmiddle\ngoodbye"},
		struct{ In, Out string }{"hello\n    middle\ngoodbye", "hello\n    middle\ngoodbye"},
		struct{ In, Out string }{"    hello\n        middle\n    goodbye\n        again\n    bye for real", "hello\n    middle\ngoodbye\n    again\nbye for real"},
	}

	for _, c := range cases {
		assert.Equal(c.Out, (&Step{}).trimCommonIndent(c.In))
	}
}

func TestStep_Pos(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	pcd := NewProcedure()
	pcd.AddStep(func(step *Step) {})
	pcd.AddStep(func(step *Step) {
		step.Name("grandparent")
		step.AddStep(func(step *Step) {})
		step.AddStep(func(step *Step) {})
		step.AddStep(func(step *Step) {
			step.Name("parent")
			step.AddStep(func(step *Step) {
				step.Name("myStep")
			})
		})
	})

	must := func(step *Step, err error) *Step {
		assert.Nil(err)
		return step
	}

	assert.Equal([]int{}, must(pcd.GetStepByName("root")).Pos())
	assert.Equal([]int{1}, must(pcd.GetStepByName("root.grandparent")).Pos())
	assert.Equal([]int{1, 2}, must(pcd.GetStepByName("root.grandparent.parent")).Pos())
	assert.Equal([]int{1, 2, 0}, must(pcd.GetStepByName("root.grandparent.parent.myStep")).Pos())
}

// Pos panics if step's parent's "children" slice doesn't contain step
func TestStep_Pos_MissingFromParent(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	pcd0 := NewProcedure()
	pcd1 := NewProcedure()
	pcd1.AddStep(func(step *Step) {
		step.Name("foo")
	})
	assert.Panics(func() {
		fooStep, err := pcd1.GetStepByName("root.foo")
		assert.Nil(err)
		fooStep.parent = pcd0.rootStep
		_ = fooStep.Pos()
	})
}
