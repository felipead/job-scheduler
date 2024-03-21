package main

import "github.com/felipead/job-scheduler/pkg/scheduler"

func main() {
	// Job 1 is executed only once per hour at minute 17 (HH:17)
	// Job 2 is executed every 4th minute. (HH:04, HH:08, HH:12, etc)
	// Job 3 is executed every 6th minute with a 1-minute offset. (HH:07, HH:13,...)
	// Job 4 is executed every 25th minute (HH:25, HH:50, HH:15, HH:40)
	// Job 5 is executed every 100th minute (HH:25, HH:50, HH:15, HH:40)

	scheduler.ScheduleHourlyJob("Job 1", 17, nil)
	scheduler.ScheduleIntervalJob("Job 2", 4, 0, nil)
	scheduler.ScheduleIntervalJob("Job 3", 6, 1, nil)
	scheduler.ScheduleIntervalJob("Job 4", 25, 0, nil)
	scheduler.JobLoop()
}
