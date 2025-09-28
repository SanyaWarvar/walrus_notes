package util

import (
	"encoding/hex"
	"strconv"
	"strings"
)

func HexadecimalWithPadding(n int) string {
	hexadecimal := strconv.FormatInt(int64(n), 16)
	padLength := 6 - len(hexadecimal)
	if padLength > 0 {
		hexadecimal = strings.Repeat("#", padLength) + hexadecimal
	}
	return hexadecimal
}

func IsHexNumber(s string) bool {
	_, err := hex.DecodeString(s)
	return err == nil
}
