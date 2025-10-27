package domain

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

type Document struct {
	Number string `json:"number"`
}

func RestoreDocument(raw string) Document {
	return Document{
		Number: raw,
	}
}

func (p *Document) GetDocumentNumber() string {
	return p.Number
}

func NewDocument(document string) (*Document, error) {

	doc := &Document{
		Number: document,
	}

	if !doc.ValidateCpf() {
		return nil, errors.New("Document number is invalid")
	}

	return doc, nil
}

func (d *Document) NormalizeDocument(document string) string {
	re := regexp.MustCompile(`\D`)
	return re.ReplaceAllString(document, "")
}

func (d *Document) ValidateCpf() bool {

	if d.Number == "" {
		return false
	}

	cpf := d.NormalizeDocument(d.Number)

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

	return secondDigit == int(cpf[10]-'0')
}
