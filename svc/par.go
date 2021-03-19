package svc

import (
	"path/filepath"
	"text/template"
	"time"

	"github.com/google/uuid"
)

type Par struct {
	Value string
	Unit  string
	Utime int
}

type Result struct {
	Data string
}

func (s *Par) getUniqueId() string {
	return uuid.New().String()
}

func (s *Par) expirationDate() time.Duration {
	switch s.Unit {
	case "h":
		return time.Duration(s.Utime) * time.Hour
	case "d":
		return time.Duration(s.Utime*24) * time.Hour
	}
	return time.Duration(s.Utime) * time.Minute
}

type TemplateParser interface {
	ParseFiles() (*template.Template, error)
}

type DefaultTemplateParser struct{}

func (r *DefaultTemplateParser) ParseFiles() (*template.Template, error) {
	return template.ParseFiles("./ui/secret.html")
}

type TestTemplateParser struct{}

func (r *TestTemplateParser) ParseFiles() (*template.Template, error) {
	secret, err := filepath.Abs("./ui/secret.html")
	if err != nil {
		return nil, err
	}
	return template.ParseFiles(secret)
}
