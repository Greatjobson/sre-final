package repository

import (
	"context"
	"fmt"
	"sync"
	"time"

	models2 "github.com/Tedra-ez/AdvancedProgramming_Final/frontend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OrderStore interface {
	Save(ctx context.Context, order *models2.Order) error
	FindByUser(ctx context.Context, userID string) ([]*models2.Order, error)
	FindAll(ctx context.Context) ([]*models2.Order, error)
	UpdateStatus(ctx context.Context, orderID, status string) error
	FindByID(ctx context.Context, orderID string) (*models2.Order, error)
}

type OrderRepositoryMongo struct {
	coll     *mongo.Collection
	itemRepo OrderItemStore
}

func NewOrderRepositoryMongo(coll *mongo.Collection, itemRepo OrderItemStore) *OrderRepositoryMongo {
	return &OrderRepositoryMongo{coll: coll, itemRepo: itemRepo}
}

func (r *OrderRepositoryMongo) Save(ctx context.Context, order *models2.Order) error {
	if order.CreatedAt.IsZero() {
		order.CreatedAt = time.Now()
	}
	order.UpdatedAt = time.Now()
	if order.ID == "" {
		id := primitive.NewObjectID()
		order.ID = id.Hex()
		for i := range order.Items {
			order.Items[i].OrderID = order.ID
		}
		doc := orderDocFromModel(order)
		doc.ID = id
		if _, err := r.coll.InsertOne(ctx, doc); err != nil {
			return err
		}
		return r.itemRepo.CreateMany(ctx, order.Items)
	}
	oid, err := primitive.ObjectIDFromHex(order.ID)
	if err != nil {
		return err
	}
	for i := range order.Items {
		order.Items[i].OrderID = order.ID
	}
	if err := r.itemRepo.CreateMany(ctx, order.Items); err != nil {
		return err
	}
	doc := orderDocFromModel(order)
	doc.ID = oid
	_, err = r.coll.ReplaceOne(ctx, bson.M{"_id": oid}, doc)
	return err
}

func (r *OrderRepositoryMongo) FindByUser(ctx context.Context, userID string) ([]*models2.Order, error) {
	return r.findOrders(ctx, bson.M{"userId": userID})
}

func (r *OrderRepositoryMongo) FindAll(ctx context.Context) ([]*models2.Order, error) {
	return r.findOrders(ctx, bson.M{})
}

func (r *OrderRepositoryMongo) FindRecent(ctx context.Context, limit int) ([]*models2.Order, error) {
	if limit <= 0 {
		limit = 10
	}
	cur, err := r.coll.Find(ctx, bson.M{}, options.Find().SetSort(bson.M{"createdAt": -1}).SetLimit(int64(limit)))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var out []*models2.Order
	var orderIDs []string
	for cur.Next(ctx) {
		var doc orderDoc
		if err := cur.Decode(&doc); err != nil {
			return nil, err
		}
		o := doc.toModel()
		out = append(out, o)
		orderIDs = append(orderIDs, o.ID)
	}
	itemsByOrderID, err := r.itemRepo.FindByOrderIds(ctx, orderIDs)
	if err != nil {
		return nil, err
	}
	for _, o := range out {
		items := itemsByOrderID[o.ID]
		o.Items = make([]models2.OrderItem, 0, len(items))
		for _, it := range items {
			o.Items = append(o.Items, *it)
		}
	}
	return out, cur.Err()
}

func (r *OrderRepositoryMongo) findOrders(ctx context.Context, filter bson.M) ([]*models2.Order, error) {
	cur, err := r.coll.Find(ctx, filter, options.Find().SetSort(bson.M{"createdAt": -1}))
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var out []*models2.Order
	var orderIDs []string
	for cur.Next(ctx) {
		var doc orderDoc
		if err := cur.Decode(&doc); err != nil {
			return nil, err
		}
		o := doc.toModel()
		out = append(out, o)
		orderIDs = append(orderIDs, o.ID)
	}
	itemsByOrderID, err := r.itemRepo.FindByOrderIds(ctx, orderIDs)
	if err != nil {
		return nil, err
	}
	for _, o := range out {
		items := itemsByOrderID[o.ID]
		o.Items = make([]models2.OrderItem, 0, len(items))
		for _, it := range items {
			o.Items = append(o.Items, *it)
		}
	}
	return out, cur.Err()
}

func (r *OrderRepositoryMongo) FindByID(ctx context.Context, orderID string) (*models2.Order, error) {
	oid, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		return nil, nil
	}
	var doc orderDoc
	err = r.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&doc)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	o := doc.toModel()
	items, _ := r.itemRepo.FindByOrderId(ctx, orderID)
	o.Items = make([]models2.OrderItem, 0, len(items))
	for _, it := range items {
		o.Items = append(o.Items, *it)
	}
	return o, nil
}

func (r *OrderRepositoryMongo) UpdateStatus(ctx context.Context, orderID, status string) error {
	oid, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		return err
	}
	_, err = r.coll.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": bson.M{"status": status, "updatedAt": primitive.NewDateTimeFromTime(time.Now())}})
	return err
}

type orderDoc struct {
	ID              primitive.ObjectID `bson:"_id,omitempty"`
	UserID          string             `bson:"userId"`
	Status          string             `bson:"status"`
	PaymentMethod   string             `bson:"paymentMethod"`
	DeliveryMethod  string             `bson:"deliveryMethod"`
	DeliveryAddress string             `bson:"deliveryAddress"`
	Comment         string             `bson:"comment"`
	Subtotal        float64            `bson:"subtotal"`
	DeliveryFee     float64            `bson:"deliveryFee"`
	Total           float64            `bson:"total"`
	CreatedAt       primitive.DateTime `bson:"createdAt"`
	UpdatedAt       primitive.DateTime `bson:"updatedAt"`
}

func orderDocFromModel(o *models2.Order) *orderDoc {
	return &orderDoc{
		UserID:          o.UserID,
		Status:          o.Status,
		PaymentMethod:   o.PaymentMethod,
		DeliveryMethod:  o.DeliveryMethod,
		DeliveryAddress: o.DeliveryAddress,
		Comment:         o.Comment,
		Subtotal:        o.Subtotal,
		DeliveryFee:     o.DeliveryFee,
		Total:           o.Total,
		CreatedAt:       primitive.NewDateTimeFromTime(o.CreatedAt),
		UpdatedAt:       primitive.NewDateTimeFromTime(o.UpdatedAt),
	}
}

func (d *orderDoc) toModel() *models2.Order {
	return &models2.Order{
		ID:              d.ID.Hex(),
		UserID:          d.UserID,
		Status:          d.Status,
		PaymentMethod:   d.PaymentMethod,
		DeliveryMethod:  d.DeliveryMethod,
		DeliveryAddress: d.DeliveryAddress,
		Comment:         d.Comment,
		Subtotal:        d.Subtotal,
		DeliveryFee:     d.DeliveryFee,
		Total:           d.Total,
		CreatedAt:       d.CreatedAt.Time(),
		UpdatedAt:       d.UpdatedAt.Time(),
	}
}

type OrderRepositoryMemory struct {
	mu    sync.RWMutex
	data  map[string]*models2.Order
	idGen int
}

func NewOrderRepositoryMemory() *OrderRepositoryMemory {
	return &OrderRepositoryMemory{data: make(map[string]*models2.Order)}
}

func (r *OrderRepositoryMemory) Save(ctx context.Context, order *models2.Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if order.ID == "" {
		r.idGen++
		order.ID = fmt.Sprintf("order-%d-%d", time.Now().UnixNano(), r.idGen)
	}
	if order.CreatedAt.IsZero() {
		order.CreatedAt = time.Now()
	}
	order.UpdatedAt = time.Now()
	for i := range order.Items {
		order.Items[i].OrderID = order.ID
	}
	r.data[order.ID] = order
	return nil
}

func (r *OrderRepositoryMemory) FindByUser(ctx context.Context, userID string) ([]*models2.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var out []*models2.Order
	for _, o := range r.data {
		if o.UserID == userID {
			out = append(out, o)
		}
	}
	return out, nil
}

func (r *OrderRepositoryMemory) FindAll(ctx context.Context) ([]*models2.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*models2.Order, 0, len(r.data))
	for _, o := range r.data {
		out = append(out, o)
	}
	return out, nil
}

func (r *OrderRepositoryMemory) FindByID(ctx context.Context, orderID string) (*models2.Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.data[orderID], nil
}

func (r *OrderRepositoryMemory) UpdateStatus(ctx context.Context, orderID, status string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if o, ok := r.data[orderID]; ok {
		o.Status = status
		o.UpdatedAt = time.Now()
	}
	return nil
}
