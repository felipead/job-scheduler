package scheduler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSchedule_AddHourlyJob(t *testing.T) {
	schedule := NewSchedule()

	callbackCalled := false
	var callback = func(name string, time Time) {
		callbackCalled = true
	}

	schedule.AddHourlyJob("foobar", 17, callback)

	jobs := schedule.GetJobsAt(17)
	assert.NotNil(t, jobs)
	assert.Equal(t, jobs.Len(), 1)

	job := jobs.Front().Value.(Job)
	assert.Equal(t, job.Name, "foobar")
	assert.Equal(t, job.IntervalMinutes, 60)
	assert.Equal(t, job.NextTime, Time(17))

	job.Trigger(1000)
	assert.True(t, callbackCalled)
}

func TestSchedule_AddIntervalJob(t *testing.T) {
	schedule := NewSchedule()

	callbackCalled := false
	var callback = func(name string, time Time) {
		callbackCalled = true
	}

	schedule.AddIntervalJob("foobar", 25, 10, callback)

	time := Time(25 + 10)

	jobs := schedule.GetJobsAt(time)
	assert.NotNil(t, jobs)
	assert.Equal(t, jobs.Len(), 1)

	job := jobs.Front().Value.(Job)
	assert.Equal(t, job.Name, "foobar")
	assert.Equal(t, job.IntervalMinutes, 25)
	assert.Equal(t, job.NextTime, time)

	job.Trigger(1000)
	assert.True(t, callbackCalled)
}

func TestSchedule_AddIntervalJob_IntervalGreaterThan60Minutes(t *testing.T) {
	schedule := NewSchedule()

	callbackCalled := false
	var callback = func(name string, time Time) {
		callbackCalled = true
	}

	schedule.AddIntervalJob("foobar", 100, 15, callback)

	time := Time(100 + 15)
	hour := time.GetHour()
	minute := time.GetMinute()

	jobs := schedule.GetJobsAt(time)
	assert.NotNil(t, jobs)
	assert.Equal(t, jobs.Len(), 1)

	job := jobs.Front().Value.(Job)
	assert.Equal(t, job.Name, "foobar")
	assert.Equal(t, job.IntervalMinutes, 100)
	assert.Equal(t, job.NextTime.GetHour(), hour)
	assert.Equal(t, job.NextTime.GetMinute(), minute)

	job.Trigger(1000)
	assert.True(t, callbackCalled)
}

func TestSchedule_Reschedule_FreshBucket(t *testing.T) {
	schedule := NewSchedule()

	job := Job{
		Name:            "foobar",
		IntervalMinutes: 17,
		NextTime:        170,
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
		NextTime:        170,
	}

	schedule.Reschedule(job)

	jobs := schedule.GetJobsAt(50)
	assert.Equal(t, jobs.Len(), 3)

	rescheduledJob := jobs.Front().Next().Next().Value.(Job)
	assert.Equal(t, rescheduledJob, job)
}
