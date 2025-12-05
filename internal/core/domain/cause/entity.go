package cause

import "time"

// TODO: when logic is complex, should not return domain object directly
type Cause struct {
	ID     uint   `json:"id"`
	Reason string `json:"reason"`
	Status bool   `json:"status"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updateAt"`
}
