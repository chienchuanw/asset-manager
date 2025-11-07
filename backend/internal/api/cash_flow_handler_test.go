package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/chienchuanw/asset-manager/internal/models"
	"github.com/chienchuanw/asset-manager/internal/repository"
	"github.com/chienchuanw/asset-manager/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCashFlowService 模擬的 CashFlowService
type MockCashFlowService struct {
	mock.Mock
}

func (m *MockCashFlowService) CreateCashFlow(input *models.CreateCashFlowInput) (*models.CashFlow, error) {
	args := m.Called(input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CashFlow), args.Error(1)
}

func (m *MockCashFlowService) GetCashFlow(id uuid.UUID) (*models.CashFlow, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CashFlow), args.Error(1)
}

func (m *MockCashFlowService) ListCashFlows(filters repository.CashFlowFilters) ([]*models.CashFlow, error) {
	args := m.Called(filters)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.CashFlow), args.Error(1)
}

func (m *MockCashFlowService) UpdateCashFlow(id uuid.UUID, input *models.UpdateCashFlowInput) (*models.CashFlow, error) {
	args := m.Called(id, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.CashFlow), args.Error(1)
}

func (m *MockCashFlowService) DeleteCashFlow(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockCashFlowService) GetSummary(startDate, endDate time.Time) (*repository.CashFlowSummary, error) {
	args := m.Called(startDate, endDate)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*repository.CashFlowSummary), args.Error(1)
}

func (m *MockCashFlowService) GetMonthlySummaryWithComparison(year, month int) (*models.MonthlyCashFlowSummary, error) {
	args := m.Called(year, month)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.MonthlyCashFlowSummary), args.Error(1)
}

func (m *MockCashFlowService) GetYearlySummaryWithComparison(year int) (*models.YearlyCashFlowSummary, error) {
	args := m.Called(year)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.YearlyCashFlowSummary), args.Error(1)
}

// MockDiscordService 模擬的 DiscordService
type MockDiscordService struct {
	mock.Mock
}

func (m *MockDiscordService) SendMessage(webhookURL string, message *models.DiscordMessage) error {
	args := m.Called(webhookURL, message)
	return args.Error(0)
}

func (m *MockDiscordService) FormatDailyReport(data *models.DailyReportData) *models.DiscordMessage {
	args := m.Called(data)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*models.DiscordMessage)
}

func (m *MockDiscordService) SendDailyBillingNotification(webhookURL string, result *service.DailyBillingResult) error {
	args := m.Called(webhookURL, result)
	return args.Error(0)
}

func (m *MockDiscordService) SendSubscriptionExpiryNotification(webhookURL string, subscriptions []*models.Subscription, days int) error {
	args := m.Called(webhookURL, subscriptions, days)
	return args.Error(0)
}

func (m *MockDiscordService) SendInstallmentCompletionNotification(webhookURL string, installments []*models.Installment, remainingCount int) error {
	args := m.Called(webhookURL, installments, remainingCount)
	return args.Error(0)
}

func (m *MockDiscordService) SendCreditCardPaymentReminder(webhookURL string, creditCards []*models.CreditCard) error {
	args := m.Called(webhookURL, creditCards)
	return args.Error(0)
}

func (m *MockDiscordService) FormatMonthlyCashFlowReport(summary *models.MonthlyCashFlowSummary) *models.DiscordMessage {
	args := m.Called(summary)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*models.DiscordMessage)
}

func (m *MockDiscordService) FormatYearlyCashFlowReport(summary *models.YearlyCashFlowSummary) *models.DiscordMessage {
	args := m.Called(summary)
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*models.DiscordMessage)
}

// setupCashFlowTestRouter 設定測試用的 router
func setupCashFlowTestRouter(handler *CashFlowHandler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	api := router.Group("/api")
	{
		cashFlows := api.Group("/cash-flows")
		{
			cashFlows.POST("", handler.CreateCashFlow)
			cashFlows.GET("", handler.ListCashFlows)
			cashFlows.GET("/summary", handler.GetSummary)
			cashFlows.GET("/monthly-summary", handler.GetMonthlySummary)
			cashFlows.GET("/yearly-summary", handler.GetYearlySummary)
			cashFlows.POST("/send-monthly-report", handler.SendMonthlyReport)
			cashFlows.POST("/send-yearly-report", handler.SendYearlyReport)
			cashFlows.GET("/:id", handler.GetCashFlow)
			cashFlows.PUT("/:id", handler.UpdateCashFlow)
			cashFlows.DELETE("/:id", handler.DeleteCashFlow)
		}
	}

	return router
}

// TestCreateCashFlow_Success 測試成功建立現金流記錄
func TestCreateCashFlow_Success(t *testing.T) {
	// Arrange
	mockService := new(MockCashFlowService)
	handler := NewCashFlowHandler(mockService)
	router := setupCashFlowTestRouter(handler)

	categoryID := uuid.New()
	input := models.CreateCashFlowInput{
		Date:        time.Date(2025, 10, 25, 0, 0, 0, 0, time.UTC),
		Type:        models.CashFlowTypeIncome,
		CategoryID:  categoryID,
		Amount:      50000,
		Description: "十月薪資",
	}

	expectedCashFlow := &models.CashFlow{
		ID:          uuid.New(),
		Date:        input.Date,
		Type:        input.Type,
		CategoryID:  input.CategoryID,
		Amount:      input.Amount,
		Currency:    models.CurrencyTWD,
		Description: input.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mockService.On("CreateCashFlow", &input).Return(expectedCashFlow, nil)

	// 準備請求
	body, _ := json.Marshal(input)
	req, _ := http.NewRequest("POST", "/api/cash-flows", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Nil(t, response.Error)
	assert.NotNil(t, response.Data)

	mockService.AssertExpectations(t)
}

// TestCreateCashFlow_InvalidInput 測試無效的輸入資料
func TestCreateCashFlow_InvalidInput(t *testing.T) {
	// Arrange
	mockService := new(MockCashFlowService)
	handler := NewCashFlowHandler(mockService)
	router := setupCashFlowTestRouter(handler)

	// 無效的 JSON
	invalidJSON := []byte(`{"invalid": json}`)

	req, _ := http.NewRequest("POST", "/api/cash-flows", bytes.NewBuffer(invalidJSON))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response.Error)
	assert.Equal(t, "INVALID_INPUT", response.Error.Code)
}

// TestGetCashFlow_Success 測試成功取得現金流記錄
func TestGetCashFlow_Success(t *testing.T) {
	// Arrange
	mockService := new(MockCashFlowService)
	handler := NewCashFlowHandler(mockService)
	router := setupCashFlowTestRouter(handler)

	cashFlowID := uuid.New()
	expectedCashFlow := &models.CashFlow{
		ID:          cashFlowID,
		Date:        time.Date(2025, 10, 25, 0, 0, 0, 0, time.UTC),
		Type:        models.CashFlowTypeIncome,
		Amount:      50000,
		Description: "薪資",
	}

	mockService.On("GetCashFlow", cashFlowID).Return(expectedCashFlow, nil)

	// 準備請求
	req, _ := http.NewRequest("GET", fmt.Sprintf("/api/cash-flows/%s", cashFlowID), nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Nil(t, response.Error)
	assert.NotNil(t, response.Data)

	mockService.AssertExpectations(t)
}

// TestGetCashFlow_InvalidID 測試無效的 ID
func TestGetCashFlow_InvalidID(t *testing.T) {
	// Arrange
	mockService := new(MockCashFlowService)
	handler := NewCashFlowHandler(mockService)
	router := setupCashFlowTestRouter(handler)

	// 準備請求
	req, _ := http.NewRequest("GET", "/api/cash-flows/invalid-id", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response.Error)
	assert.Equal(t, "INVALID_ID", response.Error.Code)
}

// TestListCashFlows_Success 測試成功取得現金流列表
func TestListCashFlows_Success(t *testing.T) {
	// Arrange
	mockService := new(MockCashFlowService)
	handler := NewCashFlowHandler(mockService)
	router := setupCashFlowTestRouter(handler)

	expectedCashFlows := []*models.CashFlow{
		{
			ID:          uuid.New(),
			Type:        models.CashFlowTypeIncome,
			Amount:      50000,
			Description: "薪資",
		},
		{
			ID:          uuid.New(),
			Type:        models.CashFlowTypeExpense,
			Amount:      1200,
			Description: "午餐",
		},
	}

	mockService.On("ListCashFlows", mock.AnythingOfType("repository.CashFlowFilters")).Return(expectedCashFlows, nil)

	// 準備請求
	req, _ := http.NewRequest("GET", "/api/cash-flows", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Nil(t, response.Error)
	assert.NotNil(t, response.Data)

	mockService.AssertExpectations(t)
}

// TestDeleteCashFlow_Success 測試成功刪除現金流記錄
func TestDeleteCashFlow_Success(t *testing.T) {
	// Arrange
	mockService := new(MockCashFlowService)
	handler := NewCashFlowHandler(mockService)
	router := setupCashFlowTestRouter(handler)

	cashFlowID := uuid.New()
	mockService.On("DeleteCashFlow", cashFlowID).Return(nil)

	// 準備請求
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/cash-flows/%s", cashFlowID), nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusNoContent, w.Code)
	mockService.AssertExpectations(t)
}

// TestGetSummary_Success 測試成功取得摘要
func TestGetSummary_Success(t *testing.T) {
	// Arrange
	mockService := new(MockCashFlowService)
	handler := NewCashFlowHandler(mockService)
	router := setupCashFlowTestRouter(handler)

	startDate := time.Date(2025, 10, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 10, 31, 0, 0, 0, 0, time.UTC)

	expectedSummary := &repository.CashFlowSummary{
		TotalIncome:  55000,
		TotalExpense: 15000,
		NetCashFlow:  40000,
	}

	mockService.On("GetSummary", startDate, endDate).Return(expectedSummary, nil)

	// 準備請求
	req, _ := http.NewRequest("GET", "/api/cash-flows/summary?start_date=2025-10-01&end_date=2025-10-31", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Nil(t, response.Error)
	assert.NotNil(t, response.Data)

	mockService.AssertExpectations(t)
}

// TestGetSummary_MissingParameters 測試缺少參數
func TestGetSummary_MissingParameters(t *testing.T) {
	// Arrange
	mockService := new(MockCashFlowService)
	handler := NewCashFlowHandler(mockService)
	router := setupCashFlowTestRouter(handler)

	// 準備請求（缺少 end_date）
	req, _ := http.NewRequest("GET", "/api/cash-flows/summary?start_date=2025-10-01", nil)
	w := httptest.NewRecorder()

	// Act
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response.Error)
	assert.Equal(t, "MISSING_PARAMETERS", response.Error.Code)
}

// TestGetMonthlySummary_Success 測試成功取得月度摘要
func TestGetMonthlySummary_Success(t *testing.T) {
	// Arrange
	mockService := new(MockCashFlowService)
	handler := NewCashFlowHandler(mockService)

	expectedSummary := &models.MonthlyCashFlowSummary{
		Year:         2024,
		Month:        1,
		TotalIncome:  50000,
		TotalExpense: 30000,
		NetCashFlow:  20000,
		IncomeCount:  5,
		ExpenseCount: 10,
	}

	mockService.On("GetMonthlySummaryWithComparison", 2024, 1).Return(expectedSummary, nil)

	router := setupCashFlowTestRouter(handler)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/cash-flows/monthly-summary?year=2024&month=1", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Nil(t, response.Error)
	assert.NotNil(t, response.Data)

	mockService.AssertExpectations(t)
}

// TestGetMonthlySummary_MissingYear 測試缺少年份參數
func TestGetMonthlySummary_MissingYear(t *testing.T) {
	// Arrange
	mockService := new(MockCashFlowService)
	handler := NewCashFlowHandler(mockService)
	router := setupCashFlowTestRouter(handler)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/cash-flows/monthly-summary?month=1", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response.Error)
	assert.Equal(t, "MISSING_YEAR", response.Error.Code)
}

// TestGetMonthlySummary_InvalidMonth 測試無效的月份參數
func TestGetMonthlySummary_InvalidMonth(t *testing.T) {
	// Arrange
	mockService := new(MockCashFlowService)
	handler := NewCashFlowHandler(mockService)
	router := setupCashFlowTestRouter(handler)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/cash-flows/monthly-summary?year=2024&month=13", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.NotNil(t, response.Error)
	assert.Equal(t, "INVALID_MONTH", response.Error.Code)
}

// TestGetYearlySummary_Success 測試成功取得年度摘要
func TestGetYearlySummary_Success(t *testing.T) {
	// Arrange
	mockService := new(MockCashFlowService)
	handler := NewCashFlowHandler(mockService)

	expectedSummary := &models.YearlyCashFlowSummary{
		Year:         2024,
		TotalIncome:  600000,
		TotalExpense: 360000,
		NetCashFlow:  240000,
		IncomeCount:  60,
		ExpenseCount: 120,
	}

	mockService.On("GetYearlySummaryWithComparison", 2024).Return(expectedSummary, nil)

	router := setupCashFlowTestRouter(handler)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/cash-flows/yearly-summary?year=2024", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Nil(t, response.Error)
	assert.NotNil(t, response.Data)

	mockService.AssertExpectations(t)
}

// TestSendMonthlyReport_Success 測試成功發送月度報告
func TestSendMonthlyReport_Success(t *testing.T) {
	// Arrange
	mockCashFlowService := new(MockCashFlowService)
	mockDiscordService := new(MockDiscordService)
	handler := NewCashFlowHandler(mockCashFlowService)
	handler.SetDiscordService(mockDiscordService)

	summary := &models.MonthlyCashFlowSummary{
		Year:         2024,
		Month:        1,
		TotalIncome:  50000,
		TotalExpense: 30000,
		NetCashFlow:  20000,
	}

	message := &models.DiscordMessage{
		Content: "Test monthly report",
	}

	mockCashFlowService.On("GetMonthlySummaryWithComparison", 2024, 1).Return(summary, nil)
	mockDiscordService.On("FormatMonthlyCashFlowReport", summary).Return(message)
	mockDiscordService.On("SendMessage", "https://discord.webhook.url", message).Return(nil)

	router := setupCashFlowTestRouter(handler)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/cash-flows/send-monthly-report?year=2024&month=1&webhook_url=https://discord.webhook.url", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Nil(t, response.Error)
	assert.Equal(t, "Monthly report sent successfully", response.Data)

	mockCashFlowService.AssertExpectations(t)
	mockDiscordService.AssertExpectations(t)
}

// TestSendYearlyReport_Success 測試成功發送年度報告
func TestSendYearlyReport_Success(t *testing.T) {
	// Arrange
	mockCashFlowService := new(MockCashFlowService)
	mockDiscordService := new(MockDiscordService)
	handler := NewCashFlowHandler(mockCashFlowService)
	handler.SetDiscordService(mockDiscordService)

	summary := &models.YearlyCashFlowSummary{
		Year:         2024,
		TotalIncome:  600000,
		TotalExpense: 360000,
		NetCashFlow:  240000,
	}

	message := &models.DiscordMessage{
		Content: "Test yearly report",
	}

	mockCashFlowService.On("GetYearlySummaryWithComparison", 2024).Return(summary, nil)
	mockDiscordService.On("FormatYearlyCashFlowReport", summary).Return(message)
	mockDiscordService.On("SendMessage", "https://discord.webhook.url", message).Return(nil)

	router := setupCashFlowTestRouter(handler)

	// Act
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/cash-flows/send-yearly-report?year=2024&webhook_url=https://discord.webhook.url", nil)
	router.ServeHTTP(w, req)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)

	var response APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Nil(t, response.Error)
	assert.Equal(t, "Yearly report sent successfully", response.Data)

	mockCashFlowService.AssertExpectations(t)
	mockDiscordService.AssertExpectations(t)
}
