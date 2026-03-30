#!/bin/bash

BASE="http://localhost:9000"
PASS=0
FAIL=0

check() {
  local desc="$1"
  local expected="$2"
  local actual="$3"

  if echo "$actual" | grep -q "$expected"; then
    echo "  ✅ $desc"
    ((PASS++))
  else
    echo "  ❌ $desc"
    echo "     Oczekiwano: $expected"
    echo "     Otrzymano:  $actual"
    ((FAIL++))
  fi
}

echo "=== PRODUKTY ==="

echo "GET /products (lista)"
RES=$(curl -s "$BASE/products")
check "Zwraca Laptop" "Laptop" "$RES"

echo "GET /products/1 (po ID)"
RES=$(curl -s "$BASE/products/1")
check "Zwraca produkt o id 1" "Laptop" "$RES"

echo "GET /products/999 (nieistniejący)"
RES=$(curl -s "$BASE/products/999")
check "Zwraca NotFound" "Nie znaleziono" "$RES"

echo "POST /products (dodaj)"
RES=$(curl -s -X POST "$BASE/products" -H "Content-Type: application/json" -d '{"id":3,"name":"Klawiatura","price":250.00}')
check "Tworzy Klawiatura" "Klawiatura" "$RES"

echo "PUT /products/3 (aktualizuj)"
RES=$(curl -s -X PUT "$BASE/products/3" -H "Content-Type: application/json" -d '{"id":3,"name":"Klawiatura RGB","price":350.00}')
check "Aktualizuje na Klawiatura RGB" "Klawiatura RGB" "$RES"

echo "DELETE /products/3 (usuń)"
RES=$(curl -s -X DELETE "$BASE/products/3")
check "Usuwa produkt" "Usunięto" "$RES"

echo ""
echo "=== KATEGORIE ==="

echo "GET /categories (lista)"
RES=$(curl -s "$BASE/categories")
check "Zwraca Elektronika" "Elektronika" "$RES"

echo "GET /categories/1 (po ID)"
RES=$(curl -s "$BASE/categories/1")
check "Zwraca kategorię o id 1" "Elektronika" "$RES"

echo "GET /categories/999 (nieistniejący)"
RES=$(curl -s "$BASE/categories/999")
check "Zwraca NotFound" "Nie znaleziono" "$RES"

echo "POST /categories (dodaj)"
RES=$(curl -s -X POST "$BASE/categories" -H "Content-Type: application/json" -d '{"id":3,"name":"Odzież"}')
check "Tworzy Odzież" "Odzież" "$RES"

echo "PUT /categories/3 (aktualizuj)"
RES=$(curl -s -X PUT "$BASE/categories/3" -H "Content-Type: application/json" -d '{"id":3,"name":"Ubrania"}')
check "Aktualizuje na Ubrania" "Ubrania" "$RES"

echo "DELETE /categories/3 (usuń)"
RES=$(curl -s -X DELETE "$BASE/categories/3")
check "Usuwa kategorię" "Usunięto" "$RES"

echo ""
echo "=== KOSZYK ==="

echo "GET /cart (lista)"
RES=$(curl -s "$BASE/cart")
check "Zwraca elementy koszyka" "productId" "$RES"

echo "GET /cart/1 (po ID)"
RES=$(curl -s "$BASE/cart/1")
check "Zwraca element o id 1" "productId" "$RES"

echo "GET /cart/999 (nieistniejący)"
RES=$(curl -s "$BASE/cart/999")
check "Zwraca NotFound" "Nie znaleziono" "$RES"

echo "POST /cart (dodaj)"
RES=$(curl -s -X POST "$BASE/cart" -H "Content-Type: application/json" -d '{"id":2,"productId":2,"quantity":3}')
check "Tworzy element koszyka" "\"id\":2" "$RES"

echo "PUT /cart/2 (aktualizuj)"
RES=$(curl -s -X PUT "$BASE/cart/2" -H "Content-Type: application/json" -d '{"id":2,"productId":2,"quantity":5}')
check "Aktualizuje ilość na 5" "\"quantity\":5" "$RES"

echo "DELETE /cart/2 (usuń)"
RES=$(curl -s -X DELETE "$BASE/cart/2")
check "Usuwa element z koszyka" "Usunięto" "$RES"

echo ""
echo "================================"
echo "Wynik: $PASS ✅  |  $FAIL ❌"
echo "================================"
