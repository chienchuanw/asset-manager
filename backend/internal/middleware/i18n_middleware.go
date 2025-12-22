// Package middleware provides HTTP middleware for the API
package middleware

import (
	"github.com/chienchuanw/asset-manager/internal/i18n"
	"github.com/gin-gonic/gin"
)

// LocaleKey is the context key for storing the locale
const LocaleKey = "locale"

// I18nMiddleware parses Accept-Language header and sets the locale in context
func I18nMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Accept-Language header
		acceptLanguage := c.GetHeader("Accept-Language")

		// Parse the header and get the best matching locale
		locale := i18n.ParseAcceptLanguage(acceptLanguage)

		// Store locale in context
		c.Set(LocaleKey, locale)

		c.Next()
	}
}

// GetLocale retrieves the locale from the Gin context
func GetLocale(c *gin.Context) i18n.Locale {
	if locale, exists := c.Get(LocaleKey); exists {
		if l, ok := locale.(i18n.Locale); ok {
			return l
		}
	}
	return i18n.DefaultLocale
}

// T is a helper function to get translated message from context
func T(c *gin.Context, key string) string {
	locale := GetLocale(c)
	return i18n.T(locale, key)
}

