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

	db.AutoMigrate(&models.Product{})

	productController := controllers.NewProductController(db)
	api := r.Group("/api")
	{
		version := api.Group("/v1")
		{
			version.POST("/product", productController.CreateProduct)
			version.GET("/products", productController.GetProducts)
			version.PUT("/products/:id", productController.UpdateProduct)
			version.DELETE("/products/:id", productController.DeleteProduct)
		}

		r.Run(":8080")
	}
}
