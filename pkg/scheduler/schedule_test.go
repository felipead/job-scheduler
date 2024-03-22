package scheduler

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSchedule_AddHourlyJob(t *testing.T) {
	schedule := NewSchedule()

	callbackCalled := false
	var callback = func(name string, time int, hour int, minute int) {
		callbackCalled = true
	}

	schedule.AddHourlyJob("foobar", 17, callback)

	jobs := schedule.GetJobsAt(17)
	assert.NotNil(t, jobs)
	assert.Equal(t, jobs.Len(), 1)

	job := jobs.Front().Value.(Job)
	assert.Equal(t, job.Name, "foobar")
	assert.Equal(t, job.IntervalMinutes, 60)
	assert.Equal(t, job.NextHour, 0)
	assert.Equal(t, job.NextMinute, 17)

	job.Trigger(1000, 1000/60, 1000%60)
	assert.True(t, callbackCalled)
}

func TestSchedule_AddIntervalJob(t *testing.T) {
	schedule := NewSchedule()

	callbackCalled := false
	var callback = func(name string, time int, hour int, minute int) {
		callbackCalled = true
	}

	schedule.AddIntervalJob("foobar", 25, 10, callback)

	minute := 25 + 10

	jobs := schedule.GetJobsAt(minute)
	assert.NotNil(t, jobs)
	assert.Equal(t, jobs.Len(), 1)

	job := jobs.Front().Value.(Job)
	assert.Equal(t, job.Name, "foobar")
	assert.Equal(t, job.IntervalMinutes, 25)
	assert.Equal(t, job.NextHour, 0)
	assert.Equal(t, job.NextMinute, minute)

	job.Trigger(1000, 1000/60, 1000%60)
	assert.True(t, callbackCalled)
}

func TestSchedule_AddIntervalJob_IntervalGreaterThan60Minutes(t *testing.T) {
	schedule := NewSchedule()

	callbackCalled := false
	var callback = func(name string, time int, hour int, minute int) {
		callbackCalled = true
	}

	schedule.AddIntervalJob("foobar", 100, 15, callback)

	hour := (100 + 15) / 60
	minute := (100 + 15) % 60

	jobs := schedule.GetJobsAt(minute)
	assert.NotNil(t, jobs)
	assert.Equal(t, jobs.Len(), 1)

	job := jobs.Front().Value.(Job)
	assert.Equal(t, job.Name, "foobar")
	assert.Equal(t, job.IntervalMinutes, 100)
	assert.Equal(t, job.NextHour, hour)
	assert.Equal(t, job.NextMinute, minute)

	job.Trigger(1000, 1000/60, 1000%60)
	assert.True(t, callbackCalled)
}

func TestSchedule_Reschedule_FreshBucket(t *testing.T) {
	schedule := NewSchedule()

	job := Job{
		Name:            "foobar",
		IntervalMinutes: 17,
		NextMinute:      50,
		NextHour:        2,
	}

	schedule.Reschedule(job)

	jobs := schedule.GetJobsAt(50)
	assert.Equal(t, jobs.Len(), 1)

	rescheduledJob := jobs.Front().Value.(Job)
	assert.Equal(t, rescheduledJob, job)
}

func TestSchedule_Reschedule_DirtyBucket(t *testing.T) {
	schedule := NewSchedule()
	schedule.AddIntervalJob("alice", 40, 10, nil)
	schedule.AddHourlyJob("bob", 50, nil)

	job := Job{
		Name:            "foobar",
		IntervalMinutes: 17,
		NextMinute:      50,
		NextHour:        2,
	}

	schedule.Reschedule(job)

	jobs := schedule.GetJobsAt(50)
	assert.Equal(t, jobs.Len(), 3)

	rescheduledJob := jobs.Front().Next().Next().Value.(Job)
	assert.Equal(t, rescheduledJob, job)
}
