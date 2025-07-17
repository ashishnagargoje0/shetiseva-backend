#!/bin/bash

echo "====== Logging in to fetch token... ======"

TOKEN=$(curl -s -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"ashu@admin.com", "password":"Asz@1212"}' | jq -r '.token')

if [[ -z "$TOKEN" || "$TOKEN" == "null" ]]; then
  echo "❌ Failed to fetch token"
  exit 1
fi

echo "✅ Token fetched"

# Test data
PRODUCT_ID="6872063ff56ace040fc335d5"
ORDER_ID="6872068f68b7fb176b9fa306"
SUPPORT_ID="6872063ff56ace040fc335da"
CATEGORY_SLUG="seeds"

echo "====== Starting API Tests ======"

# Function to run test
test_api() {
  METHOD=$1
  ENDPOINT=$2
  DATA=$3

  echo -e "\n🔸 $METHOD $ENDPOINT"
  if [ "$METHOD" == "GET" ] || [ "$METHOD" == "DELETE" ]; then
    curl -s -X $METHOD http://localhost:8080$ENDPOINT -H "Authorization: Bearer $TOKEN"
  else
    curl -s -X $METHOD http://localhost:8080$ENDPOINT \
      -H "Authorization: Bearer $TOKEN" \
      -H "Content-Type: application/json" \
      -d "$DATA"
  fi
}

# === CART ===
test_api POST "/cart/add" '{"productId":"'"$PRODUCT_ID"'", "quantity": 2}'
test_api GET "/cart"
test_api POST "/cart/update" '{"productId":"'"$PRODUCT_ID"'", "quantity": 3}'
test_api DELETE "/cart/remove" '{"productId":"'"$PRODUCT_ID"'"}'

# === ORDER ===
test_api POST "/order/create" '{"items":[{"productId":"'"$PRODUCT_ID"'", "quantity": 1}]}'
test_api GET "/order/$ORDER_ID"
test_api GET "/order/history"
test_api POST "/order/deliver" '{"orderId":"'"$ORDER_ID"'"}'

# === PAYMENT/INVOICE ===
test_api GET "/payment/methods"
test_api GET "/invoice/$ORDER_ID"

# === DELIVERY ===
test_api GET "/delivery/status/$ORDER_ID"

# === USER ===
test_api PUT "/user/profile/update" '{"name": "Ashish", "phone": "9876543210"}'
test_api POST "/user/language/set" '{"language":"mr"}'
test_api GET "/user/kyc/status"

# === PRODUCT ===
test_api GET "/products/filters"
test_api GET "/product/$PRODUCT_ID"
test_api GET "/category/$CATEGORY_SLUG/products"

# === WISHLIST ===
test_api POST "/wishlist/add" '{"productId":"'"$PRODUCT_ID"'"}'
test_api GET "/wishlist"
test_api DELETE "/wishlist/remove" '{"productId":"'"$PRODUCT_ID"'"}'

# === MARKETPLACE ===
test_api POST "/marketplace/ad" '{"title":"Farm Tools", "price":500, "location":"Pune"}'
test_api GET "/marketplace/ads"
test_api DELETE "/marketplace/ad/remove" '{"adId":"dummyAdId"}'

# === SUPPORT ===
test_api GET "/support/tickets"
test_api POST "/support/ticket" '{"subject":"Test Ticket", "message":"Testing system"}'
test_api GET "/support/status/$SUPPORT_ID"

# === REVIEW / FEEDBACK ===
test_api POST "/review/submit" '{"productId":"'"$PRODUCT_ID"'", "rating":4, "comment":"Very good quality seeds!"}'
test_api POST "/feedback/voice" '{"userId":"687204fc68b7fb176b9fa2ff", "audioURL":"https://cdn.shetiseva.in/feedback/voice123.mp3"}'

# === META ===
test_api GET "/meta/districts"
test_api GET "/meta/languages"

# === HEALTH ===
test_api GET "/health"

echo -e "\n====== ✅ All API Tests Finished ======"
