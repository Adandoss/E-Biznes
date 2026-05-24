package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v5"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"sklep/handlers"
	"sklep/models"
)

// helpers 

func testToken(t *testing.T) string {
	t.Helper()
	claims := jwt.MapClaims{
		"user_id": float64(1),
		"email":   "test@test.com",
		"name":    "Test",
		"exp":     time.Now().Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := token.SignedString([]byte("default-dev-secret-change-in-production"))
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func newTestEnv(t *testing.T) *echo.Echo {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatal(err)
	}
	if err := db.AutoMigrate(&models.Product{}, &models.Cart{}, &models.CartItem{}, &models.Payment{}); err != nil {
		t.Fatal(err)
	}

	e := echo.New()
	ph := &handlers.ProductHandler{DB: db}
	ch := &handlers.CartHandler{DB: db}
	pay := &handlers.PaymentHandler{DB: db}

	e.POST("/products", ph.CreateProduct)
	e.GET("/products", ph.GetProducts)
	e.GET("/products/:id", ph.GetProduct)
	e.PUT("/products/:id", ph.UpdateProduct)
	e.DELETE("/products/:id", ph.DeleteProduct)

	carts := e.Group("/carts", handlers.JWTMiddleware)
	carts.POST("", ch.CreateCart)
	carts.GET("/:id", ch.GetCart)
	carts.POST("/:id/items", ch.AddItem)
	carts.DELETE("/:id/items/:itemId", ch.RemoveItem)
	carts.DELETE("/:id", ch.DeleteCart)

	payments := e.Group("/payments", handlers.JWTMiddleware)
	payments.GET("", pay.GetPayments)
	payments.POST("", pay.CreatePayment)

	return e
}

func doJSON(e *echo.Echo, method, path string, body any, token ...string) *httptest.ResponseRecorder {
	var req *http.Request
	if body != nil {
		b, _ := json.Marshal(body)
		req = httptest.NewRequest(method, path, bytes.NewReader(b))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	if len(token) > 0 && token[0] != "" {
		req.Header.Set("Authorization", "Bearer "+token[0])
	}
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec
}

func doRaw(e *echo.Echo, method, path, rawBody string, token ...string) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, bytes.NewReader([]byte(rawBody)))
	req.Header.Set("Content-Type", "application/json")
	if len(token) > 0 && token[0] != "" {
		req.Header.Set("Authorization", "Bearer "+token[0])
	}
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	return rec
}

func bodyMap(rec *httptest.ResponseRecorder) map[string]any {
	var m map[string]any
	_ = json.Unmarshal(rec.Body.Bytes(), &m)
	return m
}

func bodySlice(rec *httptest.ResponseRecorder) []any {
	var a []any
	_ = json.Unmarshal(rec.Body.Bytes(), &a)
	return a
}

func assertStatus(t *testing.T, rec *httptest.ResponseRecorder, want int) {
	t.Helper()
	if rec.Code != want {
		t.Fatalf("status: want %d, got %d", want, rec.Code)
	}
}

func assertField(t *testing.T, m map[string]any, key string, want any) {
	t.Helper()
	got := m[key]
	if got != want {
		t.Errorf("%s: want %v, got %v", key, want, got)
	}
}

func assertFieldGT(t *testing.T, m map[string]any, key string, min float64) {
	t.Helper()
	if m[key].(float64) <= min {
		t.Errorf("%s: want > %v, got %v", key, min, m[key])
	}
}

func product(name, desc string, price float64) map[string]any {
	return map[string]any{"name": name, "description": desc, "price": price}
}

func cartItem(productID, qty int) map[string]any {
	return map[string]any{"product_id": productID, "quantity": qty}
}

func payment(amount float64, status string) map[string]any {
	return map[string]any{"amount": amount, "status": status}
}

// PRODUCTS 

func TestCreateProduct(t *testing.T) {
	e := newTestEnv(t)
	rec := doJSON(e, "POST", "/products", product("Testowy", "Opis", 19.99))
	body := bodyMap(rec)

	assertStatus(t, rec, 201)
	assertField(t, body, "name", "Testowy")
	assertField(t, body, "description", "Opis")
	assertField(t, body, "price", 19.99)
	assertFieldGT(t, body, "ID", 0)
}

func TestCreateProductInvalidBody(t *testing.T) {
	e := newTestEnv(t)
	assertStatus(t, doRaw(e, "POST", "/products", "invalid"), 400)
}

func TestGetProductsEmpty(t *testing.T) {
	e := newTestEnv(t)
	rec := doJSON(e, "GET", "/products", nil)
	assertStatus(t, rec, 200)
	if len(bodySlice(rec)) != 0 {
		t.Error("expected empty list")
	}
}

func TestGetProductsWithData(t *testing.T) {
	e := newTestEnv(t)
	doJSON(e, "POST", "/products", product("P1", "D1", 10.0))
	doJSON(e, "POST", "/products", product("P2", "D2", 20.0))

	rec := doJSON(e, "GET", "/products", nil)
	arr := bodySlice(rec)

	assertStatus(t, rec, 200)
	if len(arr) != 2 {
		t.Errorf("want 2 products, got %d", len(arr))
	}
	first := arr[0].(map[string]any)
	assertField(t, first, "name", "P2") // sorted price DESC
}

func TestGetProduct(t *testing.T) {
	e := newTestEnv(t)
	doJSON(e, "POST", "/products", product("Single", "D", 5.0))

	rec := doJSON(e, "GET", "/products/1", nil)
	body := bodyMap(rec)

	assertStatus(t, rec, 200)
	assertField(t, body, "name", "Single")
	assertField(t, body, "price", 5.0)
}

func TestGetProductNotFound(t *testing.T) {
	e := newTestEnv(t)
	assertStatus(t, doJSON(e, "GET", "/products/999", nil), 404)
}

func TestUpdateProduct(t *testing.T) {
	e := newTestEnv(t)
	doJSON(e, "POST", "/products", product("Old", "D", 1.0))

	rec := doJSON(e, "PUT", "/products/1", map[string]any{"name": "New", "price": 99.0})
	body := bodyMap(rec)

	assertStatus(t, rec, 200)
	assertField(t, body, "name", "New")
	assertField(t, body, "price", 99.0)
}

func TestUpdateProductNotFound(t *testing.T) {
	e := newTestEnv(t)
	assertStatus(t, doJSON(e, "PUT", "/products/999", map[string]any{"name": "X", "price": 1}), 404)
}

func TestDeleteProduct(t *testing.T) {
	e := newTestEnv(t)
	doJSON(e, "POST", "/products", product("ToDelete", "D", 1.0))
	assertStatus(t, doJSON(e, "DELETE", "/products/1", nil), 204)
}

func TestDeleteProductNotFound(t *testing.T) {
	e := newTestEnv(t)
	rec := doJSON(e, "DELETE", "/products/999", nil)
	assertStatus(t, rec, 404)
	if bodyMap(rec)["error"] == nil {
		t.Error("expected error in body")
	}
}

func TestGetProductAfterDelete(t *testing.T) {
	e := newTestEnv(t)
	doJSON(e, "POST", "/products", product("Gone", "D", 1.0))
	doJSON(e, "DELETE", "/products/1", nil)
	assertStatus(t, doJSON(e, "GET", "/products/1", nil), 404)
}

// CARTS 

func TestCreateCart(t *testing.T) {
	e := newTestEnv(t)
	tk := testToken(t)
	rec := doJSON(e, "POST", "/carts", nil, tk)
	body := bodyMap(rec)

	assertStatus(t, rec, 201)
	assertField(t, body, "status", "aktywny")
	assertFieldGT(t, body, "ID", 0)
}

func TestGetCart(t *testing.T) {
	e := newTestEnv(t)
	tk := testToken(t)
	doJSON(e, "POST", "/carts", nil, tk)

	rec := doJSON(e, "GET", "/carts/1", nil, tk)
	assertStatus(t, rec, 200)
	assertField(t, bodyMap(rec), "status", "aktywny")
}

func TestGetCartNotFound(t *testing.T) {
	e := newTestEnv(t)
	tk := testToken(t)
	assertStatus(t, doJSON(e, "GET", "/carts/999", nil, tk), 404)
}

func TestAddItemToCart(t *testing.T) {
	e := newTestEnv(t)
	tk := testToken(t)
	doJSON(e, "POST", "/carts", nil, tk)
	doJSON(e, "POST", "/products", product("Item", "D", 10.0))

	rec := doJSON(e, "POST", "/carts/1/items", cartItem(1, 3), tk)
	body := bodyMap(rec)

	assertStatus(t, rec, 201)
	assertField(t, body, "quantity", float64(3))
	assertField(t, body, "product_id", float64(1))
}

func TestAddItemCartNotFound(t *testing.T) {
	e := newTestEnv(t)
	tk := testToken(t)
	assertStatus(t, doJSON(e, "POST", "/carts/999/items", cartItem(1, 1), tk), 404)
}

func TestAddItemProductNotFound(t *testing.T) {
	e := newTestEnv(t)
	tk := testToken(t)
	doJSON(e, "POST", "/carts", nil, tk)
	assertStatus(t, doJSON(e, "POST", "/carts/1/items", cartItem(999, 1), tk), 404)
}

func TestGetCartWithItems(t *testing.T) {
	e := newTestEnv(t)
	tk := testToken(t)
	doJSON(e, "POST", "/carts", nil, tk)
	doJSON(e, "POST", "/products", product("X", "D", 5.0))
	doJSON(e, "POST", "/carts/1/items", cartItem(1, 1), tk)

	rec := doJSON(e, "GET", "/carts/1", nil, tk)
	assertStatus(t, rec, 200)

	items := bodyMap(rec)["items"].([]any)
	if len(items) != 1 {
		t.Errorf("want 1 item, got %d", len(items))
	}
}

func TestRemoveItem(t *testing.T) {
	e := newTestEnv(t)
	tk := testToken(t)
	doJSON(e, "POST", "/carts", nil, tk)
	doJSON(e, "POST", "/products", product("X", "D", 5.0))
	doJSON(e, "POST", "/carts/1/items", cartItem(1, 1), tk)

	assertStatus(t, doJSON(e, "DELETE", "/carts/1/items/1", nil, tk), 204)
}

func TestRemoveItemNotFound(t *testing.T) {
	e := newTestEnv(t)
	tk := testToken(t)
	doJSON(e, "POST", "/carts", nil, tk)
	assertStatus(t, doJSON(e, "DELETE", "/carts/1/items/999", nil, tk), 404)
}

func TestDeleteCart(t *testing.T) {
	e := newTestEnv(t)
	tk := testToken(t)
	doJSON(e, "POST", "/carts", nil, tk)
	assertStatus(t, doJSON(e, "DELETE", "/carts/1", nil, tk), 204)
}

func TestDeleteCartNotFound(t *testing.T) {
	e := newTestEnv(t)
	tk := testToken(t)
	assertStatus(t, doJSON(e, "DELETE", "/carts/999", nil, tk), 404)
}

func TestGetCartAfterDelete(t *testing.T) {
	e := newTestEnv(t)
	tk := testToken(t)
	doJSON(e, "POST", "/carts", nil, tk)
	doJSON(e, "DELETE", "/carts/1", nil, tk)
	assertStatus(t, doJSON(e, "GET", "/carts/1", nil, tk), 404)
}

// PAYMENTS 

func TestCreatePayment(t *testing.T) {
	e := newTestEnv(t)
	tk := testToken(t)
	rec := doJSON(e, "POST", "/payments", payment(250.0, "PAID"), tk)
	body := bodyMap(rec)

	assertStatus(t, rec, 201)
	assertField(t, body, "amount", 250.0)
	assertField(t, body, "status", "PAID")
}

func TestCreatePaymentInvalidBody(t *testing.T) {
	e := newTestEnv(t)
	tk := testToken(t)
	rec := doRaw(e, "POST", "/payments", "bad", tk)
	if rec.Code != 400 && rec.Code != 500 {
		t.Fatalf("want 400 or 500, got %d", rec.Code)
	}
}

func TestGetPaymentsEmpty(t *testing.T) {
	e := newTestEnv(t)
	tk := testToken(t)
	rec := doJSON(e, "GET", "/payments", nil, tk)
	assertStatus(t, rec, 200)
	if len(bodySlice(rec)) != 0 {
		t.Error("expected empty list")
	}
}

func TestGetPaymentsWithData(t *testing.T) {
	e := newTestEnv(t)
	tk := testToken(t)
	doJSON(e, "POST", "/payments", payment(50.0, "PAID"), tk)
	doJSON(e, "POST", "/payments", payment(75.0, "PENDING"), tk)

	rec := doJSON(e, "GET", "/payments", nil, tk)
	assertStatus(t, rec, 200)
	if len(bodySlice(rec)) != 2 {
		t.Error("want 2 payments")
	}
}

func TestMultipleCartItems(t *testing.T) {
	e := newTestEnv(t)
	tk := testToken(t)
	doJSON(e, "POST", "/carts", nil, tk)
	doJSON(e, "POST", "/products", product("A", "D", 10.0))
	doJSON(e, "POST", "/products", product("B", "D", 20.0))
	doJSON(e, "POST", "/carts/1/items", cartItem(1, 1), tk)
	doJSON(e, "POST", "/carts/1/items", cartItem(2, 2), tk)

	rec := doJSON(e, "GET", "/carts/1", nil, tk)
	body := bodyMap(rec)

	assertStatus(t, rec, 200)
	items := body["items"].([]any)
	if len(items) != 2 {
		t.Errorf("want 2 items, got %d", len(items))
	}
	second := items[1].(map[string]any)
	assertField(t, second, "quantity", float64(2))
}
