package order

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jshan411/ddd/domain/customer"
	"github.com/jshan411/ddd/domain/product"
)

func init_products(t *testing.T) []product.Product {
	yebisu, err := product.NewProduct("YEBISU", "Beer", 4000)
	if err != nil {
		t.Fatal(err)
	}

	denroku, err := product.NewProduct("Denroku", "Snack", 1500)
	if err != nil {
		t.Fatal(err)
	}

	evaFigure, err := product.NewProduct("EvaFigure", "Figure", 100000)
	if err != nil {
		t.Fatal(err)
	}

	return []product.Product{
		yebisu, denroku, evaFigure,
	}
}

func TestOrder_NewOrderService(t *testing.T) {
	// 함수 이름이 Test로 시작하면, test 기능 추가됨
	products := init_products(t)

	os, err := NewOrderService(
		WithMemoryCustomerRepository(),
		WithMemoryProductRepository(products),
		// WithSqlserverCustomerRepository(),
		// WithSqlserverProductRepository(products),
	)

	if err != nil {
		t.Fatal(err)
	}

	cust, err := customer.NewCustomer("jshan")

	if err != nil {
		t.Fatal(err)
	}

	err = os.customers.Add(cust)
	if err != nil {
		t.Error(err)
	}

	order := []uuid.UUID{
		products[0].GetID(),
	}

	_, err = os.CreateOrder(cust.GetID(), order)

	if err != nil {
		t.Error(err)
	}
}
