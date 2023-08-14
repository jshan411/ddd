package memory

import (
	"sync"

	"github.com/google/uuid"
	"github.com/jshan411/ddd/domain/product"
)

type MemoryProductRepository struct {
	products map[uuid.UUID]product.Product
	sync.Mutex
}

func New() *MemoryProductRepository {
	return &MemoryProductRepository{
		products: make(map[uuid.UUID]product.Product),
	}
}

func (mpr *MemoryProductRepository) GetAll() ([]product.Product, error) {
	var products []product.Product

	for _, product := range mpr.products {
		products = append(products, product)
	}
	return products, nil
}

func (mpr *MemoryProductRepository) GetById(id uuid.UUID) (product.Product, error) {
	if product, ok := mpr.products[id]; ok {
		return product, nil
	}
	return product.Product{}, product.ErrProductNotFound
}

func (mpr *MemoryProductRepository) Add(newProduct product.Product) error {
	mpr.Lock()
	defer mpr.Unlock()

	if _, ok := mpr.products[newProduct.GetID()]; ok {
		return product.ErrProductAlreadyExist
	}

	mpr.products[newProduct.GetID()] = newProduct

	return nil
}

func (mpr *MemoryProductRepository) Update(update product.Product) error {
	// 여기서, Update(product product.Product)가 아니라 Update(update product.Product)라고 선언한 이유는,
	// 밑에 product.ErrProductNotFound 이 부분에서 product라는 package를 호출해야 하기 때문.
	mpr.Lock()
	defer mpr.Unlock()

	if _, ok := mpr.products[update.GetID()]; !ok {
		return product.ErrProductNotFound
	}
	mpr.products[update.GetID()] = update
	return nil
}

func (mpr *MemoryProductRepository) Delete(id uuid.UUID) error {
	mpr.Lock()
	defer mpr.Unlock()

	if _, ok := mpr.products[id]; !ok {
		return product.ErrProductNotFound
	}
	delete(mpr.products, id)

	return nil
}
