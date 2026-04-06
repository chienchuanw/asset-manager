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

// ParseResult holds the structured result from NLP parsing.
type ParseResult struct {
	IsBookkeeping bool     `json:"is_bookkeeping"`
	Type          string   `json:"type"`
	Amount        float64  `json:"amount"`
	Description   string   `json:"description"`
	CategoryID    string   `json:"category_id"`
	CategoryName  string   `json:"category_name"`
	Date          string   `json:"date"`
	SourceType    string   `json:"source_type"`
	SourceID      string   `json:"-"`
	SourceName    string   `json:"-"`
	MissingFields []string `json:"missing_fields"`
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

	return &result, nil
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
  "type": "income" | "expense" | "",
  "amount": number,
  "description": string,
  "category_id": string,
  "category_name": string,
  "date": string,
  "source_type": "cash" | "bank_account" | "credit_card" | "",
  "missing_fields": string[]
}

Rules:
- Support both Chinese and English input.
- Set "is_bookkeeping" to false for greetings, chit-chat, or messages that are not bookkeeping.
- Only use "income" or "expense" for "type" when the message is bookkeeping.
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
- If required bookkeeping information is missing, keep "is_bookkeeping" true and list missing keys in "missing_fields".
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
