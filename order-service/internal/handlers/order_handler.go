package handlers

import (
	"net/http"
	"time"

	"github.com/Tedra-ez/AdvancedProgramming_Final/order-service/internal/models"
	"github.com/Tedra-ez/AdvancedProgramming_Final/order-service/internal/services"
	"github.com/Tedra-ez/AdvancedProgramming_Final/pkg/events"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	svc       *services.OrderService
	publisher *events.Publisher
}

func NewOrderHandler(svc *services.OrderService, publisher *events.Publisher) *OrderHandler {
	return &OrderHandler{svc: svc, publisher: publisher}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req models.CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	order, err := h.svc.Create(c.Request.Context(), &req)
	if err != nil {
		if err == services.ErrUserNotFound || err == services.ErrProductNotFound {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	h.publisher.PublishOrderCreated(events.OrderCreatedEvent{
		EventID:   time.Now().UTC().Format("20060102150405.000000000"),
		OrderID:   order.ID,
		UserID:    order.UserID,
		Status:    order.Status,
		Total:     order.Total,
		ItemCount: len(order.Items),
		CreatedAt: order.CreatedAt,
	})
	c.JSON(http.StatusCreated, order)
}

func (h *OrderHandler) GetOrderStatus(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order id required"})
		return
	}
	order, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if order == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}
	c.JSON(http.StatusOK, order)
}

func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "order id required"})
		return
	}
	var body struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.UpdateStatus(c.Request.Context(), id, body.Status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id, "status": body.Status})
}

func (h *OrderHandler) ListOrdersByUser(c *gin.Context) {
	userID := c.Param("userId")
	if userID == "" {
		userID = c.Query("user_id")
	}
	if userID == "" {
		if role, _ := c.Get("user_role"); role == "admin" {
			orders, err := h.svc.ListAll(c.Request.Context())
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, orders)
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id required (query or path)"})
		return
	}
	orders, err := h.svc.ListByUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}
