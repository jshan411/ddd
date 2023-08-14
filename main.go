package main

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jshan411/ddd/domain/product"
	"github.com/jshan411/ddd/services/order"
	tavernService "github.com/jshan411/ddd/services/tavern"
)

var (
	server          = "neoheliostest.database.windows.net"
	port            = 1433
	user            = "panstar_admin"
	password        = "Vostmxk0712!"
	db_name         = "fullstack_test"
	errDBConnection = errors.New("fail to connect DB")
)

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
		products[1].GetID(),
	}

	err = tavern.Order(uid, order)

	if err != nil {
		panic(err)
	}

	// //Initialize database, is currently connected to devTest database server
	// Instance, dbError := gorm.Open(sqlserver.Open(fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s", user, password, server, port, db_name)), &gorm.Config{})
	// if dbError != nil {
	// 	fmt.Println(errDBConnection)
	// }

	// Instance.AutoMigrate(&tavernDomain.Item{})
	// Instance.AutoMigrate(&product.Product{})
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
