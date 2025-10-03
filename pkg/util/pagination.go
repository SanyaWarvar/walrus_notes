package util

import (
	"wn/pkg/constants"
)

func CalculateOffset(page int) int {
	return (page - 1) * constants.PageSize
}

func CalculateLimit() int {
	return constants.PageSize
}
