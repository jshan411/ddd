package tavern

import (
	"github.com/google/uuid"
)

// Item is an entity that represents an item in all domains
type Item struct {
	// ID is the identifier of the entity
	ID          uuid.UUID `gorm:"type:uniqueidentifier; primary_key;"`
	Name        string    `gorm:"size:255"`
	Description string    `gorm:"size:255"`
}
