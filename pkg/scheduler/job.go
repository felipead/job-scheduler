package scheduler

import "fmt"

// TriggerCallback receives the job name, the absolute time (minutes), the hour and the minute of the hour - all
// relative to the time the job was triggered.
type TriggerCallback = func(string, int)

type Job struct {
	Name            string
	OnTrigger       TriggerCallback
	IntervalMinutes int
	NextMinute      int
	NextHour        int
}

func (job *Job) Trigger(time int) {
	job.logTrigger(time)

	if job.OnTrigger != nil {
		job.OnTrigger(job.Name, time)
	}
}

func (job *Job) logTrigger(time int) {
	hour := time / 60
	minute := time % 60

	intervalText := "hour"
	if interval := job.IntervalMinutes; interval != 60 {
		word := "minutes"
		if interval == 1 {
			word = "minute"
		}
		intervalText = fmt.Sprintf("%d %s", interval, word)
	}

	fmt.Printf("[%02d:%02d] %v triggered → {every %s}\n", hour, minute, job.Name, intervalText)
}
