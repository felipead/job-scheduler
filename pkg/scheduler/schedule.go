package scheduler

import "container/list"

var hourlySchedule = make(map[int][]HourlyJob)
var intervalSchedule = make(map[int]*list.List)

func ScheduleHourlyJob(name string, minute int, onTrigger TriggerCallback) {
	// TODO: in a production system, I would instead add validation before accepting this job: 0 <= minute < 60
	minute = minute % 60

	jobs := hourlySchedule[minute]
	if jobs == nil {
		jobs = make([]HourlyJob, 0)
	}

	hourlySchedule[minute] = append(jobs, HourlyJob{
		Name:      name,
		OnTrigger: onTrigger,
	})
}

func ScheduleIntervalJob(name string, intervalMinutes int, offset int, onTrigger TriggerCallback) {
	nextTime := offset + intervalMinutes
	nextMinute := nextTime % 60
	nextHour := nextTime / 60

	jobs := intervalSchedule[nextMinute]
	if jobs == nil {
		jobs = list.New()
		intervalSchedule[nextMinute] = jobs
	}

	jobs.PushBack(IntervalJob{
		Name:            name,
		OnTrigger:       onTrigger,
		IntervalMinutes: intervalMinutes,
		NextMinute:      nextMinute,
		NextHour:        nextHour,
	})
}

func rescheduleIntervalJob(job IntervalJob) {
	jobs := intervalSchedule[job.NextMinute]
	if jobs == nil {
		jobs = list.New()
		intervalSchedule[job.NextMinute] = jobs
	}
	jobs.PushBack(job)
}

func getHourlyScheduleAt(minute int) []HourlyJob {
	return hourlySchedule[minute]
}

func getIntervalScheduleAt(minute int) *list.List {
	return intervalSchedule[minute]
}
