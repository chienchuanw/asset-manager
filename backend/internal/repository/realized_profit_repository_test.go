package repository

import (
"database/sql"
"testing"
"time"

"github.com/chienchuanw/asset-manager/internal/models"
"github.com/stretchr/testify/assert"
"github.com/stretchr/testify/require"
"github.com/stretchr/testify/suite"
)

// RealizedProfitRepositoryTestSuite 測試套件
type RealizedProfitRepositoryTestSuite struct {
suite.Suite
db              *sql.DB
repo            RealizedProfitRepository
transactionRepo TransactionRepository
}

// SetupSuite 在所有測試開始前執行一次
func (suite *RealizedProfitRepositoryTestSuite) SetupSuite() {
db, err := setupTestDB()
if err != nil {
suite.T().Fatalf("Failed to setup test database: %v", err)
}
suite.db = db
suite.repo = NewRealizedProfitRepository(db)
suite.transactionRepo = NewTransactionRepository(db)
}

// TearDownSuite 在所有測試結束後執行一次
func (suite *RealizedProfitRepositoryTestSuite) TearDownSuite() {
if suite.db != nil {
suite.db.Close()
}
}

// SetupTest 在每個測試開始前執行
func (suite *RealizedProfitRepositoryTestSuite) SetupTest() {
_, err := suite.db.Exec("TRUNCATE TABLE transactions CASCADE")
if err != nil {
suite.T().Fatalf("Failed to truncate tables: %v", err)
}
}

// TestCreate 測試建立已實現損益記錄
func (suite *RealizedProfitRepositoryTestSuite) TestCreate() {
sellDate := time.Date(2025, 10, 24, 0, 0, 0, 0, time.UTC)
fee := 28.0
transaction, err := suite.transactionRepo.Create(&models.CreateTransactionInput{
Date:            sellDate,
AssetType:       models.AssetTypeTWStock,
Symbol:          "2330",
Name:            "台積電",
TransactionType: models.TransactionTypeSell,
Quantity:        100,
Price:           620,
Amount:          62000,
Fee:             &fee,
Currency:        "TWD",
})
require.NoError(suite.T(), err)

input := &models.CreateRealizedProfitInput{
TransactionID: transaction.ID.String(),
Symbol:        "2330",
AssetType:     models.AssetTypeTWStock,
SellDate:      sellDate,
Quantity:      100,
SellPrice:     620,
SellAmount:    62000,
SellFee:       28,
CostBasis:     50000,
Currency:      "TWD",
}

result, err := suite.repo.Create(input)

require.NoError(suite.T(), err)
assert.NotEmpty(suite.T(), result.ID)
assert.Equal(suite.T(), transaction.ID.String(), result.TransactionID)
assert.Equal(suite.T(), "2330", result.Symbol)
assert.Equal(suite.T(), models.AssetTypeTWStock, result.AssetType)
assert.Equal(suite.T(), 100.0, result.Quantity)
assert.Equal(suite.T(), 620.0, result.SellPrice)
assert.Equal(suite.T(), 62000.0, result.SellAmount)
assert.Equal(suite.T(), 28.0, result.SellFee)
assert.Equal(suite.T(), 50000.0, result.CostBasis)

expectedPL := 11972.0
assert.Equal(suite.T(), expectedPL, result.RealizedPL)

expectedPLPct := 23.944
assert.InDelta(suite.T(), expectedPLPct, result.RealizedPLPct, 0.001)

assert.Equal(suite.T(), "TWD", result.Currency)
assert.NotZero(suite.T(), result.CreatedAt)
assert.NotZero(suite.T(), result.UpdatedAt)
}

// TestGetByTransactionID 測試根據交易 ID 取得已實現損益
func (suite *RealizedProfitRepositoryTestSuite) TestGetByTransactionID() {
sellDate := time.Date(2025, 10, 24, 0, 0, 0, 0, time.UTC)
fee := 28.0
transaction, _ := suite.transactionRepo.Create(&models.CreateTransactionInput{
Date:            sellDate,
AssetType:       models.AssetTypeTWStock,
Symbol:          "2330",
Name:            "台積電",
TransactionType: models.TransactionTypeSell,
Quantity:        100,
Price:           620,
Amount:          62000,
Fee:             &fee,
Currency:        "TWD",
})

created, _ := suite.repo.Create(&models.CreateRealizedProfitInput{
TransactionID: transaction.ID.String(),
Symbol:        "2330",
AssetType:     models.AssetTypeTWStock,
SellDate:      sellDate,
Quantity:      100,
SellPrice:     620,
SellAmount:    62000,
SellFee:       28,
CostBasis:     50000,
Currency:      "TWD",
})

result, err := suite.repo.GetByTransactionID(transaction.ID.String())

require.NoError(suite.T(), err)
assert.Equal(suite.T(), created.ID, result.ID)
assert.Equal(suite.T(), transaction.ID.String(), result.TransactionID)
assert.Equal(suite.T(), "2330", result.Symbol)
}

// TestGetByTransactionID_NotFound 測試取得不存在的記錄
func (suite *RealizedProfitRepositoryTestSuite) TestGetByTransactionID_NotFound() {
result, err := suite.repo.GetByTransactionID("non-existent-id")

assert.Error(suite.T(), err)
assert.Nil(suite.T(), result)
}

// TestRealizedProfitRepository 執行測試套件
func TestRealizedProfitRepository(t *testing.T) {
suite.Run(t, new(RealizedProfitRepositoryTestSuite))
}
