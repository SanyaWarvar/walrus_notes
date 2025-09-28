package util

import (
	"wn/pkg/constants"
	"net/http"
	"strings"
)

var headerToMask = []string{constants.AuthorizationHeader, constants.RefreshHeader}

func MaskHeaders(header http.Header) http.Header {
	// Копируем исходные заголовки
	maskedHeaders := make(http.Header)
	for k, v := range header {
		maskedHeaders[k] = v
	}

	// Маскирование заголовков из слайса headerToMask
	for _, headerName := range headerToMask {
		if value := header.Get(headerName); value != "" {
			maskedHeaders.Set(headerName, maskBySize(value)) // Замените на то, что нужно для маскировки
		}
	}

	return maskedHeaders
}

func maskBySize(val string) string {
	if len(val) < 1 {
		return ""
	} else if len(val) < 4 {
		return "***"
	} else if len(val) < 6 {
		return val[:3] + strings.Repeat("*", len(val)-3)
	} else if len(val) < 10 {
		return val[:4] + strings.Repeat("*", len(val)-4)
	} else if len(val) < 16 {
		return val[:5] + strings.Repeat("*", len(val)-5)
	}
	return val[:3] + "***" + val[len(val)-3:]
}
