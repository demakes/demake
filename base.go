package klaro

import (
	"time"
)

type Base struct {
	Backend   Backend
	ID        uint64    `json:"-"`
	ExtID     []byte    `json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
