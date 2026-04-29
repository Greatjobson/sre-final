package handlers

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Tedra-ez/AdvancedProgramming_Final/backend/internal/models"
	"github.com/Tedra-ez/AdvancedProgramming_Final/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type ProductHandler struct {
	productService *services.ProductService
}

func NewProductHandler(svc *services.ProductService) *ProductHandler {
	return &ProductHandler{productService: svc}
}

func (h *ProductHandler) GetProducts(c *gin.Context) {
	products, err := h.productService.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if products == nil {
		products = []*models.Product{}
	}
	c.JSON(http.StatusOK, products)
}

func (h *ProductHandler) GetProductByID(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id required"})
		return
	}
	p, err := h.productService.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if p == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}
	c.JSON(http.StatusOK, p)
}

func (h *ProductHandler) CreateProduct(c *gin.Context) {
	name := strings.TrimSpace(c.PostForm("name"))
	priceStr := strings.TrimSpace(c.PostForm("price"))
	category := strings.TrimSpace(c.PostForm("category"))
	gender := strings.TrimSpace(c.PostForm("gender"))
	description := strings.TrimSpace(c.PostForm("description"))
	sizesStr := c.PostForm("sizes")
	colorsStr := c.PostForm("colors")
	stockStr := c.PostForm("stock")

	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
		return
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil || price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "valid price is required and must be > 0"})
		return
	}

	file, err := c.FormFile("image")
	var imagePath string
	if err == nil {
		uploadDir := resolveUploadDir(
			filepath.Join("static", "assets", "products"),
			filepath.Join("frontend", "static", "assets", "products"),
		)
		if err := os.MkdirAll(uploadDir, 0o755); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to prepare image directory"})
			return
		}
		path := filepath.Join(uploadDir, file.Filename)
		if err := c.SaveUploadedFile(file, path); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save image"})
			return
		}
		imagePath = "/static/assets/products/" + file.Filename
	}

	req := models.CreateProductRequest{
		Name:        name,
		Description: description,
		Category:    category,
		Gender:      gender,
		Price:       price,
		Sizes:       parseCommaString(sizesStr),
		Colors:      parseCommaString(colorsStr),
		StockBySize: parseStockString(stockStr),
	}
	if imagePath != "" {
		req.Images = []string{imagePath}
	}

	p, err := h.productService.Create(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, p)
}

func parseCommaString(s string) []string {
	if s == "" {
		return []string{}
	}
	parts := strings.Split(s, ",")
	var result []string
	for _, p := range parts {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func parseStockString(s string) map[string]int {
	result := make(map[string]int)
	if s == "" {
		return result
	}
	parts := strings.Split(s, ",")
	for _, entry := range parts {
		pair := strings.Split(strings.TrimSpace(entry), ":")
		if len(pair) == 2 {
			size := strings.TrimSpace(pair[0])
			qty, err := strconv.Atoi(strings.TrimSpace(pair[1]))
			if err == nil && size != "" {
				result[size] = qty
			}
		}
	}
	return result
}

func resolveUploadDir(candidates ...string) string {
	for _, candidate := range candidates {
		info, err := os.Stat(candidate)
		if err == nil && info.IsDir() {
			return candidate
		}
	}
	return candidates[0]
}

func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id required"})
		return
	}
	var p models.Product
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	p.ID = id
	if err := h.productService.Update(c.Request.Context(), id, &p); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, p)
}

func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id required"})
		return
	}
	if err := h.productService.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusNoContent, nil)
}
