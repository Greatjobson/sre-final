package handlers

import (
	"net/http"
	"time"

	"github.com/Tedra-ez/AdvancedProgramming_Final/order-service/internal/services"
	"github.com/gin-gonic/gin"
)

type AnalyticsHandler struct {
	svc *services.AnalyticsService
}

func NewAnalyticsHandler(svc *services.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{svc: svc}
}

func (h *AnalyticsHandler) TopProductsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		stats, err := h.svc.GetDashboardStats(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, stats.TopProducts)
	}
}

func (h *AnalyticsHandler) RevenueHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		startDate := c.Query("start_date")
		endDate := c.Query("end_date")

		var result interface{}
		var err error

		if startDate != "" && endDate != "" {
			start, _ := time.Parse("2006-01-02", startDate)
			end, _ := time.Parse("2006-01-02", endDate)
			result, err = h.svc.GetRevenueByPeriod(c.Request.Context(), start, end)
		} else {
			stats, statErr := h.svc.GetDashboardStats(c.Request.Context())
			if statErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": statErr.Error()})
				return
			}
			result = gin.H{
				"total_revenue":  stats.TotalRevenue,
				"revenue_by_day": stats.RevenueByDay,
			}
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

func (h *AnalyticsHandler) OrdersByStatusHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		stats, err := h.svc.GetDashboardStats(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, stats.OrdersByStatus)
	}
}

func (h *AnalyticsHandler) DashboardStatsHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		stats, err := h.svc.GetDashboardStats(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, stats)
	}
}
