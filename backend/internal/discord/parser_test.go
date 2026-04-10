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
		return &ParseResult{Action: "chat", IsBookkeeping: false}, nil
	case "繳中信卡 15000":
		return &ParseResult{Action: "cc_payment", Amount: 15000, PaymentType: "custom", TargetCardHint: "中信", Date: "2026-04-05"}, nil
	case "繳玉山卡全額":
		return &ParseResult{Action: "cc_payment", PaymentType: "full", TargetCardHint: "玉山", Date: "2026-04-05"}, nil
	case "繳中信卡最低 3000":
		return &ParseResult{Action: "cc_payment", Amount: 3000, PaymentType: "minimum", TargetCardHint: "中信", Date: "2026-04-05"}, nil
	case "繳中信卡":
		return &ParseResult{Action: "cc_payment", PaymentType: "custom", TargetCardHint: "中信", Date: "2026-04-05", MissingFields: []string{"amount"}}, nil
	case "pay credit card 15000":
		return &ParseResult{Action: "cc_payment", Amount: 15000, PaymentType: "custom", Date: "2026-04-05"}, nil
	case "幫我買台積電 10 股":
		return &ParseResult{Action: "unsupported", IsBookkeeping: false}, nil
	case "幫我設定每月預算 30000":
		return &ParseResult{Action: "unsupported", IsBookkeeping: false}, nil
	case "buy 10 shares of TSMC":
		return &ParseResult{Action: "unsupported", IsBookkeeping: false}, nil
	case "嗨":
		return &ParseResult{Action: "chat", IsBookkeeping: false}, nil
	case "hello":
		return &ParseResult{Action: "chat", IsBookkeeping: false}, nil
	case "謝謝":
		return &ParseResult{Action: "chat", IsBookkeeping: false}, nil
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
	require.Equal(t, "chat", result.Action)
}

func TestMockParser_CCPayment_CustomAmount(t *testing.T) {
	parser := MockParser{}

	result, err := parser.Parse(t.Context(), "繳中信卡 15000", testCategories())

	require.NoError(t, err)
	require.Equal(t, "cc_payment", result.Action)
	require.Equal(t, 15000.0, result.Amount)
	require.Equal(t, "custom", result.PaymentType)
	require.Contains(t, result.TargetCardHint, "中信")
}

func TestMockParser_CCPayment_FullPayment(t *testing.T) {
	parser := MockParser{}

	result, err := parser.Parse(t.Context(), "繳玉山卡全額", testCategories())

	require.NoError(t, err)
	require.Equal(t, "cc_payment", result.Action)
	require.Equal(t, "full", result.PaymentType)
}

func TestMockParser_CCPayment_MinimumPayment(t *testing.T) {
	parser := MockParser{}

	result, err := parser.Parse(t.Context(), "繳中信卡最低 3000", testCategories())

	require.NoError(t, err)
	require.Equal(t, "cc_payment", result.Action)
	require.Equal(t, 3000.0, result.Amount)
	require.Equal(t, "minimum", result.PaymentType)
}

func TestMockParser_CCPayment_MissingAmount(t *testing.T) {
	parser := MockParser{}

	result, err := parser.Parse(t.Context(), "繳中信卡", testCategories())

	require.NoError(t, err)
	require.Equal(t, "cc_payment", result.Action)
	require.Equal(t, []string{"amount"}, result.MissingFields)
}

func TestMockParser_CCPayment_English(t *testing.T) {
	parser := MockParser{}

	result, err := parser.Parse(t.Context(), "pay credit card 15000", testCategories())

	require.NoError(t, err)
	require.Equal(t, "cc_payment", result.Action)
}

func TestMockParser_BackwardCompat_CreateAction(t *testing.T) {
	parser := MockParser{}

	result, err := parser.Parse(t.Context(), "午餐吃拉麵 180", testCategories())

	require.NoError(t, err)
	require.Equal(t, "", result.PaymentType)
	require.Equal(t, "", result.TargetCardHint)
}

func TestMockParser_Unsupported_BuyStock(t *testing.T) {
	parser := MockParser{}

	result, err := parser.Parse(t.Context(), "幫我買台積電 10 股", testCategories())

	require.NoError(t, err)
	require.Equal(t, "unsupported", result.Action)
}

func TestMockParser_Unsupported_SetBudget(t *testing.T) {
	parser := MockParser{}

	result, err := parser.Parse(t.Context(), "幫我設定每月預算 30000", testCategories())

	require.NoError(t, err)
	require.Equal(t, "unsupported", result.Action)
}

func TestMockParser_Unsupported_English(t *testing.T) {
	parser := MockParser{}

	result, err := parser.Parse(t.Context(), "buy 10 shares of TSMC", testCategories())

	require.NoError(t, err)
	require.Equal(t, "unsupported", result.Action)
}

func TestMockParser_Chat_Greeting_ZhTW(t *testing.T) {
	parser := MockParser{}

	result, err := parser.Parse(t.Context(), "嗨", testCategories())

	require.NoError(t, err)
	require.Equal(t, "chat", result.Action)
}

func TestMockParser_Chat_Greeting_En(t *testing.T) {
	parser := MockParser{}

	result, err := parser.Parse(t.Context(), "hello", testCategories())

	require.NoError(t, err)
	require.Equal(t, "chat", result.Action)
}

func TestMockParser_Chat_Thanks(t *testing.T) {
	parser := MockParser{}

	result, err := parser.Parse(t.Context(), "謝謝", testCategories())

	require.NoError(t, err)
	require.Equal(t, "chat", result.Action)
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

func TestGeminiParser_CCPayment(t *testing.T) {
	originalGenerate := geminiGenerateContent
	originalNow := geminiNow
	t.Cleanup(func() {
		geminiGenerateContent = originalGenerate
		geminiNow = originalNow
	})

	geminiNow = func() time.Time { return time.Date(2026, 4, 5, 0, 0, 0, 0, time.UTC) }
	geminiGenerateContent = func(ctx context.Context, apiKey, model, prompt string) (string, error) {
		return `{"is_bookkeeping":false,"action":"cc_payment","type":"","amount":15000,"description":"","category_id":"","category_name":"","date":"2026-04-05","source_type":"","missing_fields":[],"payment_type":"custom","target_card_hint":"中信"}`, nil
	}

	parser := NewGeminiParser("test-key")

	result, err := parser.Parse(t.Context(), "繳中信卡 15000", testCategories())

	require.NoError(t, err)
	require.Equal(t, "cc_payment", result.Action)
	require.Equal(t, 15000.0, result.Amount)
	require.Equal(t, "custom", result.PaymentType)
	require.Equal(t, "中信", result.TargetCardHint)
	require.Empty(t, result.MissingFields)
}

func TestGeminiParser_Unsupported(t *testing.T) {
	originalGenerate := geminiGenerateContent
	originalNow := geminiNow
	t.Cleanup(func() {
		geminiGenerateContent = originalGenerate
		geminiNow = originalNow
	})

	geminiNow = func() time.Time { return time.Date(2026, 4, 5, 0, 0, 0, 0, time.UTC) }
	geminiGenerateContent = func(ctx context.Context, apiKey, model, prompt string) (string, error) {
		return `{"is_bookkeeping":false,"action":"unsupported","type":"","amount":0,"description":"","category_id":"","category_name":"","date":"2026-04-05","source_type":"","missing_fields":[]}`, nil
	}

	parser := NewGeminiParser("test-key")

	result, err := parser.Parse(t.Context(), "幫我買台積電 10 股", testCategories())

	require.NoError(t, err)
	require.Equal(t, "unsupported", result.Action)
	require.Empty(t, result.MissingFields)
}

func TestGeminiParser_Chat(t *testing.T) {
	originalGenerate := geminiGenerateContent
	originalNow := geminiNow
	t.Cleanup(func() {
		geminiGenerateContent = originalGenerate
		geminiNow = originalNow
	})

	geminiNow = func() time.Time { return time.Date(2026, 4, 5, 0, 0, 0, 0, time.UTC) }
	geminiGenerateContent = func(ctx context.Context, apiKey, model, prompt string) (string, error) {
		return `{"is_bookkeeping":false,"action":"chat","type":"","amount":0,"description":"hello","category_id":"","category_name":"","date":"2026-04-05","source_type":"","missing_fields":[]}`, nil
	}

	parser := NewGeminiParser("test-key")

	result, err := parser.Parse(t.Context(), "hello", testCategories())

	require.NoError(t, err)
	require.Equal(t, "chat", result.Action)
	require.Empty(t, result.MissingFields)
}

func TestGeminiParser_PromptContainsSourceType(t *testing.T) {
	prompt := buildGeminiPrompt("午餐 180", testCategories(), "2026-04-05")

	require.Contains(t, prompt, "source_type")
}

func TestBuildGeminiPrompt_IncludesCCPaymentFields(t *testing.T) {
	prompt := buildGeminiPrompt("test", nil, "2026-04-06")

	require.Contains(t, prompt, "payment_type")
	require.Contains(t, prompt, "target_card_hint")
	require.Contains(t, prompt, "cc_payment")
	require.Contains(t, prompt, "unsupported")
	require.Contains(t, prompt, "chat")
}

func TestBuildGeminiPrompt_IncludesQueryFields(t *testing.T) {
	prompt := buildGeminiPrompt("test", nil, "2026-04-06")

	require.Contains(t, prompt, "action")
	require.Contains(t, prompt, "query_type")
	require.Contains(t, prompt, "query_params")
}

func TestGeminiParser_PromptContainsDateFormats(t *testing.T) {
	prompt := buildGeminiPrompt("午餐 180", testCategories(), "2026-04-05")

	require.Contains(t, prompt, "YYYY-MM-DD")
	require.Contains(t, prompt, "昨天")
	require.Contains(t, prompt, "前天")
}

func TestParseResult_ActionDefaults(t *testing.T) {
	originalGenerate := geminiGenerateContent
	originalNow := geminiNow
	t.Cleanup(func() {
		geminiGenerateContent = originalGenerate
		geminiNow = originalNow
	})

	geminiNow = func() time.Time { return time.Date(2026, 4, 6, 0, 0, 0, 0, time.UTC) }

	tests := []struct {
		name           string
		response       string
		expectedAction string
		expectedMonth  int
		expectedYear   int
	}{
		{
			name:           "Given bookkeeping result with empty action When Parse Then defaults action to create",
			response:       `{"is_bookkeeping":true,"action":"","type":"expense","amount":180,"description":"午餐","category_id":"expense-food","category_name":"飲食","date":"2026-04-06","source_type":"cash","missing_fields":[]}`,
			expectedAction: "create",
		},
		{
			name:           "Given query result with missing year When Parse Then defaults year to current year",
			response:       `{"is_bookkeeping":true,"action":"query","type":"","amount":0,"description":"","category_id":"","category_name":"","date":"2026-04-06","source_type":"","missing_fields":[],"query_type":"cash_flow_summary","query_params":{"month":3,"year":0,"category":""}}`,
			expectedAction: "query",
			expectedMonth:  3,
			expectedYear:   2026,
		},
		{
			name:           "Given query result with missing month When Parse Then defaults month to current month",
			response:       `{"is_bookkeeping":true,"action":"query","type":"","amount":0,"description":"","category_id":"","category_name":"","date":"2026-04-06","source_type":"","missing_fields":[],"query_type":"cash_flow_summary","query_params":{"month":0,"year":2025,"category":""}}`,
			expectedAction: "query",
			expectedMonth:  4,
			expectedYear:   2025,
		},
		{
			name:           "Given non-bookkeeping result with empty action When Parse Then action stays empty",
			response:       `{"is_bookkeeping":false,"action":"","type":"","amount":0,"description":"hello","category_id":"","category_name":"","date":"2026-04-06","source_type":"","missing_fields":[]}`,
			expectedAction: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			geminiGenerateContent = func(ctx context.Context, apiKey, model, prompt string) (string, error) {
				return tt.response, nil
			}

			parser := NewGeminiParser("test-key")

			result, err := parser.Parse(t.Context(), "test message", testCategories())

			require.NoError(t, err)
			require.Equal(t, tt.expectedAction, result.Action)

			if tt.expectedAction == "query" {
				require.NotNil(t, result.QueryParams)
				require.Equal(t, tt.expectedYear, result.QueryParams.Year)
				require.Equal(t, tt.expectedMonth, result.QueryParams.Month)
			} else {
				require.Nil(t, result.QueryParams)
			}
		})
	}
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

func TestGeminiNow_ReturnsAsiaTaipeiTime(t *testing.T) {
	// geminiNow() should return time in Asia/Taipei timezone
	now := geminiNow()
	taipei, _ := time.LoadLocation("Asia/Taipei")

	require.Equal(t, taipei, now.Location(), "geminiNow should return Asia/Taipei time")
}

func TestParse_DateInPromptReflectsTimezone(t *testing.T) {
	originalGenerate := geminiGenerateContent
	originalNow := geminiNow
	t.Cleanup(func() {
		geminiGenerateContent = originalGenerate
		geminiNow = originalNow
	})

	// UTC: 2026-04-05 23:00 → Asia/Taipei: 2026-04-06 07:00
	geminiNow = func() time.Time {
		taipei, _ := time.LoadLocation("Asia/Taipei")
		return time.Date(2026, 4, 5, 23, 0, 0, 0, time.UTC).In(taipei)
	}

	var capturedPrompt string
	geminiGenerateContent = func(ctx context.Context, apiKey, model, prompt string) (string, error) {
		capturedPrompt = prompt
		return `{"is_bookkeeping":true,"type":"expense","amount":180,"description":"午餐","category_id":"expense-food","category_name":"飲食","date":"2026-04-06","source_type":"","missing_fields":[]}`, nil
	}

	parser := NewGeminiParser("test-key")
	_, err := parser.Parse(t.Context(), "午餐 180", testCategories())

	require.NoError(t, err)
	require.Contains(t, capturedPrompt, "2026-04-06", "prompt should contain Asia/Taipei date, not UTC date")
	require.NotContains(t, capturedPrompt, "today (2026-04-05)", "prompt should not contain UTC date")
}

func TestApplyQueryParamDefaults_UsesTimezoneMonth(t *testing.T) {
	originalGenerate := geminiGenerateContent
	originalNow := geminiNow
	t.Cleanup(func() {
		geminiGenerateContent = originalGenerate
		geminiNow = originalNow
	})

	// UTC: 2026-03-31 22:00 → Asia/Taipei: 2026-04-01 06:00
	geminiNow = func() time.Time {
		taipei, _ := time.LoadLocation("Asia/Taipei")
		return time.Date(2026, 3, 31, 22, 0, 0, 0, time.UTC).In(taipei)
	}

	geminiGenerateContent = func(ctx context.Context, apiKey, model, prompt string) (string, error) {
		return `{"is_bookkeeping":true,"action":"query","type":"","amount":0,"description":"","category_id":"","category_name":"","date":"2026-04-01","source_type":"","missing_fields":[],"query_type":"cash_flow_summary","query_params":{"month":0,"year":0,"category":""}}`, nil
	}

	parser := NewGeminiParser("test-key")
	result, err := parser.Parse(t.Context(), "這個月花了多少", testCategories())

	require.NoError(t, err)
	require.Equal(t, "query", result.Action)
	require.NotNil(t, result.QueryParams)
	require.Equal(t, 4, result.QueryParams.Month, "month should be 4 (April in Asia/Taipei), not 3 (March in UTC)")
	require.Equal(t, 2026, result.QueryParams.Year)
}
