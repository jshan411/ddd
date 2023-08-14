package tavern

import (
	"time"

	"github.com/google/uuid"
)

// Transaction is a valueobject, not entity since the transaction has no identifier and is immutable.
type Transaction struct {
	amount    int
	from      uuid.UUID
	to        uuid.UUID
	createdAt time.Time
}
