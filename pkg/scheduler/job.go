package scheduler

import "fmt"

type Job interface {
	Trigger(hour int, minute int)
}

type HourlyJob struct {
	Job
	Name string
}

func (job *HourlyJob) Trigger(hour int, minute int) {
	fmt.Printf("[%02d:%02d] %v triggered (hourly)\n", hour, minute, job.Name)
}

type IntervalJob struct {
	Job
	Name            string
	IntervalMinutes int
	NextMinute      int
	NextHour        int
}

func (job *IntervalJob) Trigger(hour int, minute int) {
	fmt.Printf("[%02d:%02d] %v triggered (%v mins interval)\n", hour, minute, job.Name, job.IntervalMinutes)
}
