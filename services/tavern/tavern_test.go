package tavern

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jshan411/ddd/domain/product"
	"github.com/jshan411/ddd/services/order"
)

func Test_Tavern(t *testing.T) {
	products := init_products(t)

	os, err := order.NewOrderService(
		order.WithMemoryCustomerRepository(), // 이 부분을 내 DB로 바꾸면 된다.
		order.WithMemoryProductRepository(products),
		// order.WithSqlserverCustomerRepository(context.Background(), "sqlhost쓰자"),
		// order.WithSqlserverProductRepository(products),
	)

	if err != nil {
		t.Fatal(err)
	}

	tavern, err := NewTavern(WithOrderService(os))

	if err != nil {
		t.Fatal(err)
	}

	uid, err := os.AddCustomer("jshan")
	if err != nil {
		t.Fatal(err)
	}

	order := []uuid.UUID{
		products[0].GetID(),
	}

	err = tavern.Order(uid, order)
	if err != nil {
		t.Fatal(err)
	}
}

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
