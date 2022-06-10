package domain

import "time"

type Event struct {
	TaskUUID  string    `json:"task_uuid"`
	EventType string    `json:"event"`
	UserUUID  string    `json:"user_uuid"`
	Timestamp time.Time `json:"timestamp"`
}

type SignedCount int

func (s SignedCount) Match() {}

type NotSignedYetCount int

func (ns NotSignedYetCount) Match() {}

type SignationTotalTime struct {
	TaskUUID string
}
