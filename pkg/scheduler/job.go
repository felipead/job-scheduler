package scheduler

import "fmt"

// TriggerCallback receives the job name, the absolute time (minutes), the hour and the minute of the hour - all
// relative to the time the job was triggered.
type TriggerCallback = func(string, Time)

type Job struct {
	Name            string
	OnTrigger       TriggerCallback
	IntervalMinutes int
	NextTime        Time
}

func (job *Job) Trigger(time Time) {
	job.logTrigger(time)

	if job.OnTrigger != nil {
		job.OnTrigger(job.Name, time)
	}
}

func (job *Job) String() string {
	intervalText := "hour"
	if interval := job.IntervalMinutes; interval != 60 {
		word := "minutes"
		if interval == 1 {
			word = "minute"
		}
		intervalText = fmt.Sprintf("%d %s", interval, word)
	}

	return fmt.Sprintf("%s {every %s}", job.Name, intervalText)
}

func (job *Job) logTrigger(time Time) {
	fmt.Printf("[%s] triggered → %s\n", time.String(), job.String())
}
