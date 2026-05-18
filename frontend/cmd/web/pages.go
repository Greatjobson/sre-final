package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/Tedra-ez/AdvancedProgramming_Final/frontend/internal/config"
	"github.com/Tedra-ez/AdvancedProgramming_Final/frontend/internal/models"
	"github.com/Tedra-ez/AdvancedProgramming_Final/frontend/internal/services"
	"github.com/gin-gonic/gin"
)

type serviceClient struct {
	httpClient *http.Client

	authURL      string
	productURL   string
	orderURL     string
	userURL      string
	analyticsURL string
}

func newServiceClient(cfg *config.Config) *serviceClient {
	return &serviceClient{
		httpClient:   &http.Client{Timeout: 10 * time.Second},
		authURL:      strings.TrimRight(cfg.AuthServiceURL, "/"),
		productURL:   strings.TrimRight(cfg.ProductServiceURL, "/"),
		orderURL:     strings.TrimRight(cfg.OrderServiceURL, "/"),
		userURL:      strings.TrimRight(cfg.UserServiceURL, "/"),
		analyticsURL: strings.TrimRight(cfg.AnalyticsURL, "/"),
	}
}

func (c *serviceClient) products(ctx *gin.Context) ([]*models.Product, error) {
	var products []*models.Product
	if err := c.getJSON(ctx, c.productURL, "/api/product", "", &products); err != nil {
		return nil, err
	}
	if products == nil {
		products = []*models.Product{}
	}
	return products, nil
}

func (c *serviceClient) product(ctx *gin.Context, id string) (*models.Product, error) {
	var product models.Product
	if err := c.getJSON(ctx, c.productURL, "/api/product/"+url.PathEscape(id), "", &product); err != nil {
		return nil, err
	}
	return &product, nil
}

func (c *serviceClient) users(ctx *gin.Context) ([]*models.User, error) {
	var users []*models.User
	if err := c.getJSON(ctx, c.userURL, "/api/users", "", &users); err != nil {
		return nil, err
	}
	if users == nil {
		users = []*models.User{}
	}
	return users, nil
}

func (c *serviceClient) user(ctx *gin.Context, id string) (*models.User, error) {
	var user models.User
	if err := c.getJSON(ctx, c.userURL, "/api/users/"+url.PathEscape(id), "", &user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (c *serviceClient) orders(ctx *gin.Context, query string) ([]*models.Order, error) {
	var orders []*models.Order
	if err := c.getJSON(ctx, c.orderURL, "/orders", query, &orders); err != nil {
		return nil, err
	}
	if orders == nil {
		orders = []*models.Order{}
	}
	return orders, nil
}

func (c *serviceClient) stats(ctx *gin.Context) (*services.DashboardStats, error) {
	var stats services.DashboardStats
	if err := c.getJSON(ctx, c.analyticsURL, "/api/analytics/stats", "", &stats); err != nil {
		return nil, err
	}
	if stats.OrdersByStatus == nil {
		stats.OrdersByStatus = map[string]int{}
	}
	return &stats, nil
}

func (c *serviceClient) getJSON(ctx *gin.Context, baseURL, path, rawQuery string, out interface{}) error {
	endpoint := baseURL + path
	if rawQuery != "" {
		endpoint += "?" + rawQuery
	}
	req, err := http.NewRequestWithContext(ctx.Request.Context(), http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}
	for _, cookie := range ctx.Request.Cookies() {
		req.AddCookie(cookie)
	}
	if auth := ctx.GetHeader("Authorization"); auth != "" {
		req.Header.Set("Authorization", auth)
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		var body struct {
			Error string `json:"error"`
		}
		_ = json.NewDecoder(res.Body).Decode(&body)
		if body.Error == "" {
			body.Error = res.Status
		}
		return errors.New(body.Error)
	}
	return json.NewDecoder(res.Body).Decode(out)
}

func (p *pageServer) render(c *gin.Context, page string, data gin.H) {
	tmpl, ok := p.templates[page]
	if !ok {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.Header("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(c.Writer, "base.html", data); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
}

func (p *pageServer) userData(c *gin.Context) gin.H {
	data := gin.H{}
	if id, ok := c.Get("user_id"); ok && id != "" {
		data["User"] = map[string]string{
			"id":    asString(id),
			"role":  contextString(c, "user_role"),
			"email": contextString(c, "user_email"),
			"name":  contextString(c, "user_name"),
		}
	}
	return data
}

func (p *pageServer) Index(c *gin.Context) {
	p.render(c, "index", p.userData(c))
}

func (p *pageServer) Shop(c *gin.Context) {
	products, err := p.client.products(c)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}

	q := c.Query("q")
	order := c.Query("sort")
	if order == "" {
		order = "recommended"
	}
	categories := c.QueryArray("category")
	genders := c.QueryArray("gender")
	colors := c.QueryArray("color")
	sizes := c.QueryArray("size")
	filtered := filterProducts(products, q, categories, colors, sizes, genders)
	ordered := sortProducts(filtered, order)
	chips, clearURL := buildFilterChips(c.Request.URL.Query())

	data := p.userData(c)
	data["Products"] = ordered
	data["SearchQuery"] = q
	data["Sort"] = order
	data["ShowSidebar"] = true
	data["SelectedCategoryList"] = categories
	data["SelectedGenderList"] = genders
	data["SelectedColorList"] = colors
	data["SelectedSizeList"] = sizes
	data["SelectedCategories"] = toSelectionMap(categories)
	data["SelectedGenders"] = toSelectionMap(genders)
	data["Categories"] = buildCategoryCounts(products)
	data["Colors"] = buildColorCounts(products)
	data["SelectedColors"] = toSelectionMap(colors)
	data["SelectedSizes"] = toSelectionMap(sizes)
	data["FilterChips"] = chips
	data["ClearFiltersURL"] = clearURL

	p.render(c, "shop", data)
}

func (p *pageServer) Product(c *gin.Context) {
	product, err := p.client.product(c, c.Param("id"))
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	data := p.userData(c)
	data["Product"] = product
	p.render(c, "product", data)
}

func (p *pageServer) Account(c *gin.Context) {
	p.render(c, "account", p.userData(c))
}

func (p *pageServer) Wishlist(c *gin.Context) {
	p.render(c, "wishlist", p.userData(c))
}

func (p *pageServer) Cart(c *gin.Context) {
	p.render(c, "cart", p.userData(c))
}

func (p *pageServer) Checkout(c *gin.Context) {
	p.render(c, "checkout", p.userData(c))
}

func (p *pageServer) Login(c *gin.Context) {
	data := p.userData(c)
	data["Error"] = readableAuthError(c.Query("error"))
	p.render(c, "login", data)
}

func (p *pageServer) Register(c *gin.Context) {
	data := p.userData(c)
	data["Error"] = readableAuthError(c.Query("error"))
	p.render(c, "register", data)
}

func (p *pageServer) AccountOrders(c *gin.Context) {
	userID := contextString(c, "user_id")
	if userID == "" {
		c.Redirect(http.StatusFound, "/login")
		return
	}
	values := url.Values{}
	values.Set("user_id", userID)
	orders, err := p.client.orders(c, values.Encode())
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	data := p.userData(c)
	data["Orders"] = orders
	p.render(c, "account_orders", data)
}

func (p *pageServer) AdminDashboard(c *gin.Context) {
	p.renderAdminStats(c, "admin_dashboard")
}

func (p *pageServer) AdminAnalytics(c *gin.Context) {
	p.renderAdminStats(c, "admin_analytics")
}

func (p *pageServer) renderAdminStats(c *gin.Context, page string) {
	stats, err := p.client.stats(c)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	data := p.userData(c)
	data["Stats"] = stats
	p.render(c, page, data)
}

func (p *pageServer) AdminOrders(c *gin.Context) {
	orders, err := p.client.orders(c, "")
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	data := p.userData(c)
	data["Orders"] = orders
	p.render(c, "admin_orders", data)
}

func (p *pageServer) AdminProducts(c *gin.Context) {
	products, err := p.client.products(c)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	data := p.userData(c)
	data["Products"] = products
	p.render(c, "admin_products", data)
}

func (p *pageServer) AdminUsers(c *gin.Context) {
	users, err := p.client.users(c)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	data := p.userData(c)
	data["Users"] = users
	p.render(c, "admin_users", data)
}

func (p *pageServer) AdminUserOrders(c *gin.Context) {
	userID := c.Param("userId")
	user, err := p.client.user(c, userID)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	values := url.Values{}
	values.Set("user_id", userID)
	orders, err := p.client.orders(c, values.Encode())
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	data := p.userData(c)
	data["Orders"] = orders
	data["FilterUser"] = user
	p.render(c, "admin_orders", data)
}

func contextString(c *gin.Context, key string) string {
	if value, ok := c.Get(key); ok && value != nil {
		return asString(value)
	}
	return ""
}

func asString(value interface{}) string {
	if s, ok := value.(string); ok {
		return s
	}
	return ""
}

func readableAuthError(value string) string {
	switch strings.ReplaceAll(value, "+", " ") {
	case "invalid credentials":
		return "Invalid credentials"
	case "invalid input":
		return "Invalid input"
	case "email exists":
		return "Email already registered"
	default:
		return value
	}
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

func filterProducts(products []*models.Product, q string, categories, colors, sizes, genders []string) []*models.Product {
	var out []*models.Product
	for _, product := range products {
		if q != "" && !strings.Contains(strings.ToLower(product.Name), strings.ToLower(q)) && !strings.Contains(strings.ToLower(product.Category), strings.ToLower(q)) {
			continue
		}
		if !matchesString(product.Category, categories, false) {
			continue
		}
		if !matchesGender(product.Gender, genders) {
			continue
		}
		if !matchesAny(product.Colors, colors) {
			continue
		}
		if !matchesAny(product.Sizes, sizes) {
			continue
		}
		out = append(out, product)
	}
	return out
}

func matchesString(value string, filters []string, emptyMatches bool) bool {
	if len(filters) == 0 {
		return true
	}
	if strings.TrimSpace(value) == "" {
		return emptyMatches
	}
	for _, filter := range filters {
		if strings.EqualFold(value, filter) {
			return true
		}
	}
	return false
}

func matchesGender(value string, filters []string) bool {
	if len(filters) == 0 {
		return true
	}
	for _, filter := range filters {
		if strings.EqualFold(strings.TrimSpace(filter), "universal") && strings.TrimSpace(value) == "" {
			return true
		}
		if strings.EqualFold(value, filter) {
			return true
		}
	}
	return false
}

func matchesAny(values, filters []string) bool {
	if len(filters) == 0 {
		return true
	}
	for _, filter := range filters {
		for _, value := range values {
			if strings.EqualFold(value, filter) {
				return true
			}
		}
	}
	return false
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

func toSelectionMap(values []string) map[string]bool {
	out := make(map[string]bool, len(values))
	for _, value := range values {
		if value != "" {
			out[value] = true
		}
	}
	return out
}

func buildCategoryCounts(products []*models.Product) []CategoryCount {
	counts := make(map[string]int)
	for _, product := range products {
		name := strings.TrimSpace(product.Category)
		if name != "" {
			counts[name]++
		}
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
	for _, product := range products {
		seen := make(map[string]struct{})
		for _, color := range product.Colors {
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
	for _, key := range []string{"category", "gender", "color", "size"} {
		for _, value := range values[key] {
			if value == "" {
				continue
			}
			chips = append(chips, FilterChip{
				Label: strings.Title(key) + ": " + value,
				URL:   buildShopURL(removeQueryValue(values, key, value)),
			})
		}
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
	for key, vals := range values {
		for _, value := range vals {
			out.Add(key, value)
		}
	}
	return out
}

func removeQueryValue(values url.Values, key, value string) url.Values {
	out := cloneValues(values)
	existing := out[key]
	out.Del(key)
	for _, current := range existing {
		if current != value {
			out.Add(key, current)
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
