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
	case "刷卡買衣服 2000":
		return &ParseResult{
			IsBookkeeping: true,
			Type:          "expense",
			Amount:        2000,
			Description:   "刷卡買衣服",
			CategoryID:    "expense-other",
			CategoryName:  "其他支出",
			Date:          "2026-04-05",
			SourceType:    "credit_card",
		}, nil
	case "轉帳繳房租 15000":
		return &ParseResult{
			IsBookkeeping: true,
			Type:          "expense",
			Amount:        15000,
			Description:   "轉帳繳房租",
			CategoryID:    "expense-other",
			CategoryName:  "其他支出",
			Date:          "2026-04-05",
			SourceType:    "bank_account",
		}, nil
	case "昨天午餐 180":
		return &ParseResult{
			IsBookkeeping: true,
			Type:          "expense",
			Amount:        180,
			Description:   "午餐",
			CategoryID:    "expense-food",
			CategoryName:  "飲食",
			Date:          "2026-04-04",
		}, nil
	case "前天晚餐 300":
		return &ParseResult{
			IsBookkeeping: true,
			Type:          "expense",
			Amount:        300,
			Description:   "晚餐",
			CategoryID:    "expense-food",
			CategoryName:  "飲食",
			Date:          "2026-04-03",
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

func TestMockParser_CreditCardSourceType(t *testing.T) {
	parser := MockParser{}

	result, err := parser.Parse(t.Context(), "刷卡買衣服 2000", testCategories())

	require.NoError(t, err)
	require.True(t, result.IsBookkeeping)
	require.Equal(t, "expense", result.Type)
	require.Equal(t, 2000.0, result.Amount)
	require.Equal(t, "credit_card", result.SourceType)
}

func TestMockParser_BankAccountSourceType(t *testing.T) {
	parser := MockParser{}

	result, err := parser.Parse(t.Context(), "轉帳繳房租 15000", testCategories())

	require.NoError(t, err)
	require.True(t, result.IsBookkeeping)
	require.Equal(t, "expense", result.Type)
	require.Equal(t, 15000.0, result.Amount)
	require.Equal(t, "bank_account", result.SourceType)
}

func TestMockParser_EmptySourceType(t *testing.T) {
	parser := MockParser{}

	result, err := parser.Parse(t.Context(), "午餐吃拉麵 180", testCategories())

	require.NoError(t, err)
	require.True(t, result.IsBookkeeping)
	require.Equal(t, "", result.SourceType)
}

func TestMockParser_YesterdayDate(t *testing.T) {
	parser := MockParser{}

	result, err := parser.Parse(t.Context(), "昨天午餐 180", testCategories())

	require.NoError(t, err)
	require.True(t, result.IsBookkeeping)
	require.Equal(t, "2026-04-04", result.Date)
}

func TestMockParser_DayBeforeYesterdayDate(t *testing.T) {
	parser := MockParser{}

	result, err := parser.Parse(t.Context(), "前天晚餐 300", testCategories())

	require.NoError(t, err)
	require.True(t, result.IsBookkeeping)
	require.Equal(t, "2026-04-03", result.Date)
}

func TestGeminiParser_SourceTypeInResponse(t *testing.T) {
	originalGenerate := geminiGenerateContent
	originalNow := geminiNow
	t.Cleanup(func() {
		geminiGenerateContent = originalGenerate
		geminiNow = originalNow
	})

	geminiNow = func() time.Time { return time.Date(2026, 4, 5, 0, 0, 0, 0, time.UTC) }
	geminiGenerateContent = func(ctx context.Context, apiKey, model, prompt string) (string, error) {
		return `{"is_bookkeeping":true,"type":"expense","amount":2000,"description":"刷卡買衣服","category_id":"expense-other","category_name":"其他支出","date":"2026-04-05","source_type":"credit_card","missing_fields":[]}`, nil
	}

	parser := NewGeminiParser("test-key")

	result, err := parser.Parse(t.Context(), "刷卡買衣服 2000", testCategories())

	require.NoError(t, err)
	require.True(t, result.IsBookkeeping)
	require.Equal(t, "credit_card", result.SourceType)
}

func TestGeminiParser_EmptySourceType(t *testing.T) {
	originalGenerate := geminiGenerateContent
	originalNow := geminiNow
	t.Cleanup(func() {
		geminiGenerateContent = originalGenerate
		geminiNow = originalNow
	})

	geminiNow = func() time.Time { return time.Date(2026, 4, 5, 0, 0, 0, 0, time.UTC) }
	geminiGenerateContent = func(ctx context.Context, apiKey, model, prompt string) (string, error) {
		return `{"is_bookkeeping":true,"type":"expense","amount":180,"description":"午餐","category_id":"expense-food","category_name":"飲食","date":"2026-04-05","source_type":"","missing_fields":[]}`, nil
	}

	parser := NewGeminiParser("test-key")

	result, err := parser.Parse(t.Context(), "午餐 180", testCategories())

	require.NoError(t, err)
	require.True(t, result.IsBookkeeping)
	require.Equal(t, "", result.SourceType)
}

func TestGeminiParser_RelativeDate(t *testing.T) {
	originalGenerate := geminiGenerateContent
	originalNow := geminiNow
	t.Cleanup(func() {
		geminiGenerateContent = originalGenerate
		geminiNow = originalNow
	})

	geminiNow = func() time.Time { return time.Date(2026, 4, 5, 0, 0, 0, 0, time.UTC) }
	geminiGenerateContent = func(ctx context.Context, apiKey, model, prompt string) (string, error) {
		return `{"is_bookkeeping":true,"type":"expense","amount":180,"description":"午餐","category_id":"expense-food","category_name":"飲食","date":"2026-04-04","source_type":"","missing_fields":[]}`, nil
	}

	parser := NewGeminiParser("test-key")

	result, err := parser.Parse(t.Context(), "昨天午餐 180", testCategories())

	require.NoError(t, err)
	require.Equal(t, "2026-04-04", result.Date)
}

func TestGeminiParser_PromptContainsSourceType(t *testing.T) {
	prompt := buildGeminiPrompt("午餐 180", testCategories(), "2026-04-05")

	require.Contains(t, prompt, "source_type")
}

func TestGeminiParser_PromptContainsDateFormats(t *testing.T) {
	prompt := buildGeminiPrompt("午餐 180", testCategories(), "2026-04-05")

	require.Contains(t, prompt, "YYYY-MM-DD")
	require.Contains(t, prompt, "昨天")
	require.Contains(t, prompt, "前天")
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
