package main

import (
	"inventory-management/config"
	"inventory-management/controllers"
	"inventory-management/models"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error Loading Environment")
	}

	r := gin.Default()
	db := config.ConnectDatabase()

	db.AutoMigrate(&models.Product{}, &models.Inventory{}, &models.Order{}, &models.OrderDetail{}, &models.ProductImage{})

	productController := controllers.NewProductController(db)
	inventoryController := controllers.NewInventoryController(db)
	orderController := controllers.NewOrderController(db)
	api := r.Group("/api")
	{
		version := api.Group("/v1")
		{
			/*
				PRODUCT ROUTES
			*/
			version.POST("/product", productController.CreateProduct)
			version.GET("/products", productController.GetProducts)
			version.PUT("/products/:id", productController.UpdateProduct)
			version.DELETE("/products/:id", productController.DeleteProduct)
			version.POST("/products/:id/image", productController.UploadProductImage)
			version.GET("/products/:id/images", productController.DownloadProductImages)
			version.GET("/products/image", productController.DownloadProductImage)

			/*
				INVENTORY ROUTES
			*/
			version.PUT("/inventories/:id", inventoryController.UpdateInventory)
			version.GET("/inventories", inventoryController.GetInventory)

			/*
				ORDER ROUTES
			*/
			version.POST("/order", orderController.CreateOrder)
			version.GET("/orders/:id", orderController.GetOrder)
		}

		r.Run(":8080")
	}
}
