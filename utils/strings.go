package utils

import (
	"golang.org/x/crypto/bcrypt"
	"strings"
	"unicode"
)

// CleanEmail email adresini temizler ve küçük harfe çevirir
func CleanEmail(email string) string {
	return strings.ToLower(strings.TrimSpace(email))
}

// CleanPhone telefon numarasını temizler
func CleanPhone(phone string) string {
	return strings.TrimSpace(phone)
}

// ToTitle metni title case yapar (Her Kelimenin İlk Harfi Büyük)
func ToTitle(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}

	words := strings.Fields(strings.ToLower(s))
	for i, word := range words {
		runes := []rune(word)
		if len(runes) > 0 {
			runes[0] = unicode.ToUpper(runes[0])
		}
		words[i] = string(runes)
	}

	return strings.Join(words, " ")
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
