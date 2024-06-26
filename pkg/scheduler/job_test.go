package scheduler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJob_Trigger_ShouldCallCallback(t *testing.T) {
	id := "foobar"
	time := Time(1000)

	callbackCalled := false
	onTrigger := func(_id string, _time Time) {
		assert.Equal(t, id, _id)
		assert.Equal(t, time, _time)
		callbackCalled = true
	}

	job := Job{
		ID:              id,
		OnTrigger:       onTrigger,
		IntervalMinutes: 50,
	}

	job.Trigger(time)

	assert.True(t, callbackCalled)
}

func TestJob_Trigger_WhenCallbackIsNotDefined(t *testing.T) {
	job := Job{
		ID:              "foobar",
		OnTrigger:       nil,
		IntervalMinutes: 50,
	}

	time := Time(1000)

	called := false
	test := func() {
		job.Trigger(time)
		called = true
	}

	assert.NotPanics(t, test)
	assert.True(t, called)
}

func TestJob_String_Hourly(t *testing.T) {
	job := Job{
		ID:              "foo-bar",
		IntervalMinutes: 60,
	}

	assert.Equal(t, job.String(), "foo-bar {every hour}")
}

func TestJob_String_IntervalPlural(t *testing.T) {
	job := Job{
		ID:              "foo-bar",
		IntervalMinutes: 17,
	}

	assert.Equal(t, job.String(), "foo-bar {every 17 minutes}")
}

func TestJob_String_IntervalSingular(t *testing.T) {
	job := Job{
		ID:              "foo-bar",
		IntervalMinutes: 1,
	}

	assert.Equal(t, job.String(), "foo-bar {every 1 minute}")
}
