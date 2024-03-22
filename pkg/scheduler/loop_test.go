package scheduler

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJobLoop_RunSchedule_HourlyJob(t *testing.T) {
	schedule := NewSchedule()

	callTimes := make([]int, 0)
	schedule.AddHourlyJob("foobar", 17, func(_ string, time int, _ int, _ int) {
		callTimes = append(callTimes, time)
	})

	for time := 0; time < 180; time++ {
		RunSchedule(schedule, time)
	}

	assert.Equal(t, callTimes, []int{
		17,
		17 + 60,
		17 + 60*2,
	})
}

func TestJobLoop_RunSchedule_FewHourlyJobs(t *testing.T) {
	schedule := NewSchedule()

	callTimes1 := make([]int, 0)
	schedule.AddHourlyJob("one", 17, func(name string, time int, hour int, minute int) {
		assert.Equal(t, name, "one")
		callTimes1 = append(callTimes1, time)
	})

	callTimes2 := make([]int, 0)
	schedule.AddHourlyJob("two", 59, func(name string, time int, hour int, minute int) {
		assert.Equal(t, name, "two")
		callTimes2 = append(callTimes2, time)
	})

	callTimes3 := make([]int, 0)
	schedule.AddHourlyJob("three", 17, func(name string, time int, hour int, minute int) {
		assert.Equal(t, name, "three")
		callTimes3 = append(callTimes3, time)
	})

	callTimes4 := make([]int, 0)
	schedule.AddHourlyJob("four", 0, func(name string, time int, hour int, minute int) {
		assert.Equal(t, name, "four")
		callTimes4 = append(callTimes4, time)
	})

	for time := 0; time < 180; time++ {
		RunSchedule(schedule, time)
	}

	assert.Equal(t, callTimes1, []int{
		17,
		17 + 60,
		17 + 60*2,
	})

	assert.Equal(t, callTimes2, []int{
		59,
		59 + 60,
		59 + 60*2,
	})

	assert.Equal(t, callTimes3, []int{
		17,
		17 + 60,
		17 + 60*2,
	})

	assert.Equal(t, callTimes4, []int{
		0,
		60,
		60 * 2,
	})
}

func TestJobLoop_RunSchedule_IntervalJob(t *testing.T) {
	schedule := NewSchedule()

	interval := 17

	callTimes := make([]int, 0)
	schedule.AddIntervalJob("foobar", interval, 0, func(_ string, time int, _ int, _ int) {
		callTimes = append(callTimes, time)
	})

	for time := 0; time < 180; time++ {
		RunSchedule(schedule, time)
	}

	assert.Equal(t, callTimes, []int{
		interval,
		interval * 2,
		interval * 3,
		interval * 4,
		interval * 5,
		interval * 6,
		interval * 7,
		interval * 8,
		interval * 9,
		interval * 10,
	})
}

func TestJobLoop_RunSchedule_IntervalJobWithOffset(t *testing.T) {
	schedule := NewSchedule()

	interval := 17
	offset := 10

	callTimes := make([]int, 0)
	schedule.AddIntervalJob("foobar", interval, offset, func(_ string, time int, _ int, _ int) {
		callTimes = append(callTimes, time)
	})

	for time := 0; time < 180; time++ {
		RunSchedule(schedule, time)
	}

	assert.Equal(t, callTimes, []int{
		interval + offset,
		(interval * 2) + offset,
		(interval * 3) + offset,
		(interval * 4) + offset,
		(interval * 5) + offset,
		(interval * 6) + offset,
		(interval * 7) + offset,
		(interval * 8) + offset,
		(interval * 9) + offset,
	})
}

func TestJobLoop_RunSchedule_IntervalJobGreaterThan1Hour(t *testing.T) {
	schedule := NewSchedule()

	interval := 100
	offset := 15

	callTimes := make([]int, 0)
	schedule.AddIntervalJob("foobar", interval, offset, func(_ string, time int, _ int, _ int) {
		callTimes = append(callTimes, time)
	})

	for time := 0; time < 360; time++ {
		RunSchedule(schedule, time)
	}

	assert.Equal(t, callTimes, []int{
		interval + offset,
		(interval * 2) + offset,
		(interval * 3) + offset,
	})
}

func TestJobLoop_RunSchedule_FewIntervalJobs(t *testing.T) {
	schedule := NewSchedule()

	callTimes1 := make([]int, 0)
	schedule.AddIntervalJob("one", 17, 0, func(name string, time int, _ int, _ int) {
		assert.Equal(t, name, "one")
		callTimes1 = append(callTimes1, time)
	})

	callTimes2 := make([]int, 0)
	schedule.AddIntervalJob("two", 15, 2, func(name string, time int, _ int, _ int) {
		assert.Equal(t, name, "two")
		callTimes2 = append(callTimes2, time)
	})

	callTimes3 := make([]int, 0)
	schedule.AddIntervalJob("three", 45, 6, func(name string, time int, _ int, _ int) {
		assert.Equal(t, name, "three")
		callTimes3 = append(callTimes3, time)
	})

	callTimes4 := make([]int, 0)
	schedule.AddIntervalJob("four", 65, 3, func(name string, time int, _ int, _ int) {
		assert.Equal(t, name, "four")
		callTimes4 = append(callTimes4, time)
	})

	callTimes5 := make([]int, 0)
	schedule.AddIntervalJob("five", 17, 0, func(name string, time int, _ int, _ int) {
		assert.Equal(t, name, "five")
		callTimes5 = append(callTimes5, time)
	})

	for time := 0; time < 180; time++ {
		RunSchedule(schedule, time)
	}

	assert.Equal(t, callTimes1, []int{
		17,
		17 * 2,
		17 * 3,
		17 * 4,
		17 * 5,
		17 * 6,
		17 * 7,
		17 * 8,
		17 * 9,
		17 * 10,
	})

	assert.Equal(t, callTimes2, []int{
		15 + 2,
		15*2 + 2,
		15*3 + 2,
		15*4 + 2,
		15*5 + 2,
		15*6 + 2,
		15*7 + 2,
		15*8 + 2,
		15*9 + 2,
		15*10 + 2,
		15*11 + 2,
	})

	assert.Equal(t, callTimes3, []int{
		45 + 6,
		45*2 + 6,
		45*3 + 6,
	})

	assert.Equal(t, callTimes4, []int{
		65 + 3,
		65*2 + 3,
	})

	assert.Equal(t, callTimes5, []int{
		17,
		17 * 2,
		17 * 3,
		17 * 4,
		17 * 5,
		17 * 6,
		17 * 7,
		17 * 8,
		17 * 9,
		17 * 10,
	})
}

func TestJobLoop_RunSchedule_MixedHourlyAndIntervalJobs(t *testing.T) {
	schedule := NewSchedule()

	callTimes1 := make([]int, 0)
	schedule.AddIntervalJob("one", 17, 0, func(name string, time int, _ int, _ int) {
		assert.Equal(t, name, "one")
		callTimes1 = append(callTimes1, time)
	})

	callTimes2 := make([]int, 0)
	schedule.AddIntervalJob("two", 15, 2, func(name string, time int, _ int, _ int) {
		assert.Equal(t, name, "two")
		callTimes2 = append(callTimes2, time)
	})

	callTimes3 := make([]int, 0)
	schedule.AddHourlyJob("three", 17, func(name string, time int, _ int, _ int) {
		assert.Equal(t, name, "three")
		callTimes3 = append(callTimes3, time)
	})

	callTimes4 := make([]int, 0)
	schedule.AddIntervalJob("four", 65, 3, func(name string, time int, _ int, _ int) {
		assert.Equal(t, name, "four")
		callTimes4 = append(callTimes4, time)
	})

	callTimes5 := make([]int, 0)
	schedule.AddHourlyJob("five", 25, func(name string, time int, _ int, _ int) {
		assert.Equal(t, name, "five")
		callTimes5 = append(callTimes5, time)
	})

	callTimes6 := make([]int, 0)
	schedule.AddIntervalJob("six", 45, 6, func(name string, time int, _ int, _ int) {
		assert.Equal(t, name, "six")
		callTimes6 = append(callTimes6, time)
	})

	callTimes7 := make([]int, 0)
	schedule.AddIntervalJob("seven", 17, 0, func(name string, time int, _ int, _ int) {
		assert.Equal(t, name, "seven")
		callTimes7 = append(callTimes7, time)
	})

	for time := 0; time < 180; time++ {
		RunSchedule(schedule, time)
	}

	assert.Equal(t, callTimes1, []int{
		17,
		17 * 2,
		17 * 3,
		17 * 4,
		17 * 5,
		17 * 6,
		17 * 7,
		17 * 8,
		17 * 9,
		17 * 10,
	})

	assert.Equal(t, callTimes2, []int{
		15 + 2,
		15*2 + 2,
		15*3 + 2,
		15*4 + 2,
		15*5 + 2,
		15*6 + 2,
		15*7 + 2,
		15*8 + 2,
		15*9 + 2,
		15*10 + 2,
		15*11 + 2,
	})

	assert.Equal(t, callTimes3, []int{
		17,
		17 + 60,
		17 + 60*2,
	})

	assert.Equal(t, callTimes4, []int{
		65 + 3,
		65*2 + 3,
	})

	assert.Equal(t, callTimes5, []int{
		25,
		25 + 60,
		25 + 60*2,
	})

	assert.Equal(t, callTimes6, []int{
		45 + 6,
		45*2 + 6,
		45*3 + 6,
	})

	assert.Equal(t, callTimes7, []int{
		17,
		17 * 2,
		17 * 3,
		17 * 4,
		17 * 5,
		17 * 6,
		17 * 7,
		17 * 8,
		17 * 9,
		17 * 10,
	})
}
