// Package service proporciona funciones de formato y manipulación de cadenas.
package service

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

// PadLeft rellena s a la izquierda hasta width usando padChar.
func PadLeft(s string, width int, padChar rune) string {
	if utf8.RuneCountInString(s) >= width {
		return s
	}
	padCount := width - utf8.RuneCountInString(s)
	return strings.Repeat(string(padChar), padCount) + s
}

// PadRight rellena s a la derecha hasta width usando padChar.
func PadRight(s string, width int, padChar rune) string {
	if utf8.RuneCountInString(s) >= width {
		return s
	}
	padCount := width - utf8.RuneCountInString(s)
	return s + strings.Repeat(string(padChar), padCount)
}

// PadCenter centra s en un campo de ancho width usando padChar.
func PadCenter(s string, width int, padChar rune) string {
	length := utf8.RuneCountInString(s)
	if length >= width {
		return s
	}
	totalPad := width - length
	left := totalPad / 2
	right := totalPad - left
	return strings.Repeat(string(padChar), left) + s + strings.Repeat(string(padChar), right)
}

// Substr devuelve los primeros length runas de s.
func Substr(s string, length int) string {
	runes := []rune(s)
	if len(runes) <= length {
		return s
	}
	return string(runes[:length])
}

// FormatFloat formatea f con exactly decimals decimales.
func FormatFloat(f float64, decimals int) string {
	return fmt.Sprintf("%.*f", decimals, f)
}

// SplitString divide s en partes de length caracteres (runas).
func SplitString(s string, length int) []string {
	if length <= 0 {
		return nil
	}

	var result []string
	runes := []rune(s)

	for i := 0; i < len(runes); i += length {
		end := i + length
		if end > len(runes) {
			end = len(runes)
		}
		result = append(result, string(runes[i:end]))
	}
	return result
}

// CountChars cuenta el número de caracteres (runas) en s.
func CountChars(s string) int {
	count := 0
	for range s {
		count++
	}
	return count
}
