package repository

import (
	"database/sql"
	"testing"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// ExchangeRateRepositoryTestSuite 測試套件
type ExchangeRateRepositoryTestSuite struct {
	suite.Suite
	db   *sql.DB
	repo ExchangeRateRepository
}

// SetupSuite 在所有測試開始前執行一次
func (suite *ExchangeRateRepositoryTestSuite) SetupSuite() {
	db, err := setupTestDB()
	if err != nil {
		suite.T().Fatalf("Failed to setup test database: %v", err)
	}
	suite.db = db
	suite.repo = NewExchangeRateRepository(db)
}

// TearDownSuite 在所有測試結束後執行一次
func (suite *ExchangeRateRepositoryTestSuite) TearDownSuite() {
	if suite.db != nil {
		suite.db.Close()
	}
}

// SetupTest 在每個測試開始前執行
func (suite *ExchangeRateRepositoryTestSuite) SetupTest() {
	// 清空測試資料
	_, err := suite.db.Exec("TRUNCATE TABLE exchange_rates CASCADE")
	if err != nil {
		suite.T().Fatalf("Failed to truncate exchange_rates table: %v", err)
	}
}

// TestUpsert_Create 測試 Upsert 建立新記錄
func (suite *ExchangeRateRepositoryTestSuite) TestUpsert_Create() {
	// Arrange
	input := &models.ExchangeRateInput{
		FromCurrency: models.CurrencyUSD,
		ToCurrency:   models.CurrencyTWD,
		Rate:         31.5,
		Date:         time.Date(2025, 10, 28, 0, 0, 0, 0, time.UTC),
	}

	// Act
	result, err := suite.repo.Upsert(input)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), input.FromCurrency, result.FromCurrency)
	assert.Equal(suite.T(), input.ToCurrency, result.ToCurrency)
	assert.Equal(suite.T(), input.Rate, result.Rate)
	assert.Equal(suite.T(), input.Date.Format("2006-01-02"), result.Date.Format("2006-01-02"))
	assert.NotZero(suite.T(), result.CreatedAt)
	assert.NotZero(suite.T(), result.UpdatedAt)
	// 新建立時，created_at 和 updated_at 應該相同
	assert.Equal(suite.T(), result.CreatedAt.Unix(), result.UpdatedAt.Unix())
}

// TestUpsert_Update 測試 Upsert 更新現有記錄
func (suite *ExchangeRateRepositoryTestSuite) TestUpsert_Update() {
	// Arrange - 先建立一筆記錄
	date := time.Date(2025, 10, 28, 0, 0, 0, 0, time.UTC)
	input1 := &models.ExchangeRateInput{
		FromCurrency: models.CurrencyUSD,
		ToCurrency:   models.CurrencyTWD,
		Rate:         31.5,
		Date:         date,
	}
	first, err := suite.repo.Upsert(input1)
	assert.NoError(suite.T(), err)

	// 等待一秒確保時間戳不同
	time.Sleep(1 * time.Second)

	// Act - 更新同一天的匯率
	input2 := &models.ExchangeRateInput{
		FromCurrency: models.CurrencyUSD,
		ToCurrency:   models.CurrencyTWD,
		Rate:         32.0, // 更新匯率
		Date:         date,
	}
	second, err := suite.repo.Upsert(input2)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), second)
	assert.Equal(suite.T(), first.ID, second.ID) // ID 應該相同
	assert.Equal(suite.T(), input2.Rate, second.Rate) // 匯率應該更新
	assert.NotZero(suite.T(), second.UpdatedAt)
	// updated_at 應該比 created_at 晚
	assert.True(suite.T(), second.UpdatedAt.After(second.CreatedAt),
		"updated_at (%v) should be after created_at (%v)",
		second.UpdatedAt, second.CreatedAt)
}

// TestUpsert_MultipleDates 測試不同日期的 Upsert
func (suite *ExchangeRateRepositoryTestSuite) TestUpsert_MultipleDates() {
	// Arrange
	date1 := time.Date(2025, 10, 27, 0, 0, 0, 0, time.UTC)
	date2 := time.Date(2025, 10, 28, 0, 0, 0, 0, time.UTC)

	input1 := &models.ExchangeRateInput{
		FromCurrency: models.CurrencyUSD,
		ToCurrency:   models.CurrencyTWD,
		Rate:         31.5,
		Date:         date1,
	}
	input2 := &models.ExchangeRateInput{
		FromCurrency: models.CurrencyUSD,
		ToCurrency:   models.CurrencyTWD,
		Rate:         32.0,
		Date:         date2,
	}

	// Act
	result1, err1 := suite.repo.Upsert(input1)
	result2, err2 := suite.repo.Upsert(input2)

	// Assert
	assert.NoError(suite.T(), err1)
	assert.NoError(suite.T(), err2)
	assert.NotEqual(suite.T(), result1.ID, result2.ID) // 不同日期應該有不同 ID
	assert.Equal(suite.T(), input1.Rate, result1.Rate)
	assert.Equal(suite.T(), input2.Rate, result2.Rate)
}

// TestCreate 測試建立匯率記錄
func (suite *ExchangeRateRepositoryTestSuite) TestCreate() {
	// Arrange
	input := &models.ExchangeRateInput{
		FromCurrency: models.CurrencyUSD,
		ToCurrency:   models.CurrencyTWD,
		Rate:         31.5,
		Date:         time.Date(2025, 10, 28, 0, 0, 0, 0, time.UTC),
	}

	// Act
	result, err := suite.repo.Create(input)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.NotZero(suite.T(), result.ID)
	assert.Equal(suite.T(), input.FromCurrency, result.FromCurrency)
	assert.Equal(suite.T(), input.ToCurrency, result.ToCurrency)
	assert.Equal(suite.T(), input.Rate, result.Rate)
	assert.NotZero(suite.T(), result.CreatedAt)
	assert.NotZero(suite.T(), result.UpdatedAt)
}

// TestGetByDate 測試根據日期取得匯率
func (suite *ExchangeRateRepositoryTestSuite) TestGetByDate() {
	// Arrange
	date := time.Date(2025, 10, 28, 0, 0, 0, 0, time.UTC)
	input := &models.ExchangeRateInput{
		FromCurrency: models.CurrencyUSD,
		ToCurrency:   models.CurrencyTWD,
		Rate:         31.5,
		Date:         date,
	}
	_, err := suite.repo.Create(input)
	assert.NoError(suite.T(), err)

	// Act
	result, err := suite.repo.GetByDate(models.CurrencyUSD, models.CurrencyTWD, date)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), input.Rate, result.Rate)
}

// TestGetLatest 測試取得最新匯率
func (suite *ExchangeRateRepositoryTestSuite) TestGetLatest() {
	// Arrange - 建立多筆不同日期的記錄
	dates := []time.Time{
		time.Date(2025, 10, 26, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 10, 27, 0, 0, 0, 0, time.UTC),
		time.Date(2025, 10, 28, 0, 0, 0, 0, time.UTC),
	}
	rates := []float64{31.0, 31.5, 32.0}

	for i, date := range dates {
		input := &models.ExchangeRateInput{
			FromCurrency: models.CurrencyUSD,
			ToCurrency:   models.CurrencyTWD,
			Rate:         rates[i],
			Date:         date,
		}
		_, err := suite.repo.Create(input)
		assert.NoError(suite.T(), err)
	}

	// Act
	result, err := suite.repo.GetLatest(models.CurrencyUSD, models.CurrencyTWD)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), 32.0, result.Rate) // 應該是最新的匯率
	assert.Equal(suite.T(), dates[2].Format("2006-01-02"), result.Date.Format("2006-01-02"))
}

// TestExchangeRateRepositoryTestSuite 執行測試套件
func TestExchangeRateRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(ExchangeRateRepositoryTestSuite))
}

