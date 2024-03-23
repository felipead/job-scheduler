# job-scheduler

Simplified Cron-like job scheduler algorithm

The idea came as an interview problem, which I then further refined. Design a job scheduler where it is possible to schedule jobs that:

1. Run **once per hour** at a given minute.
2. Run at a given **interval**, expressed in minutes.
3. Run at a given **interval**, expressed in minutes, with an **offset**.

Here are some examples:

- Job 1 is executed only once per hour at minute 17 → 00:17, 01:17, 02:17, …
- Job 2 is executed every 4th minute → 00:04, 00:08, 00:12, 00:16 …
- Job 3 is executed every 6th minute with a 1-minute offset → 00:07, 00:13, 00:19, …
- Job 4 is executed every 25th minute with a 2-minute offset → 00:27, 00:52, 01:17, 01:42, …
- Job 5 is executed every 100th minute → 01:40, 03:20, 5:40, …

## Interface

In order to schedule jobs, first you must create a `Schedule` instance, then add hourly or interval jobs.

```go
package main

import "scheduler"

func main() {
	schedule := scheduler.NewSchedule()

	schedule.AddHourlyJob("Job 1", 17, nil)
	schedule.AddIntervalJob("Job 2", 4, 0, nil)
	schedule.AddIntervalJob("Job 3", 6, 1, nil)
	schedule.AddIntervalJob("Job 4", 25, 2, nil)
	schedule.AddIntervalJob("Job 5", 100, 0, nil)

	scheduler.JobLoop(schedule) // blocks
}
```

When scheduling is done, the `JobLoop` is invoked. It will block the main thread.

All jobs accept an optional callback function, which is called when the job is triggered:

```go
schedule.AddHourlyJob("Job 1", 17, func(id string, time schedule.Time) {
    fmt.Printf("job %s was triggered at %s!", id, time.String())
})
```

## Algorithm Design

For a small number of jobs, we could simply keep them in an array. For every minute, we would go over the array and trigger those jobs that match that exact minute (and hour):

```
let jobs: an array containing all scheduled jobs

for any given time:
    for job in jobs:
        if job schedule matches time:
            job.Trigger()
```

However, this has `O(n²)` time performance, where `n` is the number jobs. It could be slow if we have many jobs or are running in the scale of seconds or milliseconds.

Instead, we will sort jobs into buckets. The idea is that buckets are sorted by a meaningful time unit. Since our smallest unit is minutes, We will keep the jobs sorted by 60 buckets, each of them corresponding to a minute of the hour. But that could change if indexing by minute is not appropriate for the problem at hand. We could index by any other unit of time, like hour, day or even second.

The algorithm works as follows:

```
let schedule: a hash table where the minute of the hour (0-59) maps to a list of jobs J

for any given time:
    hour ← time div 60
    minute ← time mod 60

    J ← schedule[minute]      // J is a linked-list
    for job in J:
        job.Trigger()
        remove job from J
        nextMinute ← (time + job.Interval) mod 60
        reschedule job for nextMinute
```

That solves the problem for hourly jobs (i.e.: every 17 minutes of the hour) and interval jobs with short  intervals (eg: every 25th minute). However, if the interval spans for more than 60 minutes, (i.e: repeat every 100 minutes) that interval wouldn't be respected.

We can still keep the jobs sorted by the minute of the hour they are supposed to be triggered. However, before triggering a job, we will check if the hour it is supposed to be triggered also matches the current hour. We only want to trigger jobs if both the hour and the minute match the schedule. The algorithm would then slightly change to:

```
J ← schedule[minute]
for job in J:
    if job.NextHour = hour:
        job.Trigger()
        remove job from J
        t ← (time + job.Interval)
        nextMinute ← t mod 60
        nextHour ← t div 60
        reschedule job for nextMinute and nextHour
```

The technique used here is called **indexing**, **[bucket sort](https://en.wikipedia.org/wiki/Bucket_sort)**, or **bin sort**, where we index each job by the minute of the  hour it is supposed to be run (0-59). That makes determining if a given job is supposed to be run in a given hour and minute on average:

    O(n ÷ 60)

where `n` is the total number of jobs (assuming a uniform distribution). For a small and uniform enough set, that can be close to `O(1)`.

Unlike bucket sort though, we don't care about sorting the jobs inside each bucket. That could change though if we determine that jobs can have **priorities**. In that case, we could use a **[priority queue](https://en.wikipedia.org/wiki/Priority_queue)**.

## Roadmap

This is a very simplistic and limited implementation. However, the following are planned:

- [ ] The `JobLoop` is not blocking, neither waits for the scheduled amount of time. It was coded purely for demonstrating the algorithm. That must change so that we can consider it production ready.
- [ ] Allow jobs to be scheduled while the loop is running in a separate routine. Currently, this is not possible and the scheduler is not thread-safe.
- [ ] Currently, intervals are specified in minutes and jobs are also sorted in buckets of minutes (see the algorithm design above). Ideally, this could be configurable. For some applications, it might make sense to run the jobs in intervals of seconds, hours or even days. The bucket time unit should be adjusted accordingly.
- [ ] Currently, all jobs are treated with the same priority. However, we could determine that some jobs are have higher priority than others, then keep them sorted using a **priority queue**.
- [ ] Use standard logging interface instead of `fmt.Println`.
- [ ] Use Go standard `time.Duration` for specifying time intervals.
- [ ] Support one-off or adhoc jobs, which will run just once at a specific time.
- [ ] Some jobs can be slow, which will impact the overall scheduling. Run the jobs in separate go-routines, and abort them if they exceed the timeout.
- [ ] Ability to retry failed jobs up to a specified number of times, in the next minute or so.

## Development

Compiling:
```shell
make build
```

Running tests:
```shell
make test
```

Running a quick demo:
```shell
./run-demo.sh
```
