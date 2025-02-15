package controllers

import (
	"archive/zip"
	"errors"
	"inventory-management/models"
	"inventory-management/utils"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ProductController struct {
	DB            *gorm.DB
	downloadMutex sync.Mutex
}

func NewProductController(db *gorm.DB) *ProductController {
	return &ProductController{DB: db}
}

func (pc *ProductController) CreateProduct(c *gin.Context) {
	var req models.CreateProductRequest

	if err := utils.ValidateJson(c, &req); err != nil {
		return
	}

	var existingProduct models.Product
	if err := pc.DB.Where("name = ?", req.Name).First(&existingProduct).Error; err == nil {
		c.JSON(400, gin.H{"error": "Product name is already used"})
		return
	}

	product := models.Product{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Category:    req.Category,
	}

	tx := pc.DB.Begin()
	if err := tx.Create(&product).Error; err != nil {
		tx.Rollback()
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	productResp := models.ProductDataResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Category:    product.Category,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}

	resp := models.ProductResponse{
		Message: "Product is successfully created",
		Data:    productResp,
	}

	if err := pc.DB.First(&product).Error; err != nil {
		c.JSON(400, gin.H{"error": "Error loading product data"})
		return
	}

	c.JSON(201, resp)

}

func (pc *ProductController) UpdateProduct(c *gin.Context) {
	var req models.UpdateProductRequest
	if err := utils.ValidateJson(c, &req); err != nil {
		return
	}

	var product models.Product

	if err := pc.DB.Where("id = ?", c.Param("id")).First(&product).Error; err != nil {
		c.JSON(404, gin.H{"error": "Product is not exists"})
		return
	}

	if err := pc.DB.Where("name = ?", req.Name).First(&product).Error; err == nil {
		c.JSON(400, gin.H{"error": "Product name is already used"})
		return
	}

	tx := pc.DB.Begin()

	if req.Name != "" {
		product.Name = req.Name
	}
	if req.Description != "" {
		product.Description = req.Description
	}
	if req.Price > 0 {
		product.Price = req.Price
	}
	if req.Category != "" {
		product.Category = req.Category
	}

	if err := tx.Save(&product).Error; err != nil {
		tx.Rollback()
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	if err := pc.DB.First(&product).Error; err != nil {
		c.JSON(400, gin.H{"error": "Error loading product data"})
		return
	}

	productResp := models.ProductDataResponse{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Category:    product.Category,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}

	resp := models.ProductResponse{
		Message: "Product is successfully created",
		Data:    productResp,
	}

	c.JSON(200, resp)

}

func (pc *ProductController) DeleteProduct(c *gin.Context) {
	var product models.Product

	if err := pc.DB.Where("id = ?", c.Param("id")).First(&product).Error; err != nil {
		c.JSON(404, gin.H{"error": "Product is not exists"})
		return
	}

	tx := pc.DB.Begin()

	if err := tx.Model(&product).Association("Inventories").Clear(); err != nil {
		tx.Rollback()
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := tx.Unscoped().Delete(&product).Error; err != nil {
		tx.Rollback()
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	c.JSON(204, gin.H{})

}

func (pc *ProductController) GetProducts(c *gin.Context) {
	var products []models.Product

	id := c.Query("id")
	category := c.Query("category")

	if id != "" && category != "" {
		if err := pc.DB.Where("id = ? AND category = ?", id, category).Find(&products).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	} else if id != "" {
		if err := pc.DB.Where("id = ?", id).Find(&products).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	} else if category != "" {
		if err := pc.DB.Where("category = ?", category).Find(&products).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	} else {
		if err := pc.DB.Find(&products).Error; err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}

	resp := make([]models.ProductDataResponse, len(products))
	for i, product := range products {
		resp[i] = models.ProductDataResponse{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Category:    product.Category,
			CreatedAt:   product.CreatedAt,
			UpdatedAt:   product.UpdatedAt,
		}
	}

	c.JSON(200, gin.H{"data": resp})
}

func (pc *ProductController) UploadProductImage(c *gin.Context) {
	const maxFileSize = 2 * 1024 * 1024
	const uploadDirectory = "uploads"

	rawProductId, err := strconv.ParseUint(c.Param("id"), 10, 0)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid product ID (must be a number)"})

		return
	}

	productId := uint(rawProductId)

	image, err := c.FormFile("image")
	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to upload image"})
		return
	}

	if image.Size > maxFileSize {
		c.JSON(400, gin.H{"error": "Maximum File Size is 2 mb "})
		return
	}

	if !utils.IsValidExtension(image) {
		c.JSON(400, gin.H{"error": "Please upload file with extension PNG, JPG, JPEG"})
		return
	}

	tx := pc.DB.Begin()

	imagePath := filepath.Join("uploads", image.Filename)
	if err := c.SaveUploadedFile(image, imagePath); err != nil {
		c.JSON(400, gin.H{"error": "Error uploading file"})
		return
	}

	productImage := models.ProductImage{
		ProductId:  productId,
		ImageUrl:   imagePath,
		FileName:   image.Filename,
		UploadedAt: time.Now(),
	}
	if err := tx.Create(&productImage).Error; err != nil {
		tx.Rollback()
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	c.JSON(200, gin.H{
		"message": "Image uploaded successfully",
	})

}

func (pc *ProductController) DownloadProductImage(c *gin.Context) {
	rawProductId := c.Query("productId")
	RawProductImageId := c.Query("productImageId")

	if rawProductId == "" || RawProductImageId == "" {
		c.JSON(400, gin.H{"error": "productId and productImageId are required"})
		return
	}

	productId, err := strconv.ParseUint(rawProductId, 10, 0)
	if err != nil {
		c.JSON(400, gin.H{"error": "productIdmust be a number"})
		return
	}

	productImageId, err := strconv.ParseUint(RawProductImageId, 10, 0)
	if err != nil {
		c.JSON(400, gin.H{"error": "productImageIdmust be a number"})
		return
	}

	var productImage models.ProductImage
	result := pc.DB.Where("product_id = ? AND id = ?", productId, productImageId).First(&productImage)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			c.JSON(404, gin.H{"error": "Image not found"})
		} else {
			c.JSON(500, gin.H{"error": result.Error.Error()})
		}
		return
	}

	c.File(productImage.ImageUrl)
}

func (pc *ProductController) DownloadProductImages(c *gin.Context) {
	rawProductId := c.Param("id")
	productId, err := strconv.ParseUint(rawProductId, 10, 0)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid product ID"})
		return
	}

	var productImages []models.ProductImage
	result := pc.DB.Where("product_id = ?", productId).Find(&productImages)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	var product models.Product
	result = pc.DB.Where("id = ?", productId).First(&product)
	if result.Error != nil {
		c.JSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	if len(productImages) == 0 {
		c.JSON(404, gin.H{"error": "No images found for this product"})
		return
	}

	uuid := uuid.New().String()
	zipFileName := filepath.Join("uploads", uuid+".zip")

	zipFile, err := os.Create(zipFileName)
	if err != nil {
		c.JSON(500, gin.H{"error": "Error creating ZIP file"})
		return
	}
	zipWriter := zip.NewWriter(zipFile)

	for _, img := range productImages {
		file, err := os.Open(img.ImageUrl)
		if err != nil {
			zipWriter.Close()
			zipFile.Close()
			c.JSON(500, gin.H{"error": "Error opening image file: " + img.FileName})
			return
		}

		zipEntry, err := zipWriter.Create(img.FileName)
		if err != nil {
			file.Close()
			zipWriter.Close()
			zipFile.Close()
			c.JSON(500, gin.H{"error": "Error adding file to ZIP: " + img.FileName})
			return
		}

		if _, err := io.Copy(zipEntry, file); err != nil {
			file.Close()
			zipWriter.Close()
			zipFile.Close()
			c.JSON(500, gin.H{"error": "Error writing file to ZIP: " + img.FileName})
			return
		}

		file.Close()
	}

	zipWriter.Close()
	zipFile.Close()

	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", "attachment; filename="+filepath.Base(zipFileName))
	c.File(zipFileName)

	go func() {
		time.Sleep(10 * time.Second)
		os.Remove(zipFileName)
	}()
}
