package handlers

import (
	"html/template"
	"net/http"
	"net/url"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Tedra-ez/AdvancedProgramming_Final/backend/internal/models"
	services2 "github.com/Tedra-ez/AdvancedProgramming_Final/backend/internal/services"
	"github.com/gin-gonic/gin"
)

type PageHandler struct {
	productService   *services2.ProductService
	orderService     *services2.OrderService
	authService      *services2.AuthService
	analyticsService *services2.AnalyticsService
	templates        map[string]*template.Template
}

func NewPageHandler(productService *services2.ProductService, orderService *services2.OrderService, authService *services2.AuthService, analyticsService *services2.AnalyticsService, templateDir string) (*PageHandler, error) {
	basePath := filepath.Join(templateDir, "base.html")
	pages := []string{
		"shop", "index", "account", "login", "register",
		"admin_orders", "admin_products", "admin_dashboard",
		"admin_users", "admin_analytics", "account_orders",
		"product", "wishlist", "cart", "checkout",
	}

	templates := make(map[string]*template.Template)
	for _, p := range pages {
		path := filepath.Join(templateDir, p+".html")
		templates[p] = template.Must(template.ParseFiles(basePath, path))
	}

	return &PageHandler{
		productService:   productService,
		orderService:     orderService,
		authService:      authService,
		analyticsService: analyticsService,
		templates:        templates,
	}, nil
}

func (h *PageHandler) getUserData(c *gin.Context) gin.H {
	data := gin.H{}
	if id, ok := c.Get("user_id"); ok && id != "" {
		data["User"] = map[string]string{
			"id":    id.(string),
			"role":  getStr(c, "user_role"),
			"email": getStr(c, "user_email"),
			"name":  getStr(c, "user_name"),
		}
	}
	return data
}

func (h *PageHandler) Index(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	if err := h.templates["index"].ExecuteTemplate(c.Writer, "base.html", h.getUserData(c)); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func (h *PageHandler) Shop(c *gin.Context) {
	products, err := h.productService.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	q := c.Query("q")
	sort := c.Query("sort")
	if sort == "" {
		sort = "recommended"
	}
	categories := c.QueryArray("category")
	genders := c.QueryArray("gender")
	colors := c.QueryArray("color")
	sizes := c.QueryArray("size")
	selectedCategories := toSelectionMap(categories)
	selectedGenders := toSelectionMap(genders)
	selectedColors := toSelectionMap(colors)
	selectedSizes := toSelectionMap(sizes)
	categoryCounts := buildCategoryCounts(products)
	colorCounts := buildColorCounts(products)
	filtered := filterProducts(products, q, categories, colors, sizes, genders)
	ordered := sortProducts(filtered, sort)
	chips, clearURL := buildFilterChips(c.Request.URL.Query())

	data := h.getUserData(c)
	data["Products"] = ordered
	data["SearchQuery"] = q
	data["Sort"] = sort
	data["ShowSidebar"] = true
	data["SelectedCategoryList"] = categories
	data["SelectedGenderList"] = genders
	data["SelectedColorList"] = colors
	data["SelectedSizeList"] = sizes
	data["SelectedCategories"] = selectedCategories
	data["SelectedGenders"] = selectedGenders
	data["Categories"] = categoryCounts
	data["Colors"] = colorCounts
	data["SelectedColors"] = selectedColors
	data["SelectedSizes"] = selectedSizes
	data["FilterChips"] = chips
	data["ClearFiltersURL"] = clearURL

	c.Header("Content-Type", "text/html; charset=utf-8")
	if err := h.templates["shop"].ExecuteTemplate(c.Writer, "base.html", data); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func (h *PageHandler) Account(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	if err := h.templates["account"].ExecuteTemplate(c.Writer, "base.html", h.getUserData(c)); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func (h *PageHandler) LoginPage(c *gin.Context) {
	errMsg := c.Query("error")
	if errMsg == "invalid+credentials" {
		errMsg = "Invalid credentials"
	} else if errMsg == "invalid+input" {
		errMsg = "Invalid input"
	}
	data := h.getUserData(c)
	data["Error"] = errMsg

	c.Header("Content-Type", "text/html; charset=utf-8")
	if err := h.templates["login"].ExecuteTemplate(c.Writer, "base.html", data); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func (h *PageHandler) RegisterPage(c *gin.Context) {
	errMsg := c.Query("error")
	if errMsg == "email+exists" {
		errMsg = "Email already registered"
	} else if errMsg == "invalid+input" {
		errMsg = "Invalid input"
	}
	data := h.getUserData(c)
	data["Error"] = errMsg

	c.Header("Content-Type", "text/html; charset=utf-8")
	if err := h.templates["register"].ExecuteTemplate(c.Writer, "base.html", data); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func (h *PageHandler) AdminOrders(c *gin.Context) {
	orders, err := h.orderService.ListAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	data := h.getUserData(c)
	data["Orders"] = orders

	c.Header("Content-Type", "text/html; charset=utf-8")
	if err := h.templates["admin_orders"].ExecuteTemplate(c.Writer, "base.html", data); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func (h *PageHandler) AdminProducts(c *gin.Context) {
	products, err := h.productService.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	data := h.getUserData(c)
	data["Products"] = products

	c.Header("Content-Type", "text/html; charset=utf-8")
	if err := h.templates["admin_products"].ExecuteTemplate(c.Writer, "base.html", data); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func (h *PageHandler) AdminDashboard(c *gin.Context) {
	stats, err := h.analyticsService.GetDashboardStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	data := h.getUserData(c)
	data["Stats"] = stats

	c.Header("Content-Type", "text/html; charset=utf-8")
	if err := h.templates["admin_dashboard"].ExecuteTemplate(c.Writer, "base.html", data); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func (h *PageHandler) AdminUsers(c *gin.Context) {
	users, err := h.authService.GetAllUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	data := h.getUserData(c)
	data["Users"] = users

	c.Header("Content-Type", "text/html; charset=utf-8")
	if err := h.templates["admin_users"].ExecuteTemplate(c.Writer, "base.html", data); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func (h *PageHandler) AdminUserOrders(c *gin.Context) {
	userID := c.Param("userId")
	if userID == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	user, err := h.authService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	orders, err := h.orderService.ListByUser(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	data := h.getUserData(c)
	data["Orders"] = orders
	data["FilterUser"] = user

	c.Header("Content-Type", "text/html; charset=utf-8")
	if err := h.templates["admin_orders"].ExecuteTemplate(c.Writer, "base.html", data); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func (h *PageHandler) AdminAnalytics(c *gin.Context) {
	stats, err := h.analyticsService.GetDashboardStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	data := h.getUserData(c)
	data["Stats"] = stats

	c.Header("Content-Type", "text/html; charset=utf-8")
	if err := h.templates["admin_analytics"].ExecuteTemplate(c.Writer, "base.html", data); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func (h *PageHandler) AccountOrders(c *gin.Context) {
	userID, _ := c.Get("user_id")
	if userID == nil || userID == "" {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	orders, err := h.orderService.ListByUser(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	data := h.getUserData(c)
	data["Orders"] = orders

	c.Header("Content-Type", "text/html; charset=utf-8")
	if err := h.templates["account_orders"].ExecuteTemplate(c.Writer, "base.html", data); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func (h *PageHandler) Product(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	product, err := h.productService.GetByID(c.Request.Context(), id)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	data := h.getUserData(c)
	data["Product"] = product

	c.Header("Content-Type", "text/html; charset=utf-8")
	if err := h.templates["product"].ExecuteTemplate(c.Writer, "base.html", data); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func (h *PageHandler) Wishlist(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	if err := h.templates["wishlist"].ExecuteTemplate(c.Writer, "base.html", h.getUserData(c)); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func (h *PageHandler) Cart(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	if err := h.templates["cart"].ExecuteTemplate(c.Writer, "base.html", h.getUserData(c)); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func (h *PageHandler) Checkout(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	if err := h.templates["checkout"].ExecuteTemplate(c.Writer, "base.html", h.getUserData(c)); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func getStr(c *gin.Context, key string) string {
	if v, ok := c.Get(key); ok && v != nil {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func filterProducts(products []*models.Product, q string, categories, colors, sizes, genders []string) []*models.Product {
	var out []*models.Product
	for _, p := range products {
		if q != "" && !strings.Contains(strings.ToLower(p.Name), strings.ToLower(q)) && !strings.Contains(strings.ToLower(p.Category), strings.ToLower(q)) {
			continue
		}
		if len(categories) > 0 {
			found := false
			for _, c := range categories {
				if strings.EqualFold(p.Category, c) {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		if len(genders) > 0 {
			found := false
			for _, g := range genders {
				if isUniversalGender(g) {
					if strings.TrimSpace(p.Gender) == "" {
						found = true
						break
					}
					continue
				}
				if strings.EqualFold(p.Gender, g) {
					found = true
					break
				}
			}
			if !found {
				continue
			}
		}
		if len(colors) > 0 {
			found := false
			for _, col := range colors {
				for _, pc := range p.Colors {
					if strings.EqualFold(pc, col) {
						found = true
						break
					}
				}
				if found {
					break
				}
			}
			if !found {
				continue
			}
		}
		if len(sizes) > 0 {
			found := false
			for _, size := range sizes {
				for _, ps := range p.Sizes {
					if strings.EqualFold(ps, size) {
						found = true
						break
					}
				}
				if found {
					break
				}
			}
			if !found {
				continue
			}
		}
		out = append(out, p)
	}
	return out
}

func isUniversalGender(value string) bool {
	return strings.EqualFold(strings.TrimSpace(value), "universal")
}

func sortProducts(products []*models.Product, order string) []*models.Product {
	out := make([]*models.Product, len(products))
	copy(out, products)
	switch order {
	case "price_asc":
		sort.Slice(out, func(i, j int) bool { return out[i].Price < out[j].Price })
	case "price_desc":
		sort.Slice(out, func(i, j int) bool { return out[i].Price > out[j].Price })
	case "name":
		sort.Slice(out, func(i, j int) bool { return strings.ToLower(out[i].Name) < strings.ToLower(out[j].Name) })
	}
	return out
}

type FilterChip struct {
	Label string
	URL   string
}

type CategoryCount struct {
	Name  string
	Count int
}

type ColorCount struct {
	Name  string
	Count int
}

func toSelectionMap(values []string) map[string]bool {
	out := make(map[string]bool, len(values))
	for _, v := range values {
		if v != "" {
			out[v] = true
		}
	}
	return out
}

func buildCategoryCounts(products []*models.Product) []CategoryCount {
	counts := make(map[string]int)
	for _, p := range products {
		name := strings.TrimSpace(p.Category)
		if name == "" {
			continue
		}
		counts[name]++
	}
	out := make([]CategoryCount, 0, len(counts))
	for name, count := range counts {
		out = append(out, CategoryCount{Name: name, Count: count})
	}
	sort.Slice(out, func(i, j int) bool {
		return strings.ToLower(out[i].Name) < strings.ToLower(out[j].Name)
	})
	return out
}

func buildColorCounts(products []*models.Product) []ColorCount {
	counts := make(map[string]int)
	display := make(map[string]string)
	for _, p := range products {
		seen := make(map[string]struct{})
		for _, color := range p.Colors {
			name := strings.TrimSpace(color)
			if name == "" {
				continue
			}
			key := strings.ToLower(name)
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			counts[key]++
			if _, ok := display[key]; !ok {
				display[key] = name
			}
		}
	}
	out := make([]ColorCount, 0, len(counts))
	for key, count := range counts {
		out = append(out, ColorCount{Name: display[key], Count: count})
	}
	sort.Slice(out, func(i, j int) bool {
		return strings.ToLower(out[i].Name) < strings.ToLower(out[j].Name)
	})
	return out
}

func buildFilterChips(values url.Values) ([]FilterChip, string) {
	var chips []FilterChip
	for _, v := range values["category"] {
		if v == "" {
			continue
		}
		chips = append(chips, FilterChip{
			Label: "Category: " + v,
			URL:   buildShopURL(removeQueryValue(values, "category", v)),
		})
	}
	for _, v := range values["gender"] {
		if v == "" {
			continue
		}
		chips = append(chips, FilterChip{
			Label: "Gender: " + v,
			URL:   buildShopURL(removeQueryValue(values, "gender", v)),
		})
	}
	for _, v := range values["color"] {
		if v == "" {
			continue
		}
		chips = append(chips, FilterChip{
			Label: "Color: " + v,
			URL:   buildShopURL(removeQueryValue(values, "color", v)),
		})
	}
	for _, v := range values["size"] {
		if v == "" {
			continue
		}
		chips = append(chips, FilterChip{
			Label: "Size: " + v,
			URL:   buildShopURL(removeQueryValue(values, "size", v)),
		})
	}
	clearValues := cloneValues(values)
	clearValues.Del("category")
	clearValues.Del("gender")
	clearValues.Del("color")
	clearValues.Del("size")
	return chips, buildShopURL(clearValues)
}

func cloneValues(values url.Values) url.Values {
	out := url.Values{}
	for k, vals := range values {
		for _, v := range vals {
			out.Add(k, v)
		}
	}
	return out
}

func removeQueryValue(values url.Values, key, value string) url.Values {
	out := cloneValues(values)
	existing := out[key]
	out.Del(key)
	for _, v := range existing {
		if v != value {
			out.Add(key, v)
		}
	}
	return out
}

func buildShopURL(values url.Values) string {
	if len(values) == 0 {
		return "/shop"
	}
	return "/shop?" + values.Encode()
}
