package repository

import (
	"database/sql"
	"testing"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// TransactionRepositoryTestSuite 測試套件
type TransactionRepositoryTestSuite struct {
	suite.Suite
	db   *sql.DB
	repo TransactionRepository
}

// SetupSuite 在所有測試開始前執行一次
func (suite *TransactionRepositoryTestSuite) SetupSuite() {
	// 這裡需要連接測試資料庫
	// 你需要設定測試資料庫的環境變數
	db, err := setupTestDB()
	if err != nil {
		suite.T().Fatalf("Failed to setup test database: %v", err)
	}
	suite.db = db
	suite.repo = NewTransactionRepository(db)
}

// TearDownSuite 在所有測試結束後執行一次
func (suite *TransactionRepositoryTestSuite) TearDownSuite() {
	if suite.db != nil {
		suite.db.Close()
	}
}

// SetupTest 在每個測試開始前執行
func (suite *TransactionRepositoryTestSuite) SetupTest() {
	// 清空測試資料
	_, err := suite.db.Exec("TRUNCATE TABLE transactions CASCADE")
	if err != nil {
		suite.T().Fatalf("Failed to truncate transactions table: %v", err)
	}
}

// TestCreate 測試建立交易記錄
func (suite *TransactionRepositoryTestSuite) TestCreate() {
	// Arrange
	fee := 28.0
	note := "定期定額買入"
	input := &models.CreateTransactionInput{
		Date:            time.Date(2025, 10, 22, 0, 0, 0, 0, time.UTC),
		AssetType:       models.AssetTypeTWStock,
		Symbol:          "2330",
		Name:            "台積電",
		TransactionType: models.TransactionTypeBuy,
		Quantity:        10,
		Price:           620,
		Amount:          6200,
		Fee:             &fee,
		Currency:        models.CurrencyTWD,
		Note:            &note,
	}

	// Act
	transaction, err := suite.repo.Create(input)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), transaction)
	assert.NotEqual(suite.T(), uuid.Nil, transaction.ID)
	assert.Equal(suite.T(), input.Date.Format("2006-01-02"), transaction.Date.Format("2006-01-02"))
	assert.Equal(suite.T(), input.AssetType, transaction.AssetType)
	assert.Equal(suite.T(), input.Symbol, transaction.Symbol)
	assert.Equal(suite.T(), input.Name, transaction.Name)
	assert.Equal(suite.T(), input.TransactionType, transaction.TransactionType)
	assert.Equal(suite.T(), input.Quantity, transaction.Quantity)
	assert.Equal(suite.T(), input.Price, transaction.Price)
	assert.Equal(suite.T(), input.Amount, transaction.Amount)
	assert.NotNil(suite.T(), transaction.Fee)
	assert.Equal(suite.T(), *input.Fee, *transaction.Fee)
	assert.NotNil(suite.T(), transaction.Note)
	assert.Equal(suite.T(), *input.Note, *transaction.Note)
	assert.False(suite.T(), transaction.CreatedAt.IsZero())
	assert.False(suite.T(), transaction.UpdatedAt.IsZero())
}

// TestGetByID 測試根據 ID 取得交易記錄
func (suite *TransactionRepositoryTestSuite) TestGetByID() {
	// Arrange - 先建立一筆交易記錄
	fee := 28.0
	input := &models.CreateTransactionInput{
		Date:            time.Date(2025, 10, 22, 0, 0, 0, 0, time.UTC),
		AssetType:       models.AssetTypeTWStock,
		Symbol:          "2330",
		Name:            "台積電",
		TransactionType: models.TransactionTypeBuy,
		Quantity:        10,
		Price:           620,
		Amount:          6200,
		Fee:             &fee,
		Currency:        models.CurrencyTWD,
	}
	created, err := suite.repo.Create(input)
	assert.NoError(suite.T(), err)

	// Act
	transaction, err := suite.repo.GetByID(created.ID)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), transaction)
	assert.Equal(suite.T(), created.ID, transaction.ID)
	assert.Equal(suite.T(), created.Symbol, transaction.Symbol)
}

// TestGetByID_NotFound 測試取得不存在的交易記錄
func (suite *TransactionRepositoryTestSuite) TestGetByID_NotFound() {
	// Arrange
	nonExistentID := uuid.New()

	// Act
	transaction, err := suite.repo.GetByID(nonExistentID)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), transaction)
	assert.Contains(suite.T(), err.Error(), "not found")
}

// TestGetAll 測試取得所有交易記錄
func (suite *TransactionRepositoryTestSuite) TestGetAll() {
	// Arrange - 建立多筆交易記錄
	transactions := []models.CreateTransactionInput{
		{
			Date:            time.Date(2025, 10, 22, 0, 0, 0, 0, time.UTC),
			AssetType:       models.AssetTypeTWStock,
			Symbol:          "2330",
			Name:            "台積電",
			TransactionType: models.TransactionTypeBuy,
			Quantity:        10,
			Price:           620,
			Amount:          6200,
			Currency:        models.CurrencyTWD,
		},
		{
			Date:            time.Date(2025, 10, 21, 0, 0, 0, 0, time.UTC),
			AssetType:       models.AssetTypeCrypto,
			Symbol:          "ETH",
			Name:            "Ethereum",
			TransactionType: models.TransactionTypeBuy,
			Quantity:        2,
			Price:           50000,
			Amount:          100000,
			Currency:        models.CurrencyUSD,
		},
	}

	for _, input := range transactions {
		_, err := suite.repo.Create(&input)
		assert.NoError(suite.T(), err)
	}

	// Act
	results, err := suite.repo.GetAll(TransactionFilters{})

	// Assert
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), results, 2)
	// 應該按日期降序排列
	assert.Equal(suite.T(), "2330", results[0].Symbol)
	assert.Equal(suite.T(), "ETH", results[1].Symbol)
}

// TestGetAll_WithFilters 測試使用篩選條件取得交易記錄
func (suite *TransactionRepositoryTestSuite) TestGetAll_WithFilters() {
	// Arrange
	transactions := []models.CreateTransactionInput{
		{
			Date:            time.Date(2025, 10, 22, 0, 0, 0, 0, time.UTC),
			AssetType:       models.AssetTypeTWStock,
			Symbol:          "2330",
			Name:            "台積電",
			TransactionType: models.TransactionTypeBuy,
			Quantity:        10,
			Price:           620,
			Amount:          6200,
			Currency:        models.CurrencyTWD,
		},
		{
			Date:            time.Date(2025, 10, 21, 0, 0, 0, 0, time.UTC),
			AssetType:       models.AssetTypeCrypto,
			Symbol:          "ETH",
			Name:            "Ethereum",
			TransactionType: models.TransactionTypeBuy,
			Quantity:        2,
			Price:           50000,
			Amount:          100000,
			Currency:        models.CurrencyUSD,
		},
	}

	for _, input := range transactions {
		_, err := suite.repo.Create(&input)
		assert.NoError(suite.T(), err)
	}

	// Act - 只取得台股的交易記錄
	assetType := models.AssetTypeTWStock
	results, err := suite.repo.GetAll(TransactionFilters{
		AssetType: &assetType,
	})

	// Assert
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), results, 1)
	assert.Equal(suite.T(), "2330", results[0].Symbol)
}

// TestUpdate 測試更新交易記錄
func (suite *TransactionRepositoryTestSuite) TestUpdate() {
	// Arrange - 先建立一筆交易記錄
	fee := 28.0
	input := &models.CreateTransactionInput{
		Date:            time.Date(2025, 10, 22, 0, 0, 0, 0, time.UTC),
		AssetType:       models.AssetTypeTWStock,
		Symbol:          "2330",
		Name:            "台積電",
		TransactionType: models.TransactionTypeBuy,
		Quantity:        10,
		Price:           620,
		Amount:          6200,
		Fee:             &fee,
		Currency:        models.CurrencyTWD,
	}
	created, err := suite.repo.Create(input)
	assert.NoError(suite.T(), err)

	// Act - 更新數量和價格
	newQuantity := 20.0
	newPrice := 630.0
	newAmount := 12600.0
	updateInput := &models.UpdateTransactionInput{
		Quantity: &newQuantity,
		Price:    &newPrice,
		Amount:   &newAmount,
	}
	updated, err := suite.repo.Update(created.ID, updateInput)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), updated)
	assert.Equal(suite.T(), created.ID, updated.ID)
	assert.Equal(suite.T(), newQuantity, updated.Quantity)
	assert.Equal(suite.T(), newPrice, updated.Price)
	assert.Equal(suite.T(), newAmount, updated.Amount)
	// 其他欄位應該保持不變
	assert.Equal(suite.T(), created.Symbol, updated.Symbol)
}

// TestDelete 測試刪除交易記錄
func (suite *TransactionRepositoryTestSuite) TestDelete() {
	// Arrange - 先建立一筆交易記錄
	fee := 28.0
	input := &models.CreateTransactionInput{
		Date:            time.Date(2025, 10, 22, 0, 0, 0, 0, time.UTC),
		AssetType:       models.AssetTypeTWStock,
		Symbol:          "2330",
		Name:            "台積電",
		TransactionType: models.TransactionTypeBuy,
		Quantity:        10,
		Price:           620,
		Amount:          6200,
		Fee:             &fee,
		Currency:        models.CurrencyTWD,
	}
	created, err := suite.repo.Create(input)
	assert.NoError(suite.T(), err)

	// Act
	err = suite.repo.Delete(created.ID)

	// Assert
	assert.NoError(suite.T(), err)

	// 驗證記錄已被刪除
	deleted, err := suite.repo.GetByID(created.ID)
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), deleted)
}

// TestTransactionRepositorySuite 執行測試套件
func TestTransactionRepositorySuite(t *testing.T) {
	suite.Run(t, new(TransactionRepositoryTestSuite))
}

