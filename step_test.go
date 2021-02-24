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
