package i18n

import (
	"testing"
)

// TestInitLoadsAllTranslations 測試 Init() 函式是否正確載入所有翻譯
func TestInitLoadsAllTranslations(t *testing.T) {
	// 重置全域狀態以便測試
	initialized = false
	translations = make(map[Locale]map[string]string)

	err := Init()
	if err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// 驗證兩種語言都已載入
	if _, ok := translations[LocaleZhTW]; !ok {
		t.Error("zh-TW translations not loaded")
	}
	if _, ok := translations[LocaleEn]; !ok {
		t.Error("en translations not loaded")
	}

	// 驗證翻譯數量合理（應該有 67 個 key）
	zhTWCount := len(translations[LocaleZhTW])
	enCount := len(translations[LocaleEn])

	if zhTWCount == 0 {
		t.Error("zh-TW translations are empty")
	}
	if enCount == 0 {
		t.Error("en translations are empty")
	}

	// 驗證兩種語言的 key 數量相同
	if zhTWCount != enCount {
		t.Errorf("Translation count mismatch: zh-TW=%d, en=%d", zhTWCount, enCount)
	}

	t.Logf("Loaded %d translations for each language", zhTWCount)
}

// TestInitIdempotent 測試 Init() 是否可以安全地多次呼叫
func TestInitIdempotent(t *testing.T) {
	initialized = false
	translations = make(map[Locale]map[string]string)

	err1 := Init()
	if err1 != nil {
		t.Fatalf("First Init() failed: %v", err1)
	}

	firstCount := len(translations[LocaleZhTW])

	err2 := Init()
	if err2 != nil {
		t.Fatalf("Second Init() failed: %v", err2)
	}

	secondCount := len(translations[LocaleZhTW])

	if firstCount != secondCount {
		t.Errorf("Init() is not idempotent: first=%d, second=%d", firstCount, secondCount)
	}
}

// TestTReturnsCorrectTranslation 測試 T() 函式是否返回正確的翻譯
func TestTReturnsCorrectTranslation(t *testing.T) {
	initialized = false
	translations = make(map[Locale]map[string]string)

	if err := Init(); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	tests := []struct {
		locale   Locale
		key      string
		expected string
	}{
		{LocaleZhTW, "INVALID_REQUEST", "請求格式錯誤"},
		{LocaleEn, "INVALID_REQUEST", "Invalid request format"},
		{LocaleZhTW, "CATEGORY_NOT_FOUND", "找不到分類"},
		{LocaleEn, "CATEGORY_NOT_FOUND", "Category not found"},
		{LocaleZhTW, "TRANSACTION_CREATE_FAILED", "建立交易記錄失敗"},
		{LocaleEn, "TRANSACTION_CREATE_FAILED", "Failed to create transaction"},
	}

	for _, tt := range tests {
		result := T(tt.locale, tt.key)
		if result != tt.expected {
			t.Errorf("T(%s, %q) = %q, want %q", tt.locale, tt.key, result, tt.expected)
		}
	}
}

// TestTFallbackToDefaultLocale 測試 T() 函式是否正確 fallback 到預設語言
func TestTFallbackToDefaultLocale(t *testing.T) {
	initialized = false
	translations = make(map[Locale]map[string]string)

	if err := Init(); err != nil {
		t.Fatalf("Init() failed: %v", err)
	}

	// 測試不存在的 key，應該返回 key 本身
	result := T(LocaleEn, "NONEXISTENT_KEY")
	if result != "NONEXISTENT_KEY" {
		t.Errorf("T() should return key for missing translation, got %q", result)
	}
}

// TestParseLocale 測試 ParseLocale() 函式
func TestParseLocale(t *testing.T) {
	tests := []struct {
		input    string
		expected Locale
	}{
		{"en", LocaleEn},
		{"en-US", LocaleEn},
		{"en-GB", LocaleEn},
		{"zh-TW", LocaleZhTW},
		{"zh-Hant", LocaleZhTW},
		{"zh", LocaleZhTW},
		{"unknown", DefaultLocale},
		{"", DefaultLocale},
	}

	for _, tt := range tests {
		result := ParseLocale(tt.input)
		if result != tt.expected {
			t.Errorf("ParseLocale(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

// TestParseAcceptLanguage 測試 ParseAcceptLanguage() 函式
func TestParseAcceptLanguage(t *testing.T) {
	tests := []struct {
		input    string
		expected Locale
	}{
		{"en-US,en;q=0.9", LocaleEn},
		{"zh-TW,zh;q=0.9", LocaleZhTW},
		{"fr-FR,fr;q=0.9", DefaultLocale},
		{"", DefaultLocale},
	}

	for _, tt := range tests {
		result := ParseAcceptLanguage(tt.input)
		if result != tt.expected {
			t.Errorf("ParseAcceptLanguage(%q) = %q, want %q", tt.input, result, tt.expected)
		}
	}
}

