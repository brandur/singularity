package main

import (
	"fmt"
	"testing"

	"github.com/brandur/singularity"
	"github.com/brandur/sorg/pool"
	assert "github.com/stretchr/testify/require"
)

func TestGetLocals(t *testing.T) {
	locals := getLocals("Title", map[string]interface{}{
		"Foo": "Bar",
	})

	assert.Equal(t, "Bar", locals["Foo"])
	assert.Equal(t, singularity.Release, locals["Release"])
	assert.Equal(t, "Title", locals["Title"])
}

func TestIsHidden(t *testing.T) {
	assert.Equal(t, true, isHidden(".gitkeep"))
	assert.Equal(t, false, isHidden("article"))
}

func TestRunTasks(t *testing.T) {
	conf.Concurrency = 3

	//
	// Success case
	//

	tasks := []*pool.Task{
		pool.NewTask(func() error { return nil }),
		pool.NewTask(func() error { return nil }),
		pool.NewTask(func() error { return nil }),
	}
	assert.Equal(t, true, runTasks(tasks))

	//
	// Failure case (1 error)
	//

	tasks = []*pool.Task{
		pool.NewTask(func() error { return nil }),
		pool.NewTask(func() error { return nil }),
		pool.NewTask(func() error { return fmt.Errorf("error") }),
	}
	assert.Equal(t, false, runTasks(tasks))

	//
	// Failure case (11 errors)
	//
	// Here we'll exit with a "too many errors" message.
	//

	tasks = []*pool.Task{
		pool.NewTask(func() error { return fmt.Errorf("error") }),
		pool.NewTask(func() error { return fmt.Errorf("error") }),
		pool.NewTask(func() error { return fmt.Errorf("error") }),
		pool.NewTask(func() error { return fmt.Errorf("error") }),
		pool.NewTask(func() error { return fmt.Errorf("error") }),
		pool.NewTask(func() error { return fmt.Errorf("error") }),
		pool.NewTask(func() error { return fmt.Errorf("error") }),
		pool.NewTask(func() error { return fmt.Errorf("error") }),
		pool.NewTask(func() error { return fmt.Errorf("error") }),
		pool.NewTask(func() error { return fmt.Errorf("error") }),
		pool.NewTask(func() error { return fmt.Errorf("error") }),
	}
	assert.Equal(t, false, runTasks(tasks))
}
