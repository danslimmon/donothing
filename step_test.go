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
		struct{ In, Out string }{"    hello\n    goodbye", "hello\ngoodbye"},
		struct{ In, Out string }{"    hello\n\n    goodbye", "hello\n\ngoodbye"},
		// the "middle" line has four spaces followed by a tab.
		struct{ In, Out string }{"    hello\n    	middle\n    goodbye", "hello\n	middle\ngoodbye"},
		struct{ In, Out string }{"hello\n    middle\ngoodbye", "hello\n    middle\ngoodbye"},
	}

	for _, c := range cases {
		assert.Equal(c.Out, (&Step{}).trimCommonIndent(c.In))
	}
}
