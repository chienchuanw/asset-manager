package discord

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type MockParser struct{}

func (p MockParser) Parse(ctx context.Context, message string, categories []CategoryInfo) (*ParseResult, error) {
	switch message {
	case "午餐吃拉麵 180", "lunch ramen 180":
		return &ParseResult{
			IsBookkeeping: true,
			Type:          "expense",
			Amount:        180,
			Description:   "午餐吃拉麵",
			CategoryID:    "expense-food",
			CategoryName:  "飲食",
			Date:          "2026-04-05",
		}, nil
	case "收到薪水 45000":
		return &ParseResult{
			IsBookkeeping: true,
			Type:          "income",
			Amount:        45000,
			Description:   "收到薪水",
			CategoryID:    "income-salary",
			CategoryName:  "薪資",
			Date:          "2026-04-05",
		}, nil
	case "你好":
		return &ParseResult{IsBookkeeping: false}, nil
	case "午餐吃拉麵":
		return &ParseResult{
			IsBookkeeping: true,
			Type:          "expense",
			Description:   "午餐吃拉麵",
			CategoryID:    "expense-food",
			CategoryName:  "飲食",
			Date:          "2026-04-05",
			MissingFields: []string{"amount"},
		}, nil
	default:
		return nil, errors.New("unexpected message")
	}
}

func testCategories() []CategoryInfo {
	return []CategoryInfo{
		{ID: "income-salary", Name: "薪資", Type: "income"},
		{ID: "expense-food", Name: "飲食", Type: "expense"},
		{ID: "income-other", Name: "其他收入", Type: "income"},
		{ID: "expense-other", Name: "其他支出", Type: "expense"},
	}
}

func TestMockParser_ChineseExpense(t *testing.T) {
	parser := MockParser{}

	result, err := parser.Parse(t.Context(), "午餐吃拉麵 180", testCategories())

	require.NoError(t, err)
	require.True(t, result.IsBookkeeping)
	require.Equal(t, "expense", result.Type)
	require.Equal(t, 180.0, result.Amount)
	require.Equal(t, "飲食", result.CategoryName)
}

func TestMockParser_EnglishExpense(t *testing.T) {
	parser := MockParser{}

	result, err := parser.Parse(t.Context(), "lunch ramen 180", testCategories())

	require.NoError(t, err)
	require.True(t, result.IsBookkeeping)
	require.Equal(t, "expense", result.Type)
	require.Equal(t, 180.0, result.Amount)
	require.Equal(t, "飲食", result.CategoryName)
}

func TestMockParser_Income(t *testing.T) {
	parser := MockParser{}

	result, err := parser.Parse(t.Context(), "收到薪水 45000", testCategories())

	require.NoError(t, err)
	require.True(t, result.IsBookkeeping)
	require.Equal(t, "income", result.Type)
	require.Equal(t, 45000.0, result.Amount)
	require.Equal(t, "薪資", result.CategoryName)
}

func TestMockParser_NonBookkeeping(t *testing.T) {
	parser := MockParser{}

	result, err := parser.Parse(t.Context(), "你好", testCategories())

	require.NoError(t, err)
	require.False(t, result.IsBookkeeping)
}

func TestMockParser_MissingAmount(t *testing.T) {
	parser := MockParser{}

	result, err := parser.Parse(t.Context(), "午餐吃拉麵", testCategories())

	require.NoError(t, err)
	require.True(t, result.IsBookkeeping)
	require.Equal(t, []string{"amount"}, result.MissingFields)
}

func TestGeminiParser_Errors(t *testing.T) {
	t.Run("api key not set", func(t *testing.T) {
		parser := NewGeminiParser("")

		result, err := parser.Parse(t.Context(), "午餐吃拉麵 180", testCategories())

		require.Nil(t, result)
		require.Error(t, err)
		require.ErrorContains(t, err, "Gemini API key is not set")
	})

	t.Run("invalid json response", func(t *testing.T) {
		originalGenerate := geminiGenerateContent
		originalNow := geminiNow
		t.Cleanup(func() {
			geminiGenerateContent = originalGenerate
			geminiNow = originalNow
		})

		geminiNow = func() time.Time { return time.Date(2026, 4, 5, 0, 0, 0, 0, time.UTC) }
		geminiGenerateContent = func(ctx context.Context, apiKey, model, prompt string) (string, error) {
			return "not-json", nil
		}

		parser := NewGeminiParser("test-key")

		result, err := parser.Parse(t.Context(), "午餐吃拉麵 180", testCategories())

		require.Nil(t, result)
		require.Error(t, err)
		require.ErrorContains(t, err, "invalid parser response")
		require.ErrorContains(t, err, "not-json")
	})

	t.Run("context timeout", func(t *testing.T) {
		originalGenerate := geminiGenerateContent
		t.Cleanup(func() {
			geminiGenerateContent = originalGenerate
		})

		geminiGenerateContent = func(ctx context.Context, apiKey, model, prompt string) (string, error) {
			<-ctx.Done()
			return "", ctx.Err()
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Nanosecond)
		defer cancel()
		time.Sleep(time.Millisecond)

		parser := NewGeminiParser("test-key")

		result, err := parser.Parse(ctx, "午餐吃拉麵 180", testCategories())

		require.Nil(t, result)
		require.Error(t, err)
		require.ErrorContains(t, err, context.DeadlineExceeded.Error())
	})
}
