package main

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jshan411/ddd/domain/product"
	tavernDomain "github.com/jshan411/ddd/domain/tavern"
	"github.com/jshan411/ddd/services/order"
	tavernService "github.com/jshan411/ddd/services/tavern"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var (
	server          = "neoheliostest.database.windows.net"
	port            = 1433
	user            = "panstar_admin"
	password        = "Vostmxk0712!"
	db_name         = "fullstack_test"
	errDBConnection = errors.New("fail to connect DB")
)

type BaseModel struct {
	ID        uuid.UUID  `gorm:"type:uniqueidentifier;primary_key;" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"update_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
}

func (base *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	if base.ID == uuid.Nil {
		base.ID = uuid.New()
	}
	return
}

type Pproducts struct {
	BaseModel
	ItemID   uuid.UUID          `gorm:"type:uniqueidentifier;index" json:"item_id"`
	Item     *tavernDomain.Item `gorm:"foreignKey:ItemID;references:ID" json:"item"`
	Price    float64            `json:"price"`
	Quantity int                `json:"quantity"`
}

func main() {
	products := productInventory()

	os, err := order.NewOrderService(
		order.WithMemoryCustomerRepository(),
		order.WithMemoryProductRepository(products),
	)

	fmt.Println("os: ", os)

	if err != nil {
		panic(err)
	}

	tavern, err := tavernService.NewTavern(
		tavernService.WithOrderService(os),
	)

	if err != nil {
		panic(err)
	}

	uid, err := os.AddCustomer("jshan")

	if err != nil {
		panic(err)
	}

	order := []uuid.UUID{
		products[0].GetID(),
		products[0].GetID(),
		products[1].GetID(),
	}

	err = tavern.Order(uid, order)

	if err != nil {
		panic(err)
	}

	//Initialize database, is currently connected to devTest database server
	Instance, dbError := gorm.Open(sqlserver.Open(fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s", user, password, server, port, db_name)), &gorm.Config{})
	if dbError != nil {
		fmt.Println(errDBConnection)
	}

	if err := Instance.AutoMigrate(&tavernDomain.Item{}); err != nil {
		log.Fatalf("Failed to migrate table for Item: %v", err)
	}

	if err := Instance.AutoMigrate(&Pproducts{}); err != nil {
		log.Fatalf("Failed to migrate table for Pproducts: %v", err)
	}

	// Create an item
	item := tavernDomain.Item{ID: uuid.New(), Name: "Test Item", Description: "B sample item for demonstration"}
	Instance.Create(&item)

	fmt.Println("item: ", item)

	// Create a product linked to the above item
	pproduct := Pproducts{ItemID: item.ID, Price: 9.99, Quantity: 10}
	Instance.Create(&pproduct)

	fmt.Println("pproduct: ", pproduct)

	// golang에서 uuid를 sql server에 guid로 저장할 때는 데이터 변환이 일어나지 않으며, golang에서와 같은 데이터가 저장됨.
	// sql server에 저장된 guid를 golang의 uuid로 가져오면 데이터 변환이 일어남
	// 따라서, sql server에서 id를 불러온 뒤에, convertGUIDToUUID함수를 통해 원래의 uuid로 복원 가능

	// Retrieve the product from the database with the associated Item preloaded
	var retrievedProduct Pproducts
	Instance.First(&retrievedProduct, pproduct.ID)
	retrievedProduct.ItemID = convertGUIDToUUID(retrievedProduct.ItemID)
	Instance.Model(&retrievedProduct).Association("Item").Find(&retrievedProduct.Item)

	fmt.Println("retrievedProduct: ", retrievedProduct)
	fmt.Println("retrievedProduct.Item: ", retrievedProduct.Item)

	if retrievedProduct.Item != nil {
		fmt.Println("Item: ", *retrievedProduct.Item)
	} else {
		fmt.Println("Item is nil")
	}

	// //Initialize Router

	// router := initRouter()
	// /*for a security reason apis from random urls are blocked,
	//   below is allowing function from all routes, to be replaced for specific host(url)*/
	// router.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"*"},
	// 	AllowMethods:     []string{"GET", "PUT", "PATCH", "POST", "DELETE", "OPTIONS"},
	// 	AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: true,
	// 	AllowOriginFunc: func(origin string) bool {
	// 		return origin == "https://github.com"
	// 	},
	// 	MaxAge: 12 * time.Hour,
	// }))
	// router.Run(":8080")
}

func initRouter() *gin.Engine {
	router := gin.Default()
	api := router.Group("api")
	{
		// //cargo related apis
		// api.GET("/cargo/agent/location/:user_id", controllers.AgentLocation)
	}
	fmt.Println(api)
	return router
}

func productInventory() []product.Product {
	yebisu, err := product.NewProduct("YEBISU", "Beer", 4000)
	if err != nil {
		panic(err)
	}

	denroku, err := product.NewProduct("Denroku", "Snack", 1500)
	if err != nil {
		panic(err)
	}

	evaFigure, err := product.NewProduct("EvaFigure", "Figure", 100000)
	if err != nil {
		panic(err)
	}

	return []product.Product{
		yebisu, denroku, evaFigure,
	}
}

func convertGUIDToUUID(u uuid.UUID) uuid.UUID {
	return uuid.UUID{
		u[3], u[2], u[1], u[0], // swap byte order for first segment
		u[5], u[4], // swap byte order for second segment
		u[7], u[6], // swap byte order for third segment
		u[8], u[9], u[10], u[11], u[12], u[13], u[14], u[15], // remaining segments stay the same
	}
}
