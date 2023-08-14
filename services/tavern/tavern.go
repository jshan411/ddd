package tavern

import (
	"log"

	"github.com/google/uuid"
	"github.com/jshan411/ddd/services/order"
)

type TavernConfiguration func(os *Tavern) error

type Tavern struct {
	// orderservice to takes orders
	OrderService   *order.OrderService
	BillingService interface{}
}

func NewTavern(cfgs ...TavernConfiguration) (*Tavern, error) {
	t := &Tavern{}

	for _, cfg := range cfgs {
		if err := cfg(t); err != nil {
			return nil, err
		}
	}

	return t, nil
}

func WithOrderService(os *order.OrderService) TavernConfiguration {
	return func(t *Tavern) error {
		t.OrderService = os
		return nil
	}
}

func (t *Tavern) Order(customerID uuid.UUID, productIDs []uuid.UUID) error {
	price, err := t.OrderService.CreateOrder(customerID, productIDs)
	if err != nil {
		return err
	}

	log.Printf("\nBill the customer: %0.0f\n", price)

	return nil
}
