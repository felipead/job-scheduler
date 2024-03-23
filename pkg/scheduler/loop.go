package scheduler

const _24Hours = Time(1440)

func JobLoop(schedule *Schedule) {
	var time Time = 0

	for time < _24Hours {
		RunSchedule(schedule, time)

		time++

		//
		// In a production system, we would have delays or sleeps to avoid busy work.
		// Here, we are only interested in evaluating the scheduler algorithm.
		//
	}
}

func RunSchedule(schedule *Schedule, time Time) {
	jobs := schedule.GetJobsAt(time)

	if jobs == nil || jobs.Len() == 0 {
		return
	}

	triggeredJobs := make([]Job, 0, jobs.Len())

	// iterates over the Linked List, triggering jobs and removing them from the list if they were triggered
	pointer := jobs.Front()
	for pointer != nil {
		job := pointer.Value.(Job)
		nextPointer := pointer.Next()

		// We are indexing jobs by the minute of hour, but we only want to trigger this job
		// if the hour also matches. This is important because we might have intervals that
		// are longer than 60 minutes. For example, repeats every 100 minutes, which will span
		// across 2 hours.
		if job.NextTime == time {
			job.Trigger(time)
			triggeredJobs = append(triggeredJobs, job)
			jobs.Remove(pointer)
		}

		pointer = nextPointer
	}

	// reschedule the jobs that have been triggered for the next interval
	for _, job := range triggeredJobs {
		job.NextTime = time.AddMinutes(job.IntervalMinutes)
		schedule.Reschedule(job)
	}
}
