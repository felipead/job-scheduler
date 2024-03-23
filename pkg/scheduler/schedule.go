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

func (s *Schedule) AddHourlyJob(id string, minute int, onTrigger TriggerCallback) {
	nextTime := Time(minute % 60)
	s.addJob(id, 60, nextTime, onTrigger)
}

func (s *Schedule) AddIntervalJob(id string, intervalMinutes int, offset int, onTrigger TriggerCallback) {
	nextTime := Time(offset + intervalMinutes)
	s.addJob(id, intervalMinutes, nextTime, onTrigger)
}

func (s *Schedule) addJob(id string, intervalMinutes int, nextTime Time, onTrigger TriggerCallback) {
	nextMinute := nextTime.GetMinute()

	jobs := s.buckets[nextMinute]
	if jobs == nil {
		jobs = list.New()
		s.buckets[nextMinute] = jobs
	}

	jobs.PushBack(Job{
		ID:              id,
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

func (s *Schedule) GetJobsAt(time Time) *list.List {
	minute := time.GetMinute()
	return s.buckets[minute]
}
