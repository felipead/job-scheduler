package scheduler

import "fmt"

// TriggerCallback receives the job ID, and the time when the job was triggered.
type TriggerCallback = func(string, Time)

type Job struct {
	ID              string
	OnTrigger       TriggerCallback
	IntervalMinutes int
	NextTime        Time
}

func (job *Job) Trigger(time Time) {
	job.logTrigger(time)

	if job.OnTrigger != nil {
		job.OnTrigger(job.ID, time)
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

	return fmt.Sprintf("%s {every %s}", job.ID, intervalText)
}

func (job *Job) logTrigger(time Time) {
	fmt.Printf("[%s] triggered → %s\n", time.String(), job.String())
}
