package controllers

import (
	"inventory-management/models"
	"inventory-management/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderController struct {
	DB *gorm.DB
}

func NewOrderController(db *gorm.DB) *OrderController {
	return &OrderController{DB: db}
}

func (oc *OrderController) CreateOrder(c *gin.Context) {
	var req models.CreateOrderRequest
	if err := utils.ValidateJson(c, &req); err != nil {
		return
	}

	productIds := make(map[uint]bool)
	for _, d := range req.OrderDetails {
		if productIds[d.ProductId] {
			c.JSON(400, gin.H{"error": "Duplicate Product ID: " + strconv.FormatUint(uint64(d.ProductId), 10)})
			return
		}
		productIds[d.ProductId] = true
	}

	tx := oc.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var total uint
	var orderDetails []models.OrderDetail

	for _, d := range req.OrderDetails {
		var product models.Product
		if err := tx.First(&product, d.ProductId).Error; err != nil {
			tx.Rollback()
			c.JSON(400, gin.H{"error": "Product not found: " + strconv.FormatUint(uint64(d.ProductId), 10)})
			return
		}

		total += d.Qty * product.Price
		orderDetails = append(orderDetails, models.OrderDetail{
			ProductId:   d.ProductId,
			Qty:         d.Qty,
			PriceDetail: product.Price * d.Qty,
		})

		var inventory models.Inventory
		if err := tx.Where("product_id = ?", d.ProductId).First(&inventory).Error; err != nil {
			tx.Rollback()
			c.JSON(400, gin.H{"error": "Inventory for product not available: " + strconv.FormatUint(uint64(d.ProductId), 10)})
			return
		}

		if inventory.Stock < d.Qty {
			tx.Rollback()
			c.JSON(400, gin.H{"error": "Insufficient stock for product: " + strconv.FormatUint(uint64(d.ProductId), 10)})
			return
		}

		if err := tx.Model(&inventory).Update("stock", gorm.Expr("stock - ?", d.Qty)).Error; err != nil {
			tx.Rollback()
			c.JSON(500, gin.H{"error": "failed to update inventory"})
			return
		}
	}

	order := models.Order{
		Date:            time.Now(),
		CustomerName:    req.CustomerName,
		ShippingAddress: req.ShippingAddress,
		Total:           total,
		OrderDetails:    orderDetails,
	}

	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": "failed to create order"})
		return
	}

	tx.Commit()

	orderResponse := models.OrderResponse{
		ID:                  order.ID,
		Date:                order.Date,
		CustomerName:        order.CustomerName,
		ShippingAddress:     order.ShippingAddress,
		Total:               order.Total,
		OrderDetailResponse: make([]models.OrderDetailResponse, len(order.OrderDetails)),
	}

	for i, od := range order.OrderDetails {
		orderResponse.OrderDetailResponse[i] = models.OrderDetailResponse{
			ProductId:   od.ProductId,
			Qty:         od.Qty,
			PriceDetail: od.PriceDetail,
		}
	}

	c.JSON(200, orderResponse)
}

func (oc *OrderController) GetOrder(c *gin.Context) {
	id := c.Param("id")

	var order models.Order

	if err := oc.DB.Preload("OrderDetails").First(&order, id).Error; err != nil {
		c.JSON(400, gin.H{"error": "Order not found"})
		return
	}

	orderResponse := models.OrderResponse{
		ID:                  order.ID,
		Date:                order.Date,
		CustomerName:        order.CustomerName,
		ShippingAddress:     order.ShippingAddress,
		Total:               order.Total,
		OrderDetailResponse: make([]models.OrderDetailResponse, len(order.OrderDetails)),
	}

	for i, od := range order.OrderDetails {
		orderResponse.OrderDetailResponse[i] = models.OrderDetailResponse{
			ProductId:   od.ProductId,
			Qty:         od.Qty,
			PriceDetail: od.PriceDetail,
		}
	}

	c.JSON(200, orderResponse)
}
