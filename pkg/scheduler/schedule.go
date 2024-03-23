package scheduler

import "container/list"

type Schedule struct {
	buckets map[int]*list.List
}

func NewSchedule() *Schedule {
	return &Schedule{
		buckets: make(map[int]*list.List),
	}
}

func (s *Schedule) AddHourlyJob(name string, minute int, onTrigger TriggerCallback) {
	// TODO: in a production system, I would instead add validation before accepting this job: 0 <= minute < 60
	minute = minute % 60

	nextTime := Time(minute)
	nextHour := nextTime.GetHour()
	nextMinute := nextTime.GetMinute()

	s.addJob(name, 60, nextHour, nextMinute, onTrigger)
}

func (s *Schedule) AddIntervalJob(name string, intervalMinutes int, offset int, onTrigger TriggerCallback) {
	nextTime := Time(offset + intervalMinutes)
	nextHour := nextTime.GetHour()
	nextMinute := nextTime.GetMinute()

	s.addJob(name, intervalMinutes, nextHour, nextMinute, onTrigger)
}

func (s *Schedule) addJob(name string, intervalMinutes int, nextHour int, nextMinute int, onTrigger TriggerCallback) {
	jobs := s.buckets[nextMinute]
	if jobs == nil {
		jobs = list.New()
		s.buckets[nextMinute] = jobs
	}

	jobs.PushBack(Job{
		Name:            name,
		OnTrigger:       onTrigger,
		IntervalMinutes: intervalMinutes,
		NextHour:        nextHour,
		NextMinute:      nextMinute,
	})
}

func (s *Schedule) Reschedule(job Job) {
	jobs := s.buckets[job.NextMinute]
	if jobs == nil {
		jobs = list.New()
		s.buckets[job.NextMinute] = jobs
	}
	jobs.PushBack(job)
}

func (s *Schedule) GetJobsAt(minute int) *list.List {
	return s.buckets[minute]
}
