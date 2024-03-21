package scheduler

import "fmt"

// TriggerCallback receives the job name, the absolute time (minutes), the hour and the minute of the hour - all
// relative to the time the job was triggered.
type TriggerCallback = func(string, int, int, int)

type Job struct {
	Name            string
	OnTrigger       TriggerCallback
	IntervalMinutes int
	NextMinute      int
	NextHour        int
}

func (job *Job) Trigger(time int, hour int, minute int) {
	fmt.Printf("[%02d:%02d] %v triggered (%v mins interval)\n", hour, minute, job.Name, job.IntervalMinutes)

	if job.OnTrigger != nil {
		job.OnTrigger(job.Name, time, hour, minute)
	}
}
