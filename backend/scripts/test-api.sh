#!/bin/bash

# API 測試腳本
# 用於快速測試 Transaction API 的各個端點

set -e

API_URL="http://localhost:8080"

echo "========================================="
echo "Testing Transaction API"
echo "========================================="
echo ""

# 顏色定義
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# 測試 Health Check
echo -e "${BLUE}1. Testing Health Check...${NC}"
curl -s $API_URL/health | jq .
echo -e "${GREEN}✓ Health check passed${NC}"
echo ""

# 建立交易記錄
echo -e "${BLUE}2. Creating a new transaction...${NC}"
RESPONSE=$(curl -s -X POST $API_URL/api/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "date": "2025-10-22T00:00:00Z",
    "asset_type": "tw-stock",
    "symbol": "2330",
    "name": "台積電",
    "type": "buy",
    "quantity": 10,
    "price": 620,
    "amount": 6200,
    "fee": 28,
    "note": "定期定額買入"
  }')

echo $RESPONSE | jq .

# 提取 ID
TRANSACTION_ID=$(echo $RESPONSE | jq -r '.data.id')

if [ "$TRANSACTION_ID" != "null" ] && [ -n "$TRANSACTION_ID" ]; then
    echo -e "${GREEN}✓ Transaction created with ID: $TRANSACTION_ID${NC}"
else
    echo -e "${RED}✗ Failed to create transaction${NC}"
    exit 1
fi
echo ""

# 建立第二筆交易記錄
echo -e "${BLUE}3. Creating another transaction...${NC}"
curl -s -X POST $API_URL/api/transactions \
  -H "Content-Type: application/json" \
  -d '{
    "date": "2025-10-21T00:00:00Z",
    "asset_type": "crypto",
    "symbol": "ETH",
    "name": "Ethereum",
    "type": "buy",
    "quantity": 2,
    "price": 50000,
    "amount": 100000,
    "fee": 100
  }' | jq .
echo -e "${GREEN}✓ Second transaction created${NC}"
echo ""

# 取得所有交易記錄
echo -e "${BLUE}4. Getting all transactions...${NC}"
curl -s $API_URL/api/transactions | jq .
echo -e "${GREEN}✓ Retrieved all transactions${NC}"
echo ""

# 取得單筆交易記錄
echo -e "${BLUE}5. Getting transaction by ID...${NC}"
curl -s $API_URL/api/transactions/$TRANSACTION_ID | jq .
echo -e "${GREEN}✓ Retrieved transaction by ID${NC}"
echo ""

# 更新交易記錄
echo -e "${BLUE}6. Updating transaction...${NC}"
curl -s -X PUT $API_URL/api/transactions/$TRANSACTION_ID \
  -H "Content-Type: application/json" \
  -d '{
    "quantity": 20,
    "price": 630,
    "amount": 12600
  }' | jq .
echo -e "${GREEN}✓ Transaction updated${NC}"
echo ""

# 使用篩選條件查詢
echo -e "${BLUE}7. Filtering transactions (tw-stock only)...${NC}"
curl -s "$API_URL/api/transactions?asset_type=tw-stock" | jq .
echo -e "${GREEN}✓ Filtered transactions retrieved${NC}"
echo ""

# 刪除交易記錄
echo -e "${BLUE}8. Deleting transaction...${NC}"
curl -s -X DELETE $API_URL/api/transactions/$TRANSACTION_ID | jq .
echo -e "${GREEN}✓ Transaction deleted${NC}"
echo ""

# 驗證刪除
echo -e "${BLUE}9. Verifying deletion (should return error)...${NC}"
curl -s $API_URL/api/transactions/$TRANSACTION_ID | jq .
echo -e "${GREEN}✓ Deletion verified${NC}"
echo ""

echo "========================================="
echo -e "${GREEN}All tests passed! 🎉${NC}"
echo "========================================="

