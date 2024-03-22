package scheduler

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTime_GetHourAndMinute(t *testing.T) {
	var time Time

	time = Time(0)
	assert.Equal(t, time.GetHour(), 0)
	assert.Equal(t, time.GetMinute(), 0)

	time = Time(17)
	assert.Equal(t, time.GetHour(), 0)
	assert.Equal(t, time.GetMinute(), 17)

	time = Time(100)
	assert.Equal(t, time.GetHour(), 1)
	assert.Equal(t, time.GetMinute(), 40)

	time = Time(120)
	assert.Equal(t, time.GetHour(), 2)
	assert.Equal(t, time.GetMinute(), 0)
}

func TestTime_String(t *testing.T) {
	assert.Equal(t, Time(0).String(), "00:00")
	assert.Equal(t, Time(5).String(), "00:05")
	assert.Equal(t, Time(17).String(), "00:17")
	assert.Equal(t, Time(100).String(), "01:40")
	assert.Equal(t, Time(120).String(), "02:00")
}

func TestTime_AddMinutes(t *testing.T) {
	assert.Equal(t, Time(0).AddMinutes(50), Time(50))
	assert.Equal(t, Time(17).AddMinutes(60), Time(77))
}
