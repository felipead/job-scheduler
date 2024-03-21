package scheduler

import "container/list"

type Schedule struct {
	hourly   map[int][]HourlyJob
	interval map[int]*list.List
}

func NewSchedule() *Schedule {
	return &Schedule{
		hourly:   make(map[int][]HourlyJob),
		interval: make(map[int]*list.List),
	}
}

func (s *Schedule) AddHourlyJob(name string, minute int, onTrigger TriggerCallback) {
	// TODO: in a production system, I would instead add validation before accepting this job: 0 <= minute < 60
	minute = minute % 60

	jobs := s.hourly[minute]
	if jobs == nil {
		jobs = make([]HourlyJob, 0)
	}

	s.hourly[minute] = append(jobs, HourlyJob{
		Name:      name,
		OnTrigger: onTrigger,
	})
}

func (s *Schedule) AddIntervalJob(name string, intervalMinutes int, offset int, onTrigger TriggerCallback) {
	nextTime := offset + intervalMinutes
	nextMinute := nextTime % 60
	nextHour := nextTime / 60

	jobs := s.interval[nextMinute]
	if jobs == nil {
		jobs = list.New()
		s.interval[nextMinute] = jobs
	}

	jobs.PushBack(IntervalJob{
		Name:            name,
		OnTrigger:       onTrigger,
		IntervalMinutes: intervalMinutes,
		NextMinute:      nextMinute,
		NextHour:        nextHour,
	})
}

func (s *Schedule) RescheduleIntervalJob(job IntervalJob) {
	jobs := s.interval[job.NextMinute]
	if jobs == nil {
		jobs = list.New()
		s.interval[job.NextMinute] = jobs
	}
	jobs.PushBack(job)
}

func (s *Schedule) GetHourlyJobsAt(minute int) []HourlyJob {
	return s.hourly[minute]
}

func (s *Schedule) GetIntervalJobsAt(minute int) *list.List {
	return s.interval[minute]
}
