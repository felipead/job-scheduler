package main

import "github.com/felipead/job-scheduler/pkg/scheduler"

func main() {
	// Job 1 is executed only once per hour at minute 17 → (00:17, 01:17, 02:17, …)
	// Job 2 is executed every 4th minute → (00:04, 00:08, 00:12, 00:16 …)
	// Job 3 is executed every 6th minute with a 1-minute offset → (00:07, 00:13, 00:19, …)
	// Job 4 is executed every 25th minute with a 2-minute offset → (00:27, 00:52, 01:17, 01:42)
	// Job 5 is executed every 100th minute → (01:40, 03:20, 5:40, …)

	schedule := scheduler.NewSchedule()

	schedule.AddHourlyJob("Job 1", 17, nil)
	schedule.AddIntervalJob("Job 2", 4, 0, nil)
	schedule.AddIntervalJob("Job 3", 6, 1, nil)
	schedule.AddIntervalJob("Job 4", 25, 2, nil)
	schedule.AddIntervalJob("Job 5", 100, 0, nil)

	scheduler.JobLoop(schedule)
}
