package donothing

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// readThrough reads a byte at a time from r until it encounters suffix.
//
// If timeout is reached, or there's a problem reading from r, then readThrough returns whatever
// it's read so far in addition to an error.
func readThrough(r io.ByteReader, suffix []byte, timeout time.Duration) ([]byte, error) {
	after := time.After(timeout)
	rslt := make([]byte, 0)
	for {
		select {
		case <-after:
			return rslt[:], errors.New("readThrough timed out")
		default:
			b, err := r.ReadByte()
			if err != nil {
				return rslt[:], err
			}
			rslt = append(rslt, b)
			if bytes.HasSuffix(rslt, suffix) {
				return rslt[:], nil
			}
		}
	}
}

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

	_, err := pcd.GetStepByName("root.nonexistentStep")
	assert.NotNil(err)
}

// ExecuteStep should print the step and its children, prompting after each.
//
// This tests a procedure with only a single step (the root step).
func TestProcedure_ExecuteStep_Single(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	short := "root step"
	long := "blah blah blah\n\nthis is @@all@@ very interesting to you"

	pcd := NewProcedure()
	pcd.Short(short)
	pcd.Long(long)

	stdinReader, stdinWriter := io.Pipe()
	stdoutReader, stdoutWriter := io.Pipe()
	stdoutBufReader := bufio.NewReader(stdoutReader)
	pcd.stdin = stdinReader
	pcd.stdout = stdoutWriter

	go pcd.ExecuteStep("root")

	// Validate the initial output on stdout. At the end of this stanza, ExecuteStep() should be
	// waiting for our input on stdin.
	//
	// This will be wrong if the output contains ": " before we're prompted. Cross that bridge
	// if we come to it.
	output, err := readThrough(stdoutBufReader, []byte(": "), 5*time.Second)
	assert.Nil(err)
	assert.Contains(string(output), "`all`")

	stdinWriter.Write([]byte("\n"))

	// Validate the final output on stdout.
	_, err = readThrough(stdoutBufReader, []byte("Done.\n"), 5*time.Second)
	assert.Nil(err)
}

// ExecuteStep should print the step and its children, prompting after each.
//
// This tests a procedure with three steps, one nested in another nested in another.
func TestProcedure_ExecuteStep_Nested(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	shortTpl := "short %d"
	longTpl := "long %d"

	pcd := NewProcedure()
	pcd.Short(fmt.Sprintf(shortTpl, 0))
	pcd.Long(fmt.Sprintf(longTpl, 0))
	pcd.AddStep(func(step *Step) {
		step.Name("childStep")
		step.Short(fmt.Sprintf(shortTpl, 1))
		step.Long(fmt.Sprintf(longTpl, 1))
		step.AddStep(func(step *Step) {
			step.Name("grandchildStep")
			step.Short(fmt.Sprintf(shortTpl, 2))
			step.Long(fmt.Sprintf(longTpl, 2))
		})
	})

	stdinReader, stdinWriter := io.Pipe()
	stdoutReader, stdoutWriter := io.Pipe()
	stdoutBufReader := bufio.NewReader(stdoutReader)
	pcd.stdin = stdinReader
	pcd.stdout = stdoutWriter

	go pcd.ExecuteStep("root")

	// Validate the initial output on stdout. At the end of this stanza, ExecuteStep() should be
	// waiting for our input on stdin.
	//
	// This will be wrong if the output contains ": " before we're prompted. Cross that bridge
	// if we come to it.
	output, err := readThrough(stdoutBufReader, []byte(": "), 5*time.Second)
	assert.Nil(err)
	assert.Contains(string(output), fmt.Sprintf(shortTpl, 0))
	assert.Contains(string(output), fmt.Sprintf(longTpl, 0))

	stdinWriter.Write([]byte("\n"))

	output, err = readThrough(stdoutBufReader, []byte(": "), 5*time.Second)
	assert.Nil(err)
	assert.Contains(string(output), fmt.Sprintf(shortTpl, 1))
	assert.Contains(string(output), fmt.Sprintf(longTpl, 1))

	stdinWriter.Write([]byte("\n"))

	output, err = readThrough(stdoutBufReader, []byte(": "), 5*time.Second)
	assert.Nil(err)
	assert.Contains(string(output), fmt.Sprintf(shortTpl, 2))
	assert.Contains(string(output), fmt.Sprintf(longTpl, 2))

	stdinWriter.Write([]byte("\n"))

	_, err = readThrough(stdoutBufReader, []byte("Done.\n"), 5*time.Second)
	assert.Nil(err)
}
