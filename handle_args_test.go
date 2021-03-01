package donothing

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultCLI(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	pcd := NewProcedure()
	pcd.Short("Procedure's short description")

	cli, err := NewDefaultCLI("foo", pcd, "root/stepName")
	assert.Nil(err)
	assert.Equal(`USAGE: foo [options] [STEP_NAME]

Procedure's short description

OPTIONS: 
    --markdown    Instead of executing the procedure, print its Markdown documentation to stdout
	--help        Print usage message`, cli.Usage())
}
