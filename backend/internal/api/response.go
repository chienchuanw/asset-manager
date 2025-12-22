// Package api provides HTTP handlers for the API
package api

import (
	"net/http"

	"github.com/chienchuanw/asset-manager/internal/i18n"
	"github.com/chienchuanw/asset-manager/internal/middleware"
	"github.com/gin-gonic/gin"
)

// LocalizedAPIResponse represents a localized API response
type LocalizedAPIResponse struct {
	Data  interface{}       `json:"data,omitempty"`
	Error *LocalizedAPIError `json:"error,omitempty"`
}

// LocalizedAPIError represents a localized API error
type LocalizedAPIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// RespondError sends a localized error response
func RespondError(c *gin.Context, statusCode int, errorCode string, fallbackMessage string) {
	locale := middleware.GetLocale(c)
	message := i18n.T(locale, errorCode)

	// If no translation found, use fallback message
	if message == errorCode {
		message = fallbackMessage
	}

	c.JSON(statusCode, LocalizedAPIResponse{
		Error: &LocalizedAPIError{
			Code:    errorCode,
			Message: message,
		},
	})
}

// RespondErrorWithDetails sends a localized error response with additional details
func RespondErrorWithDetails(c *gin.Context, statusCode int, errorCode string, details string) {
	locale := middleware.GetLocale(c)
	message := i18n.T(locale, errorCode)

	// If no translation found, use the error code as message
	if message == errorCode {
		message = details
	} else if details != "" {
		// Append details to the translated message
		message = message + ": " + details
	}

	c.JSON(statusCode, LocalizedAPIResponse{
		Error: &LocalizedAPIError{
			Code:    errorCode,
			Message: message,
		},
	})
}

// RespondSuccess sends a success response with data
func RespondSuccess(c *gin.Context, statusCode int, data interface{}) {
	c.JSON(statusCode, LocalizedAPIResponse{
		Data: data,
	})
}

// RespondBadRequest sends a 400 Bad Request error
func RespondBadRequest(c *gin.Context, errorCode string, details string) {
	RespondErrorWithDetails(c, http.StatusBadRequest, errorCode, details)
}

// RespondUnauthorized sends a 401 Unauthorized error
func RespondUnauthorized(c *gin.Context, errorCode string, details string) {
	RespondErrorWithDetails(c, http.StatusUnauthorized, errorCode, details)
}

// RespondForbidden sends a 403 Forbidden error
func RespondForbidden(c *gin.Context, errorCode string, details string) {
	RespondErrorWithDetails(c, http.StatusForbidden, errorCode, details)
}

// RespondNotFound sends a 404 Not Found error
func RespondNotFound(c *gin.Context, errorCode string, details string) {
	RespondErrorWithDetails(c, http.StatusNotFound, errorCode, details)
}

// RespondInternalError sends a 500 Internal Server Error
func RespondInternalError(c *gin.Context, errorCode string, details string) {
	RespondErrorWithDetails(c, http.StatusInternalServerError, errorCode, details)
}

