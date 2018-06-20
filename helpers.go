package thingscloud

import "time"

// String returns a pointer to a string
func String(str string) *string {
	return &str
}

// Status returns a pointer to a TaskStatus
func Status(val TaskStatus) *TaskStatus {
	return &val
}

// Schedule returns a pointer to a TaskSchedule
func Schedule(val TaskSchedule) *TaskSchedule {
	return &val
}

// Time returns a pointer to a Time
func Time(val time.Time) *Timestamp {
	ts := Timestamp(val)
	return &ts
}
