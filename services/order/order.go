package order

import (
	"log"

	"github.com/google/uuid"
	"github.com/jshan411/ddd/domain/customer"
	"github.com/jshan411/ddd/domain/customer/memory"
	"github.com/jshan411/ddd/domain/product"
	productMemory "github.com/jshan411/ddd/domain/product/memory"
)

type OrderConfiguration func(os *OrderService) error

type OrderService struct {
	customers customer.Repository
	products  product.Repository

	// billing billing.Service
}

func NewOrderService(cfgs ...OrderConfiguration) (*OrderService, error) {
	os := &OrderService{}
	// Loop through all the cfgs and apply them
	for _, cfg := range cfgs {
		err := cfg(os)

		if err != nil {
			return nil, err
		}
	}
	return os, nil
}

func WithCustomerRepository(cr customer.Repository) OrderConfiguration {
	// Return a function that matches the OrderConfiguration alias
	return func(os *OrderService) error {
		os.customers = cr
		return nil
	}
}

func WithMemoryCustomerRepository() OrderConfiguration {
	cr := memory.New()
	return WithCustomerRepository(cr)
}

// var server = "neoheliostest.database.windows.net"
// var port = 1433
// var user = "panstar_admin"
// var password = "Vostmxk0712!"
// var db_name = "sql-line-test-gorm"

// func WithSqlserverCustomerRepository(ctx context.Context, connectionString string) OrderConfiguration {
// 	return func(os *OrderService) error {
// 		cr, err := customer.sqlserver.New("host")
// 		if err != nil {
// 			return err
// 		}
// 		os.customers = cr
// 		return nil
// 	}
// }

func WithMemoryProductRepository(products []product.Product) OrderConfiguration {
	return func(os *OrderService) error {
		pr := productMemory.New()

		for _, p := range products {
			if err := pr.Add(p); err != nil {
				return err
			}
		}

		os.products = pr
		return nil
	}
}

// func WithSqlserverProductRepository(products []product.Product) OrderConfiguration {
// 	return func(os *OrderService) error {
// 		pr := productMemory.New()

// 		for _, p := range products {
// 			if err := pr.Add(p); err != nil {
// 				return err
// 			}
// 		}

// 		os.products = pr
// 		return nil
// 	}
// }

func (o *OrderService) CreateOrder(customerID uuid.UUID, productsIDs []uuid.UUID) (float64, error) {
	// Fetch the customer
	c, err := o.customers.Get(customerID)
	if err != nil {
		return 0, err
	}

	var products []product.Product
	var totalPrice float64

	for _, id := range productsIDs {
		p, err := o.products.GetById(id)

		if err != nil {
			return 0, err
		}

		products = append(products, p)
		totalPrice += p.GetPrice()
	}

	log.Printf("Customer: %s has ordered %d products in total price %0.0f.", c.GetID(), len(products), totalPrice)
	return totalPrice, nil
}

func (o *OrderService) AddCustomer(name string) (uuid.UUID, error) {
	c, err := customer.NewCustomer(name)
	if err != nil {
		return uuid.Nil, err
	}
	err = o.customers.Add(c)
	if err != nil {
		return uuid.Nil, err
	}

	return c.GetID(), nil
}
