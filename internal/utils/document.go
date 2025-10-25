package utils

import (
	"regexp"
	"strconv"
	"strings"
)

func ValidateCpf(cpf string) bool {

	if cpf == "" {
		return false
	}

	re := regexp.MustCompile(`\D`)
	cpf = re.ReplaceAllString(cpf, "")

	if len(cpf) != 11 {
		return false
	}

	if strings.Repeat(string(cpf[0]), 11) == cpf {
		return false
	}

	sum := 0

	for i := 0; i < 9; i++ {
		num, _ := strconv.Atoi(string(cpf[i]))
		sum += num * (10 - i)
	}

	firstDigit := (sum * 10) % 11

	if firstDigit == 10 {
		firstDigit = 0
	}

	if firstDigit != int(cpf[9]-'0') {
		return false
	}

	sum = 0

	for i := 0; i < 10; i++ {
		num, _ := strconv.Atoi(string(cpf[i]))
		sum += num * (11 - i)
	}

	secondDigit := (sum * 10) % 11

	if secondDigit == 10 {
		secondDigit = 0
	}

	if secondDigit != int(cpf[10]-'0') {
		return false
	}

	return true
}
