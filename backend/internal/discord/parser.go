package discord

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

const defaultGeminiModel = "gemini-2.5-flash"

var (
	geminiNow             = time.Now
	geminiGenerateContent = defaultGeminiGenerateContent
)

// QueryParams holds filter criteria for query-type requests.
type QueryParams struct {
	Month    int    `json:"month"`    // 0 = current month
	Year     int    `json:"year"`     // 0 = current year
	Category string `json:"category"` // category name filter, "" = all
}

// ParseResult holds the structured result from NLP parsing.
type ParseResult struct {
	IsBookkeeping  bool         `json:"is_bookkeeping"`
	Action         string       `json:"action"` // "create" | "query" | "cc_payment" | "unsupported" | "chat" | ""
	Type           string       `json:"type"`
	Amount         float64      `json:"amount"`
	Description    string       `json:"description"`
	CategoryID     string       `json:"category_id"`
	CategoryName   string       `json:"category_name"`
	Date           string       `json:"date"`
	SourceType     string       `json:"source_type"`
	SourceID       string       `json:"-"`
	SourceName     string       `json:"-"`
	MissingFields  []string     `json:"missing_fields"`
	QueryType      string       `json:"query_type"` // "cash_flow_summary" | "account_balance" | ""
	QueryParams    *QueryParams `json:"query_params"`
	PaymentType    string       `json:"payment_type"`     // "full" | "minimum" | "custom" | ""
	TargetCardHint string       `json:"target_card_hint"` // credit card name hint for matching
}

// CategoryInfo represents a category passed to the parser for matching.
type CategoryInfo struct {
	ID   string
	Name string
	Type string
}

// Parser defines the interface for NLP message parsing.
type Parser interface {
	Parse(ctx context.Context, message string, categories []CategoryInfo) (*ParseResult, error)
}

// GeminiParser parses bookkeeping messages with Gemini.
type GeminiParser struct {
	apiKey string
	model  string
}

// NewGeminiParser creates a Gemini-backed parser.
func NewGeminiParser(apiKey string) *GeminiParser {
	return &GeminiParser{
		apiKey: apiKey,
		model:  defaultGeminiModel,
	}
}

// Parse converts a natural-language message into structured bookkeeping data.
func (p *GeminiParser) Parse(ctx context.Context, message string, categories []CategoryInfo) (*ParseResult, error) {
	if strings.TrimSpace(p.apiKey) == "" {
		return nil, errors.New("Gemini API key is not set")
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	raw, err := geminiGenerateContent(ctx, p.apiKey, p.model, buildGeminiPrompt(message, categories, geminiNow().Format("2006-01-02")))
	if err != nil {
		return nil, fmt.Errorf("gemini parse request failed: %w", err)
	}

	cleaned := cleanJSONResponse(raw)

	var result ParseResult
	if err := json.Unmarshal([]byte(cleaned), &result); err != nil {
		return nil, fmt.Errorf("invalid parser response: %w; raw=%q", err, snippet(cleaned, 160))
	}

	applyCategoryFallback(&result, categories)
	if result.Date == "" {
		result.Date = geminiNow().Format("2006-01-02")
	}
	if result.MissingFields == nil {
		result.MissingFields = []string{}
	}
	if result.Action == "" && result.IsBookkeeping {
		result.Action = "create"
	}
	applyQueryParamDefaults(&result)

	return &result, nil
}

func applyQueryParamDefaults(result *ParseResult) {
	if result.Action != "query" || result.QueryParams == nil {
		return
	}
	now := geminiNow()
	if result.QueryParams.Year == 0 {
		result.QueryParams.Year = now.Year()
	}
	if result.QueryParams.Month == 0 {
		result.QueryParams.Month = int(now.Month())
	}
}

func buildGeminiPrompt(message string, categories []CategoryInfo, today string) string {
	var categoryLines []string
	for _, category := range categories {
		categoryLines = append(categoryLines, fmt.Sprintf("- ID: %s, Name: %s, Type: %s", category.ID, category.Name, category.Type))
	}
	if len(categoryLines) == 0 {
		categoryLines = append(categoryLines, "- No categories provided")
	}

	return fmt.Sprintf(`Return a single JSON object matching exactly this schema:
{
  "is_bookkeeping": boolean,
  "action": "create" | "query" | "cc_payment" | "unsupported" | "chat" | "",
  "type": "income" | "expense" | "",
  "amount": number,
  "description": string,
  "category_id": string,
  "category_name": string,
  "date": string,
  "source_type": "cash" | "bank_account" | "credit_card" | "",
  "missing_fields": string[],
  "query_type": "cash_flow_summary" | "account_balance" | "",
  "query_params": {
    "month": number,
    "year": number,
    "category": string
  },
  "payment_type": "full" | "minimum" | "custom" | "",
  "target_card_hint": string
}

Rules:
- Support both Chinese and English input.
- Set "is_bookkeeping" to false for greetings, chit-chat, or messages that are not bookkeeping.
- For "action":
  - If the message contains a specific amount and is recording a transaction, use "create".
  - If the message is asking a question about spending, balance, or summary (interrogative form, no specific amount to record), use "query".
  - If the message is about paying a credit card bill (繳卡費/繳信用卡/pay credit card bill), use "cc_payment".
  - If the message has a clear action intent but is NOT bookkeeping, querying, or credit card payment (e.g., buying stocks, setting budgets), use "unsupported".
  - If the message is a greeting, chat, or casual conversation (嗨/你好/hello/thanks/謝謝), use "chat".
  - Use an empty string only when none of the above actions can be determined.
- Only use "income" or "expense" for "type" when the message is bookkeeping and the intent is transaction recording.
- Use date format YYYY-MM-DD.
- Default the date to today (%s) when the user does not specify one.
- Handle relative dates like 昨天 / yesterday, 前天 / day before yesterday, 上禮拜 / last week relative to today.
- Support absolute date formats: M/D, MM/DD, YYYY/MM/DD, and Chinese expressions like 4月3號.
- Match the best category from this list:
%s
- For source_type: infer the payment method from context clues.
  - 刷卡 / credit card → "credit_card"
  - 轉帳 / transfer / bank → "bank_account"
  - 現金 / cash → "cash"
  - If no payment method is mentioned, use an empty string.
- For "payment_type": only set it when action is "cc_payment".
  - Use "full" for full payment (全額).
  - Use "minimum" for minimum payment (最低).
  - Use "custom" for any other credit card payment amount.
- For "target_card_hint": extract the credit card name or bank keyword from the message when action is "cc_payment".
- If required bookkeeping information is missing for a transaction, keep "is_bookkeeping" true and list missing keys in "missing_fields".
- For "query_type":
  - Use "cash_flow_summary" for spending, income, or cash-flow questions.
  - Use "account_balance" for bank balance or credit card limit questions.
  - Otherwise use an empty string.
- For "query_params":
  - Resolve month/year from relative terms such as 這個月=current, 上個月=previous, last month=previous.
  - Extract the category name if mentioned.
  - Use 0 for unspecified month/year.
  - Use an empty string for an unspecified category.
- Use an empty string for unknown string fields and 0 for an unknown amount.
- Respond with JSON only. No markdown fences, no explanations.

User message:
%s`, today, strings.Join(categoryLines, "\n"), message)
}

func defaultGeminiGenerateContent(ctx context.Context, apiKey, modelName, prompt string) (string, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return "", err
	}
	defer client.Close()

	model := client.GenerativeModel(modelName)
	model.SystemInstruction = genai.NewUserContent(genai.Text("You are a bookkeeping assistant. Parse the user's message into structured data. Respond ONLY with valid JSON."))
	model.ResponseMIMEType = "application/json"

	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", err
	}

	text := extractTextFromResponse(resp)
	if strings.TrimSpace(text) == "" {
		return "", errors.New("empty Gemini response")
	}

	return text, nil
}

func extractTextFromResponse(resp *genai.GenerateContentResponse) string {
	if resp == nil {
		return ""
	}

	var parts []string
	for _, candidate := range resp.Candidates {
		if candidate == nil || candidate.Content == nil {
			continue
		}
		for _, part := range candidate.Content.Parts {
			switch v := part.(type) {
			case genai.Text:
				parts = append(parts, string(v))
			case fmt.Stringer:
				parts = append(parts, v.String())
			}
		}
	}

	return strings.TrimSpace(strings.Join(parts, "\n"))
}

func cleanJSONResponse(raw string) string {
	cleaned := strings.TrimSpace(raw)
	cleaned = strings.TrimPrefix(cleaned, "```json")
	cleaned = strings.TrimPrefix(cleaned, "```")
	cleaned = strings.TrimSuffix(cleaned, "```")
	return strings.TrimSpace(cleaned)
}

func applyCategoryFallback(result *ParseResult, categories []CategoryInfo) {
	if result == nil || !result.IsBookkeeping {
		return
	}

	categoryByID := make(map[string]CategoryInfo, len(categories))
	categoryByNameAndType := make(map[string]CategoryInfo, len(categories))
	for _, category := range categories {
		categoryByID[category.ID] = category
		categoryByNameAndType[category.Type+":"+category.Name] = category
	}

	if category, ok := categoryByID[result.CategoryID]; ok {
		if result.CategoryName == "" {
			result.CategoryName = category.Name
		}
		return
	}

	defaultName := defaultCategoryName(result.Type)
	if category, ok := categoryByNameAndType[result.Type+":"+defaultName]; ok {
		result.CategoryID = category.ID
		result.CategoryName = category.Name
		return
	}

	result.CategoryID = ""
	result.CategoryName = defaultName
}

func defaultCategoryName(categoryType string) string {
	if categoryType == "income" {
		return "其他收入"
	}
	return "其他支出"
}

func snippet(value string, max int) string {
	if max <= 0 || len(value) <= max {
		return value
	}
	return value[:max]
}
