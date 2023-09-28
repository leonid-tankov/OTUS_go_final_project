package kafka

import "time"

type eventType string

const (
	Click eventType = "click"
	Show  eventType = "show"
)

type Message struct {
	EventType        eventType `json:"eventType"`
	BannerID         int64     `json:"bannerId"`
	SlotID           int64     `json:"slotId"`
	SocialDemGroupID int64     `json:"socialDemGroupId"`
	Timestamp        time.Time `json:"timestamp"`
}
