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
	nextTime := Time(minute % 60)
	s.addJob(name, 60, nextTime, onTrigger)
}

func (s *Schedule) AddIntervalJob(name string, intervalMinutes int, offset int, onTrigger TriggerCallback) {
	nextTime := Time(offset + intervalMinutes)
	s.addJob(name, intervalMinutes, nextTime, onTrigger)
}

func (s *Schedule) addJob(name string, intervalMinutes int, nextTime Time, onTrigger TriggerCallback) {
	nextMinute := nextTime.GetMinute()

	jobs := s.buckets[nextMinute]
	if jobs == nil {
		jobs = list.New()
		s.buckets[nextMinute] = jobs
	}

	jobs.PushBack(Job{
		Name:            name,
		OnTrigger:       onTrigger,
		IntervalMinutes: intervalMinutes,
		NextTime:        nextTime,
	})
}

func (s *Schedule) Reschedule(job Job) {
	nextMinute := job.NextTime.GetMinute()

	jobs := s.buckets[nextMinute]
	if jobs == nil {
		jobs = list.New()
		s.buckets[nextMinute] = jobs
	}
	jobs.PushBack(job)
}

func (s *Schedule) GetJobsAt(minute int) *list.List {
	return s.buckets[minute]
}
