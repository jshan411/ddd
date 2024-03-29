// Package memory is an in-memory implementation of customer repository
package memory

import (
	"fmt"
	"sync"

	"github.com/google/uuid"
	"github.com/jshan411/ddd/domain/customer"
)

type MemoryRepository struct {
	customers map[uuid.UUID]customer.Customer
	sync.Mutex
}

func New() *MemoryRepository {
	return &MemoryRepository{
		customers: make(map[uuid.UUID]customer.Customer),
	}
}

func (mr *MemoryRepository) Get(id uuid.UUID) (customer.Customer, error) {
	if customer, ok := mr.customers[id]; ok {
		return customer, nil
	}
	return customer.Customer{}, customer.ErrCustomerNotFound
}

func (mr *MemoryRepository) Add(c customer.Customer) error {
	if mr.customers == nil {
		// lock the Mutex to protect the data. See the go routine.
		mr.Lock()
		mr.customers = make(map[uuid.UUID]customer.Customer)
		mr.Unlock()
	}
	// Make sure that customer is already in repository
	if _, ok := mr.customers[c.GetID()]; ok {
		return fmt.Errorf("custoemr already exists :%w", customer.ErrFailedToAddCustomer)
	}
	mr.Lock()
	mr.customers[c.GetID()] = c
	mr.Unlock()
	return nil
}

func (mr *MemoryRepository) Update(c customer.Customer) error {
	if _, ok := mr.customers[c.GetID()]; !ok {
		return fmt.Errorf("customer does not exist: %w", customer.ErrUpdateCustomer)
	}
	mr.Lock()
	mr.customers[c.GetID()] = c
	mr.Unlock()
	return nil
}
