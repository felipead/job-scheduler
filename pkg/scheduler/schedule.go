package scheduler

import "container/list"

type Schedule struct {
	interval map[int]*list.List
}

func NewSchedule() *Schedule {
	return &Schedule{
		interval: make(map[int]*list.List),
	}
}

func (s *Schedule) AddHourlyJob(name string, minute int, onTrigger TriggerCallback) {
	// TODO: in a production system, I would instead add validation before accepting this job: 0 <= minute < 60
	minute = minute % 60

	nextTime := minute
	nextHour := nextTime / 60
	nextMinute := nextTime % 60

	s.addJob(name, 60, nextHour, nextMinute, onTrigger)
}

func (s *Schedule) AddIntervalJob(name string, intervalMinutes int, offset int, onTrigger TriggerCallback) {
	nextTime := offset + intervalMinutes
	nextHour := nextTime / 60
	nextMinute := nextTime % 60

	s.addJob(name, intervalMinutes, nextHour, nextMinute, onTrigger)
}

func (s *Schedule) addJob(name string, intervalMinutes int, nextHour int, nextMinute int, onTrigger TriggerCallback) {
	jobs := s.interval[nextMinute]
	if jobs == nil {
		jobs = list.New()
		s.interval[nextMinute] = jobs
	}

	jobs.PushBack(Job{
		Name:            name,
		OnTrigger:       onTrigger,
		IntervalMinutes: intervalMinutes,
		NextHour:        nextHour,
		NextMinute:      nextMinute,
	})
}

func (s *Schedule) RescheduleIntervalJob(job Job) {
	jobs := s.interval[job.NextMinute]
	if jobs == nil {
		jobs = list.New()
		s.interval[job.NextMinute] = jobs
	}
	jobs.PushBack(job)
}

func (s *Schedule) GetScheduledJobsAt(minute int) *list.List {
	return s.interval[minute]
}
