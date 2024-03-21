package scheduler

import "container/list"

var hourlySchedule = make(map[int][]HourlyJob)
var intervalSchedule = make(map[int]*list.List)

func ScheduleHourlyJob(jobName string, minute int) {
	// TODO: in a production system, I would instead add validation before accepting this job: 0 <= minute < 60
	minute = minute % 60

	jobs := hourlySchedule[minute]
	if jobs == nil {
		jobs = make([]HourlyJob, 0)
	}

	hourlySchedule[minute] = append(jobs, HourlyJob{
		Name: jobName,
	})
}

func ScheduleIntervalJob(jobName string, intervalMinutes int, offsetMinutes int) {
	nextTime := intervalMinutes + offsetMinutes
	nextMinute := nextTime % 60
	nextHour := nextTime / 60

	jobs := intervalSchedule[nextMinute]
	if jobs == nil {
		jobs = list.New()
		intervalSchedule[nextMinute] = jobs
	}

	jobs.PushBack(IntervalJob{
		Name:            jobName,
		IntervalMinutes: intervalMinutes,
		NextMinute:      nextMinute,
		NextHour:        nextHour,
	})
}

func reschedule(job IntervalJob) {
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
