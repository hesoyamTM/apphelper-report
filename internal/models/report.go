package models

import (
	"time"

	"github.com/google/uuid"
)

type Report struct {
	GroupId     uuid.UUID
	StudentId   uuid.UUID
	TrainerId   uuid.UUID
	Description string
	Date        time.Time
}
