// Package i18n provides internationalization support for error messages
package i18n

import (
	"embed"
	"encoding/json"
	"io/fs"
	"sync"
)

//go:embed locales/*/*.json
var localesFS embed.FS

// Locale represents a supported language
type Locale string

const (
	LocaleZhTW Locale = "zh-TW"
	LocaleEn   Locale = "en"
)

// DefaultLocale is the default language
const DefaultLocale = LocaleZhTW

// translations stores all loaded translations
var (
	translations = make(map[Locale]map[string]string)
	mu           sync.RWMutex
	initialized  bool
)

// Init loads all translation files from locale subdirectories
func Init() error {
	mu.Lock()
	defer mu.Unlock()

	if initialized {
		return nil
	}

	locales := []Locale{LocaleZhTW, LocaleEn}

	for _, locale := range locales {
		// 初始化該語言的翻譯字典
		translations[locale] = make(map[string]string)

		// 掃描該語言的目錄下的所有 JSON 檔案
		localeDir := "locales/" + string(locale)
		entries, err := fs.ReadDir(localesFS, localeDir)
		if err != nil {
			return err
		}

		// 合併所有 JSON 檔案中的翻譯
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			// 只處理 .json 檔案
			if len(entry.Name()) < 5 || entry.Name()[len(entry.Name())-5:] != ".json" {
				continue
			}

			filePath := localeDir + "/" + entry.Name()
			data, err := localesFS.ReadFile(filePath)
			if err != nil {
				return err
			}

			var msgs map[string]string
			if err := json.Unmarshal(data, &msgs); err != nil {
				return err
			}

			// 合併到該語言的翻譯字典
			for key, value := range msgs {
				translations[locale][key] = value
			}
		}
	}

	initialized = true
	return nil
}

// T returns the translated message for the given key and locale
func T(locale Locale, key string) string {
	mu.RLock()
	defer mu.RUnlock()

	// Try to get translation for the requested locale
	if msgs, ok := translations[locale]; ok {
		if msg, ok := msgs[key]; ok {
			return msg
		}
	}

	// Fallback to default locale
	if msgs, ok := translations[DefaultLocale]; ok {
		if msg, ok := msgs[key]; ok {
			return msg
		}
	}

	// Return the key itself if no translation found
	return key
}

// ParseLocale parses a locale string and returns the corresponding Locale
func ParseLocale(s string) Locale {
	switch s {
	case "en", "en-US", "en-GB":
		return LocaleEn
	case "zh-TW", "zh-Hant", "zh":
		return LocaleZhTW
	default:
		return DefaultLocale
	}
}

// ParseAcceptLanguage parses the Accept-Language header and returns the best match
func ParseAcceptLanguage(header string) Locale {
	if header == "" {
		return DefaultLocale
	}

	// Simple parsing: just check for "en" or "zh"
	// A more robust implementation would parse quality values
	for i := 0; i < len(header); i++ {
		if i+2 <= len(header) {
			prefix := header[i : i+2]
			if prefix == "en" {
				return LocaleEn
			}
			if prefix == "zh" {
				return LocaleZhTW
			}
		}
	}

	return DefaultLocale
}

