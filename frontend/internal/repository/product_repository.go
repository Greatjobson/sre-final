package repository

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"

	"github.com/Tedra-ez/AdvancedProgramming_Final/frontend/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductStore interface {
	FindAll(ctx context.Context) ([]*models.Product, error)
	FindByID(ctx context.Context, id string) (*models.Product, error)
	Insert(ctx context.Context, p *models.Product) (*models.Product, error)
	Update(ctx context.Context, id string, p *models.Product) error
	Delete(ctx context.Context, id string) error
}

type ProductRepositoryMongo struct {
	coll *mongo.Collection
}

func NewProductRepositoryMongo(coll *mongo.Collection) *ProductRepositoryMongo {
	return &ProductRepositoryMongo{coll: coll}
}

func (r *ProductRepositoryMongo) FindAll(ctx context.Context) ([]*models.Product, error) {
	cur, err := r.coll.Find(ctx, bson.M{}, options.Find())
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var out []*models.Product
	for cur.Next(ctx) {
		var doc productDoc
		if err := cur.Decode(&doc); err != nil {
			return nil, err
		}
		out = append(out, doc.toModel())
	}
	return out, cur.Err()
}

func (r *ProductRepositoryMongo) FindByID(ctx context.Context, id string) (*models.Product, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, nil
	}
	var doc productDoc
	err = r.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&doc)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return doc.toModel(), nil
}

func (r *ProductRepositoryMongo) Insert(ctx context.Context, p *models.Product) (*models.Product, error) {
	now := time.Now()
	if p.CreatedAt.IsZero() {
		p.CreatedAt = now
	}
	if p.UpdatedAt.IsZero() {
		p.UpdatedAt = now
	}
	id := primitive.NewObjectID()
	doc := productDocFromModel(p)
	doc.ID = id
	_, err := r.coll.InsertOne(ctx, doc)
	if err != nil {
		return nil, err
	}
	p.ID = id.Hex()
	return p, nil
}

func (r *ProductRepositoryMongo) Update(ctx context.Context, id string, p *models.Product) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	p.UpdatedAt = time.Now()
	doc := productDocFromModel(p)
	doc.ID = oid
	_, err = r.coll.ReplaceOne(ctx, bson.M{"_id": oid}, doc)
	return err
}

func (r *ProductRepositoryMongo) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.coll.DeleteOne(ctx, bson.M{"_id": oid})
	return err
}

type productDoc struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	Name        string             `bson:"name"`
	Description string             `bson:"description"`
	Category    string             `bson:"category"`
	Gender      string             `bson:"gender"`
	Price       float64            `bson:"price"`
	Sizes       []string           `bson:"sizes"`
	Colors      []string           `bson:"colors"`
	StockBySize map[string]int     `bson:"stockBySize"`
	Images      []string           `bson:"images"`
	IsActive    bool               `bson:"isActive"`
	CreatedAt   primitive.DateTime `bson:"createdAt"`
	UpdatedAt   primitive.DateTime `bson:"updateAt"`
}

func productDocFromModel(p *models.Product) *productDoc {
	d := &productDoc{
		Name:        p.Name,
		Description: p.Description,
		Category:    p.Category,
		Gender:      p.Gender,
		Price:       p.Price,
		Sizes:       p.Sizes,
		Colors:      p.Colors,
		StockBySize: p.StockBySize,
		Images:      p.Images,
		IsActive:    p.IsActive,
	}
	if !p.CreatedAt.IsZero() {
		d.CreatedAt = primitive.NewDateTimeFromTime(p.CreatedAt)
	}
	if !p.UpdatedAt.IsZero() {
		d.UpdatedAt = primitive.NewDateTimeFromTime(p.UpdatedAt)
	}
	return d
}

func (d *productDoc) toModel() *models.Product {
	return &models.Product{
		ID:          d.ID.Hex(),
		Name:        d.Name,
		Description: d.Description,
		Category:    d.Category,
		Gender:      d.Gender,
		Price:       d.Price,
		Sizes:       d.Sizes,
		Colors:      d.Colors,
		StockBySize: d.StockBySize,
		Images:      d.Images,
		IsActive:    d.IsActive,
		CreatedAt:   d.CreatedAt.Time(),
		UpdatedAt:   d.UpdatedAt.Time(),
	}
}

type ProductRepositoryMemory struct {
	mu   sync.RWMutex
	data map[string]*models.Product
}

func NewProductRepositoryMemory() *ProductRepositoryMemory {
	return &ProductRepositoryMemory{data: make(map[string]*models.Product)}
}

func (r *ProductRepositoryMemory) FindAll(ctx context.Context) ([]*models.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]*models.Product, 0, len(r.data))
	for _, p := range r.data {
		out = append(out, p)
	}
	return out, nil
}

func (r *ProductRepositoryMemory) FindByID(ctx context.Context, id string) (*models.Product, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.data[id], nil
}

func (r *ProductRepositoryMemory) Insert(ctx context.Context, p *models.Product) (*models.Product, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if p.ID == "" {
		b := make([]byte, 8)
		if _, err := rand.Read(b); err == nil {
			p.ID = hex.EncodeToString(b)
		} else {
			p.ID = time.Now().Format("20060102150405.000")
		}
	}
	if p.CreatedAt.IsZero() {
		p.CreatedAt = time.Now()
	}
	if p.UpdatedAt.IsZero() {
		p.UpdatedAt = time.Now()
	}
	r.data[p.ID] = p
	return p, nil
}

func (r *ProductRepositoryMemory) Update(ctx context.Context, id string, p *models.Product) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.data[id]; !ok {
		return nil
	}
	p.ID = id
	p.UpdatedAt = time.Now()
	r.data[id] = p
	return nil
}

func (r *ProductRepositoryMemory) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.data, id)
	return nil
}
