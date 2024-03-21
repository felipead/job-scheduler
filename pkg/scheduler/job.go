package scheduler

import "fmt"

// TriggerCallback receives the job name, the absolute time (minutes), the hour and the minute of the hour - all
// relative to the time the job was triggered.
type TriggerCallback = func(string, int, int, int)

type Job interface {
	Trigger(hour int, minute int)
}

type HourlyJob struct {
	Job
	Name      string
	OnTrigger TriggerCallback
}

func (job *HourlyJob) Trigger(time int, hour int, minute int) {
	fmt.Printf("[%02d:%02d] %v triggered (hourly)\n", hour, minute, job.Name)
	if job.OnTrigger != nil {
		job.OnTrigger(job.Name, time, hour, minute)
	}
}

type IntervalJob struct {
	Job
	Name            string
	OnTrigger       TriggerCallback
	IntervalMinutes int
	NextMinute      int
	NextHour        int
}

func (job *IntervalJob) Trigger(time int, hour int, minute int) {
	fmt.Printf("[%02d:%02d] %v triggered (%v mins interval)\n", hour, minute, job.Name, job.IntervalMinutes)
	if job.OnTrigger != nil {
		job.OnTrigger(job.Name, time, hour, minute)
	}
}
