package models

import "time"

type Report struct {
	GroupId     int64
	StudentId   int64
	TrainerId   int64
	Description string
	Date        time.Time
}
