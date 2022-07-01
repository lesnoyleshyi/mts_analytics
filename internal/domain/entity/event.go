package entity

import "time"

type Event struct {
	TaskUUID  string    `json:"task_uuid"`
	EventType string    `json:"event"`
	UserUUID  string    `json:"user_uuid"`
	Timestamp time.Time `json:"timestamp"`
}
