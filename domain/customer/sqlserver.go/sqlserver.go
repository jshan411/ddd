package customersql

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jshan411/ddd/domain/customer"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

type SQLCustomer struct {
	ID   uuid.UUID `gorm:"type:char(36);primaryKey"`
	Name string    `gorm:"size:255"`
}

func NewFromCustomer(c customer.Customer) SQLCustomer {
	return SQLCustomer{
		ID:   c.GetID(),
		Name: c.GetName(),
	}
}

func (m SQLCustomer) ToAggregate() customer.Customer {
	c := customer.Customer{}
	c.SetID(m.ID)
	c.SetName(m.Name)
	return c
}

type SQLRepository struct {
	db *gorm.DB
}

func New(connectionString string) (*SQLRepository, error) {
	db, err := gorm.Open(sqlserver.Open(connectionString), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// AutoMigrate will create the 'customers' table if it doesn't exist
	err = db.AutoMigrate(&SQLCustomer{})
	if err != nil {
		return nil, err
	}

	return &SQLRepository{
		db: db,
	}, nil
}

func (r *SQLRepository) Get(id uuid.UUID) (customer.Customer, error) {
	var c SQLCustomer
	if err := r.db.First(&c, "id = ?", id).Error; err != nil {
		return customer.Customer{}, err
	}
	return c.ToAggregate(), nil
}

func (r *SQLRepository) Add(c customer.Customer) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	internal := NewFromCustomer(c)

	if err := r.db.WithContext(ctx).Create(&internal).Error; err != nil {
		return err
	}
	return nil
}
