package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
)

// ExchangeRateRepository 匯率資料庫操作介面
type ExchangeRateRepository interface {
	// Create 建立新的匯率記錄
	Create(input *models.ExchangeRateInput) (*models.ExchangeRate, error)
	// GetByDate 取得指定日期的匯率
	GetByDate(fromCurrency, toCurrency models.Currency, date time.Time) (*models.ExchangeRate, error)
	// GetLatest 取得最新的匯率
	GetLatest(fromCurrency, toCurrency models.Currency) (*models.ExchangeRate, error)
	// Upsert 建立或更新匯率記錄
	Upsert(input *models.ExchangeRateInput) (*models.ExchangeRate, error)
}

// exchangeRateRepository 匯率資料庫操作實作
type exchangeRateRepository struct {
	db *sql.DB
}

// NewExchangeRateRepository 建立新的匯率 repository
func NewExchangeRateRepository(db *sql.DB) ExchangeRateRepository {
	return &exchangeRateRepository{db: db}
}

// Create 建立新的匯率記錄
func (r *exchangeRateRepository) Create(input *models.ExchangeRateInput) (*models.ExchangeRate, error) {
	query := `
		INSERT INTO exchange_rates (from_currency, to_currency, rate, date)
		VALUES ($1, $2, $3, $4)
		RETURNING id, from_currency, to_currency, rate, date, created_at, updated_at
	`

	var rate models.ExchangeRate
	err := r.db.QueryRow(
		query,
		input.FromCurrency,
		input.ToCurrency,
		input.Rate,
		input.Date,
	).Scan(&rate.ID, &rate.FromCurrency, &rate.ToCurrency, &rate.Rate, &rate.Date, &rate.CreatedAt, &rate.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create exchange rate: %w", err)
	}

	return &rate, nil
}

// GetByDate 取得指定日期的匯率
func (r *exchangeRateRepository) GetByDate(fromCurrency, toCurrency models.Currency, date time.Time) (*models.ExchangeRate, error) {
	query := `
		SELECT id, from_currency, to_currency, rate, date, created_at, updated_at
		FROM exchange_rates
		WHERE from_currency = $1 AND to_currency = $2 AND date = $3
		LIMIT 1
	`

	var rate models.ExchangeRate
	err := r.db.QueryRow(query, fromCurrency, toCurrency, date).Scan(
		&rate.ID, &rate.FromCurrency, &rate.ToCurrency, &rate.Rate, &rate.Date, &rate.CreatedAt, &rate.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get exchange rate: %w", err)
	}

	return &rate, nil
}

// GetLatest 取得最新的匯率
func (r *exchangeRateRepository) GetLatest(fromCurrency, toCurrency models.Currency) (*models.ExchangeRate, error) {
	query := `
		SELECT id, from_currency, to_currency, rate, date, created_at, updated_at
		FROM exchange_rates
		WHERE from_currency = $1 AND to_currency = $2
		ORDER BY date DESC
		LIMIT 1
	`

	var rate models.ExchangeRate
	err := r.db.QueryRow(query, fromCurrency, toCurrency).Scan(
		&rate.ID, &rate.FromCurrency, &rate.ToCurrency, &rate.Rate, &rate.Date, &rate.CreatedAt, &rate.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get latest exchange rate: %w", err)
	}

	return &rate, nil
}

// Upsert 建立或更新匯率記錄
func (r *exchangeRateRepository) Upsert(input *models.ExchangeRateInput) (*models.ExchangeRate, error) {
	query := `
		INSERT INTO exchange_rates (from_currency, to_currency, rate, date)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (from_currency, to_currency, date)
		DO UPDATE SET rate = EXCLUDED.rate, updated_at = CURRENT_TIMESTAMP
		RETURNING id, from_currency, to_currency, rate, date, created_at, updated_at
	`

	var rate models.ExchangeRate
	err := r.db.QueryRow(
		query,
		input.FromCurrency,
		input.ToCurrency,
		input.Rate,
		input.Date,
	).Scan(&rate.ID, &rate.FromCurrency, &rate.ToCurrency, &rate.Rate, &rate.Date, &rate.CreatedAt, &rate.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to upsert exchange rate: %w", err)
	}

	return &rate, nil
}

