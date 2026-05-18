package repository

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderSummaryAgg struct {
	TotalRevenue    float64 `bson:"totalRevenue"`
	TotalOrders     int     `bson:"totalOrders"`
	PendingOrders   int     `bson:"pendingOrders"`
	CompletedOrders int     `bson:"completedOrders"`
}

type OrderRevenueByDayAgg struct {
	Date    string  `bson:"_id"`
	Revenue float64 `bson:"revenue"`
	Orders  int     `bson:"orders"`
}

type OrderStatusCountAgg struct {
	Status string `bson:"_id"`
	Count  int    `bson:"count"`
}

type TopProductAgg struct {
	ProductID   string  `bson:"_id"`
	ProductName string  `bson:"productName"`
	TotalSold   int     `bson:"totalSold"`
	Revenue     float64 `bson:"revenue"`
}

func (r *OrderRepositoryMongo) AggregateDashboard(ctx context.Context) (OrderSummaryAgg, []OrderRevenueByDayAgg, []OrderStatusCountAgg, error) {
	pendingCond := bson.M{
		"$cond": bson.A{
			bson.M{"$eq": bson.A{"$status", "pending"}},
			1,
			0,
		},
	}
	completedCond := bson.M{
		"$cond": bson.A{
			bson.M{"$in": bson.A{"$status", bson.A{"completed", "delivered"}}},
			1,
			0,
		},
	}
	summaryStage := bson.D{{"$group", bson.M{
		"_id":             nil,
		"totalRevenue":    bson.M{"$sum": "$total"},
		"totalOrders":     bson.M{"$sum": 1},
		"pendingOrders":   bson.M{"$sum": pendingCond},
		"completedOrders": bson.M{"$sum": completedCond},
	}}}
	revenueByDayStage := bson.A{
		bson.D{{"$group", bson.M{
			"_id": bson.M{
				"$dateToString": bson.M{
					"format": "%Y-%m-%d",
					"date":   "$createdAt",
				},
			},
			"revenue": bson.M{"$sum": "$total"},
			"orders":  bson.M{"$sum": 1},
		}}},
		bson.D{{"$sort", bson.M{"_id": -1}}},
		bson.D{{"$limit", 30}},
	}
	ordersByStatusStage := bson.A{
		bson.D{{"$group", bson.M{
			"_id":   "$status",
			"count": bson.M{"$sum": 1},
		}}},
	}
	pipeline := mongo.Pipeline{
		bson.D{{"$facet", bson.D{
			{"summary", bson.A{summaryStage}},
			{"revenueByDay", revenueByDayStage},
			{"ordersByStatus", ordersByStatusStage},
		}}},
	}

	cur, err := r.coll.Aggregate(ctx, pipeline)
	if err != nil {
		return OrderSummaryAgg{}, nil, nil, err
	}
	var result []struct {
		Summary        []OrderSummaryAgg      `bson:"summary"`
		RevenueByDay   []OrderRevenueByDayAgg `bson:"revenueByDay"`
		OrdersByStatus []OrderStatusCountAgg  `bson:"ordersByStatus"`
	}
	if err := cur.All(ctx, &result); err != nil {
		return OrderSummaryAgg{}, nil, nil, err
	}
	if len(result) == 0 {
		return OrderSummaryAgg{}, nil, nil, nil
	}
	summary := OrderSummaryAgg{}
	if len(result[0].Summary) > 0 {
		summary = result[0].Summary[0]
	}
	return summary, result[0].RevenueByDay, result[0].OrdersByStatus, nil
}

func (r *OrderRepositoryMongo) AggregateRevenueByPeriod(ctx context.Context, startDate, endDate time.Time) ([]OrderRevenueByDayAgg, error) {
	pipeline := mongo.Pipeline{
		bson.D{{"$match", bson.D{{"createdAt", bson.D{{"$gte", startDate}, {"$lte", endDate}}}}}},
		bson.D{{"$group", bson.D{
			{"_id", bson.D{{"$dateToString", bson.D{{"format", "%Y-%m-%d"}, {"date", "$createdAt"}}}}},
			{"revenue", bson.D{{"$sum", "$total"}}},
			{"orders", bson.D{{"$sum", 1}}},
		}}},
		bson.D{{"$sort", bson.D{{"_id", 1}}}},
	}

	cur, err := r.coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	var result []OrderRevenueByDayAgg
	if err := cur.All(ctx, &result); err != nil {
		return nil, err
	}
	return result, nil
}

func (r *OrderRepositoryMongo) AggregateTopProducts(ctx context.Context, limit int) ([]TopProductAgg, error) {
	if limit <= 0 {
		limit = 10
	}
	pipeline := mongo.Pipeline{
		bson.D{{"$addFields", bson.D{{"orderId", bson.D{{"$toString", "$_id"}}}}}},
		bson.D{{"$lookup", bson.D{
			{"from", "order_items"},
			{"localField", "orderId"},
			{"foreignField", "orderId"},
			{"as", "items"},
		}}},
		bson.D{{"$unwind", "$items"}},
		bson.D{{"$group", bson.D{
			{"_id", "$items.productId"},
			{"productName", bson.D{{"$first", "$items.productName"}}},
			{"totalSold", bson.D{{"$sum", "$items.quantity"}}},
			{"revenue", bson.D{{"$sum", "$items.lineTotal"}}},
		}}},
		bson.D{{"$sort", bson.D{{"revenue", -1}}}},
		bson.D{{"$limit", limit}},
	}

	cur, err := r.coll.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	var result []TopProductAgg
	if err := cur.All(ctx, &result); err != nil {
		return nil, err
	}
	return result, nil
}
