package scheduler

func JobLoop() {
	time := 0
	for time < 1440 {
		minute := time % 60
		hour := time / 60

		if jobs := getHourlyScheduleAt(minute); jobs != nil && len(jobs) > 0 {
			for _, job := range jobs {
				job.Trigger(hour, minute)
			}
		}

		if jobs := getIntervalScheduleAt(minute); jobs != nil && jobs.Len() > 0 {
			triggeredJobs := make([]IntervalJob, 0, jobs.Len())

			// iterates over the Linked List, triggering jobs and removing them from the list if they were triggered
			pointer := jobs.Front()
			for pointer != nil {
				job := pointer.Value.(IntervalJob)
				nextPointer := pointer.Next()

				// We are indexing jobs by the minute of hour, but we only want to trigger this job
				// if the hour also matches. This is important because we might have intervals that
				// are longer than 60 minutes. For example, repeats every 100 minutes, which will span
				// across 2 hours.
				if job.NextHour == hour {
					job.Trigger(hour, minute)
					triggeredJobs = append(triggeredJobs, job)
					jobs.Remove(pointer)
				}

				pointer = nextPointer
			}

			for _, job := range triggeredJobs {
				nextTime := time + job.IntervalMinutes
				job.NextMinute = nextTime % 60
				job.NextHour = nextTime / 60
				reschedule(job)
			}
		}

		time++

		//
		// In a production system, we would have delays or sleeps to avoid busy work.
		// Here, we are only interested in evaluating the scheduler algorithm.
		//
	}
}
