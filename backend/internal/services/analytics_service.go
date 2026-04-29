package services

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/Tedra-ez/AdvancedProgramming_Final/backend/internal/models"
	repository2 "github.com/Tedra-ez/AdvancedProgramming_Final/backend/internal/repository"
)

type AnalyticsService struct {
	orderRepo   repository2.OrderStore
	productRepo repository2.ProductStore
	userRepo    *repository2.UserRepository
	cacheMu     sync.RWMutex
	cacheStats  *DashboardStats
	cacheUntil  time.Time
}

func NewAnalyticsService(orderRepo repository2.OrderStore, productRepo repository2.ProductStore, userRepo *repository2.UserRepository) *AnalyticsService {
	return &AnalyticsService{
		orderRepo:   orderRepo,
		productRepo: productRepo,
		userRepo:    userRepo,
	}
}

type DashboardStats struct {
	TotalRevenue    float64            `json:"total_revenue"`
	TotalOrders     int                `json:"total_orders"`
	TotalProducts   int                `json:"total_products"`
	TotalUsers      int64              `json:"total_users"`
	PendingOrders   int                `json:"pending_orders"`
	CompletedOrders int                `json:"completed_orders"`
	RevenueByDay    []DailyRevenue     `json:"revenue_by_day"`
	TopProducts     []ProductSales     `json:"top_products"`
	RecentOrders    []*models.Order    `json:"recent_orders"`
	OrdersByStatus  map[string]int     `json:"orders_by_status"`
	SalesByCategory map[string]float64 `json:"sales_by_category"`
}

type DailyRevenue struct {
	Date    string  `json:"date"`
	Revenue float64 `json:"revenue"`
	Orders  int     `json:"orders"`
}

type ProductSales struct {
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	TotalSold   int     `json:"total_sold"`
	Revenue     float64 `json:"revenue"`
}

func (s *AnalyticsService) GetDashboardStats(ctx context.Context) (*DashboardStats, error) {
	now := time.Now()
	s.cacheMu.RLock()
	if s.cacheStats != nil && now.Before(s.cacheUntil) {
		stats := s.cacheStats
		s.cacheMu.RUnlock()
		return stats, nil
	}
	s.cacheMu.RUnlock()

	products, err := s.productRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	userCount, err := s.userRepo.Count(ctx)
	if err != nil {
		userCount = 0
	}

	stats := &DashboardStats{
		TotalProducts:   len(products),
		TotalUsers:      userCount,
		OrdersByStatus:  make(map[string]int),
		SalesByCategory: make(map[string]float64),
	}

	productMap := make(map[string]string)
	for _, p := range products {
		productMap[p.ID] = p.Name
	}

	if mongoRepo, ok := s.orderRepo.(*repository2.OrderRepositoryMongo); ok {
		summary, revenueByDay, ordersByStatus, err := mongoRepo.AggregateDashboard(ctx)
		if err != nil {
			return nil, err
		}
		stats.TotalRevenue = summary.TotalRevenue
		stats.TotalOrders = summary.TotalOrders
		stats.PendingOrders = summary.PendingOrders
		stats.CompletedOrders = summary.CompletedOrders
		for i := len(revenueByDay) - 1; i >= 0; i-- {
			agg := revenueByDay[i]
			stats.RevenueByDay = append(stats.RevenueByDay, DailyRevenue{
				Date:    agg.Date,
				Revenue: agg.Revenue,
				Orders:  agg.Orders,
			})
		}
		for _, sc := range ordersByStatus {
			stats.OrdersByStatus[sc.Status] = sc.Count
		}
		topProducts, err := mongoRepo.AggregateTopProducts(ctx, 10)
		if err != nil {
			return nil, err
		}
		for _, tp := range topProducts {
			name := tp.ProductName
			if name == "" {
				name = productMap[tp.ProductID]
			}
			stats.TopProducts = append(stats.TopProducts, ProductSales{
				ProductID:   tp.ProductID,
				ProductName: name,
				TotalSold:   tp.TotalSold,
				Revenue:     tp.Revenue,
			})
		}
		recentOrders, err := mongoRepo.FindRecent(ctx, 10)
		if err != nil {
			return nil, err
		}
		stats.RecentOrders = recentOrders

		s.cacheMu.Lock()
		s.cacheStats = stats
		s.cacheUntil = time.Now().Add(15 * time.Second)
		s.cacheMu.Unlock()

		return stats, nil
	}

	orders, err := s.orderRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	productSalesMap := make(map[string]*ProductSales)
	dailyRevenueMap := make(map[string]*DailyRevenue)

	for _, order := range orders {
		stats.TotalOrders++
		stats.TotalRevenue += order.Total
		stats.OrdersByStatus[order.Status]++

		if order.Status == "pending" {
			stats.PendingOrders++
		} else if order.Status == "completed" || order.Status == "delivered" {
			stats.CompletedOrders++
		}

		dateKey := order.CreatedAt.Format("2006-01-02")
		if _, ok := dailyRevenueMap[dateKey]; !ok {
			dailyRevenueMap[dateKey] = &DailyRevenue{Date: dateKey}
		}
		dailyRevenueMap[dateKey].Revenue += order.Total
		dailyRevenueMap[dateKey].Orders++

		for _, item := range order.Items {
			if _, ok := productSalesMap[item.ProductID]; !ok {
				productSalesMap[item.ProductID] = &ProductSales{
					ProductID:   item.ProductID,
					ProductName: item.ProductName,
				}
			}
			if productSalesMap[item.ProductID].ProductName == "" {
				productSalesMap[item.ProductID].ProductName = productMap[item.ProductID]
			}
			productSalesMap[item.ProductID].TotalSold += item.Quantity
			productSalesMap[item.ProductID].Revenue += item.LineTotal
		}
	}

	for _, dr := range dailyRevenueMap {
		stats.RevenueByDay = append(stats.RevenueByDay, *dr)
	}
	sort.Slice(stats.RevenueByDay, func(i, j int) bool {
		return stats.RevenueByDay[i].Date < stats.RevenueByDay[j].Date
	})

	if len(stats.RevenueByDay) > 30 {
		stats.RevenueByDay = stats.RevenueByDay[len(stats.RevenueByDay)-30:]
	}

	var topProducts []ProductSales
	for _, ps := range productSalesMap {
		topProducts = append(topProducts, *ps)
	}
	sort.Slice(topProducts, func(i, j int) bool {
		return topProducts[i].Revenue > topProducts[j].Revenue
	})
	if len(topProducts) > 10 {
		topProducts = topProducts[:10]
	}
	stats.TopProducts = topProducts

	recentOrders := orders
	if len(recentOrders) > 10 {
		recentOrders = recentOrders[:10]
	}
	stats.RecentOrders = recentOrders

	s.cacheMu.Lock()
	s.cacheStats = stats
	s.cacheUntil = time.Now().Add(15 * time.Second)
	s.cacheMu.Unlock()

	return stats, nil
}

func (s *AnalyticsService) GetRevenueByPeriod(ctx context.Context, startDate, endDate time.Time) ([]DailyRevenue, error) {
	if mongoRepo, ok := s.orderRepo.(*repository2.OrderRepositoryMongo); ok {
		revenueByDay, err := mongoRepo.AggregateRevenueByPeriod(ctx, startDate, endDate)
		if err != nil {
			return nil, err
		}
		result := make([]DailyRevenue, 0, len(revenueByDay))
		for _, agg := range revenueByDay {
			result = append(result, DailyRevenue{
				Date:    agg.Date,
				Revenue: agg.Revenue,
				Orders:  agg.Orders,
			})
		}
		return result, nil
	}

	orders, err := s.orderRepo.FindAll(ctx)
	if err != nil {
		return nil, err
	}

	dailyMap := make(map[string]*DailyRevenue)
	for _, order := range orders {
		if order.CreatedAt.Before(startDate) || order.CreatedAt.After(endDate) {
			continue
		}
		dateKey := order.CreatedAt.Format("2006-01-02")
		if _, ok := dailyMap[dateKey]; !ok {
			dailyMap[dateKey] = &DailyRevenue{Date: dateKey}
		}
		dailyMap[dateKey].Revenue += order.Total
		dailyMap[dateKey].Orders++
	}

	var result []DailyRevenue
	for _, dr := range dailyMap {
		result = append(result, *dr)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Date < result[j].Date
	})

	return result, nil
}

func (s *AnalyticsService) GetOrdersByUser(ctx context.Context, userID string) ([]*models.Order, error) {
	return s.orderRepo.FindByUser(ctx, userID)
}
