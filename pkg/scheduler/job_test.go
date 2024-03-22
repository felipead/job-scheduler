package scheduler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJob_Trigger_ShouldCallCallback(t *testing.T) {
	name := "foobar"
	time := 1000

	callbackCalled := false
	onTrigger := func(_name string, _time int) {
		assert.Equal(t, name, _name)
		assert.Equal(t, time, _time)
		callbackCalled = true
	}

	job := Job{
		Name:            name,
		OnTrigger:       onTrigger,
		IntervalMinutes: 50,
	}

	job.Trigger(time)

	assert.True(t, callbackCalled)
}

func TestJob_Trigger_WhenCallbackIsNotDefined(t *testing.T) {
	job := Job{
		Name:            "foobar",
		OnTrigger:       nil,
		IntervalMinutes: 50,
	}

	time := 1000

	called := false
	test := func() {
		job.Trigger(time)
		called = true
	}

	assert.NotPanics(t, test)
	assert.True(t, called)
}
