package scheduler

const _24Hours = 1440

func JobLoop() {

	//
	// If performance was not a concern, we could simply maintain a simple list or array containing all jobs.
	// For every minute, we would go over that list and find those that matches that exact minute (and hour)
	// However, this can be very slow if we have many jobs and are running on the scale of seconds or milliseconds.
	// Instead, I am going to sort jobs into buckets. The idea is that buckets are sorted by a meaningful time unit.
	// Since our smallest unit is minutes, I will keep the jobs sorted by 60 buckets, each of them corresponding to
	// a minute of the hour.
	//
	// The algorithm works as follows:
	//
	//   let hourlySchedule: a hash table where the minute of the hour (0-59) maps to an array of jobs (hourly)
	//   let intervalSchedule: a hash table where the minute of the hour (0-59) maps to a linked-list of jobs (interval)
	//
	//   for any given time:
	//       hour ← time div 60
	//       minute ← time mod 60
	//
	//       for job in hourlySchedule[minute]:
	//           job.Trigger()
	//
	//       J ← intervalSchedule[minute]
	//       for job in J:
	//       	job.Trigger()
	//          nextMinute ← (time + job.Interval) mod 60
	//          remove job from J
	//          reschedule job for nextMinute
	//
	// That solves the problem for hourly jobs (i.e.: every 17 minutes of the hour) and interval jobs with short
	// intervals (eg: every 25th minute). However, if the interval spans for more than 60 minutes,
	// (i.e: repeat every 100 minutes) that interval wouldn't be respected.
	//
	// We can still keep the jobs sorted by the minute of the hour they are supposed to be triggered. However, before
	// triggering a job, we will check if the hour it is supposed to be triggered also matches the current hour.
	// We only want to trigger jobs if both the hour and the minute match the schedule.
	//
	// The technique used here is called "indexing", where we index each job by the minute of the hour it is supposed
	// to be run. That makes determining if a given job is supposed to be run in a given hour and minute on average:
	//
	//    O(N ÷ 60)
	//
	// where N is the total number of jobs (assuming a uniform distribution).
	//
	// If indexing by minute is not appropriate, we could index by any other unit of time, like hour, day or even
	// second.
	//

	time := 0
	for time < _24Hours {
		hour := time / 60
		minute := time % 60

		runHourlySchedule(time, hour, minute)
		runIntervalSchedule(time, hour, minute)

		time++

		//
		// In a production system, we would have delays or sleeps to avoid busy work.
		// Here, we are only interested in evaluating the scheduler algorithm.
		//
	}
}

func runHourlySchedule(time int, hour int, minute int) {
	if jobs := getHourlyScheduleAt(minute); jobs != nil && len(jobs) > 0 {
		for _, job := range jobs {
			job.Trigger(time, hour, minute)
		}
	}
}

func runIntervalSchedule(time int, hour int, minute int) {
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
				job.Trigger(time, hour, minute)
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
}
