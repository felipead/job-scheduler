package scheduler

import "fmt"

type Time int

func (t Time) GetMinute() int {
	return int(t % 60)
}

func (t Time) GetHour() int {
	return int(t / 60)
}

func (t Time) String() string {
	return fmt.Sprintf("%02d:%02d", t.GetHour(), t.GetMinute())
}

func (t Time) AddMinutes(minutes int) Time {
	return t + Time(minutes)
}
