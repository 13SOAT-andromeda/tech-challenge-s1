package email

import (
	"os"
	"path/filepath"
)

type EmailTemplates string

const (
	ORDER_APPROVAL EmailTemplates = "order-approval"
)

func LoadTemplate(template EmailTemplates) (string, error) {
	absPath, err := filepath.Abs("../../internal/adapter/email/templates/" + string(template) + ".html")

	if err != nil {
		return "", ErrTemplateNotFound
	}

	content, err := os.ReadFile(absPath)

	if err != nil {
		return "", ErrTemplateLoad
	}

	return string(content[:]), nil
}
