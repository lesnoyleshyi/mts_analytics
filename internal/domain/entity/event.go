package entity

import "time"

//type EventType int32
//
//const (
//	EVENT_UNKNOWN = iota
//	EVENT_CREATED
//	EVENT_SENT_TO
//	EVENT_APPROVED_BY
//	EVENT_REJECTED_BY
//	EVENT_SIGNED
//	EVENT_SENT
//)

//easyjson:json
type Event struct {
	TaskUUID  string    `json:"task_uuid"`
	EventType string    `json:"event"`
	UserUUID  string    `json:"user_uuid"`
	Timestamp time.Time `json:"timestamp"`
}
