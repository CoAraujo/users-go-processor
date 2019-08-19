package domains

import "time"

type QueueResponse struct {
	Id         string     `json:"id"`
	EnqueuedAt *time.Time `json:"enqueuedAt"`
	ClientID   string     `json:"clientId, omitempty"`
}
