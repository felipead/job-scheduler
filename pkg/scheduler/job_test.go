package scheduler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJob_Trigger_ShouldCallCallback(t *testing.T) {
	name := "foobar"
	time := 1000
	hour := time / 60
	minute := time % 60

	callbackCalled := false
	onTrigger := func(_name string, _time int, _hour int, _minute int) {
		assert.Equal(t, name, _name)
		assert.Equal(t, time, _time)
		assert.Equal(t, hour, _hour)
		assert.Equal(t, minute, _minute)
		callbackCalled = true
	}

	job := Job{
		Name:            name,
		OnTrigger:       onTrigger,
		IntervalMinutes: 50,
	}

	job.Trigger(time, hour, minute)

	assert.True(t, callbackCalled)
}

func TestJob_Trigger_WhenCallbackIsNotDefined(t *testing.T) {
	job := Job{
		Name:            "foobar",
		OnTrigger:       nil,
		IntervalMinutes: 50,
	}

	time := 1000
	hour := time / 60
	minute := time % 60

	called := false
	test := func() {
		job.Trigger(time, hour, minute)
		called = true
	}

	assert.NotPanics(t, test)
	assert.True(t, called)
}
