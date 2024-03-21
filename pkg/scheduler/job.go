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
	job.logTrigger(hour, minute)

	if job.OnTrigger != nil {
		job.OnTrigger(job.Name, time, hour, minute)
	}
}

func (job *Job) logTrigger(hour int, minute int) {
	var intervalText string
	if interval := job.IntervalMinutes; interval != 60 {
		word := "minutes"
		if interval == 1 {
			word = "minute"
		}
		intervalText = fmt.Sprintf("%d %s", interval, word)
	} else {
		intervalText = "hour"
	}

	fmt.Printf("[%02d:%02d] %v triggered → {every %s}\n", hour, minute, job.Name, intervalText)
}
