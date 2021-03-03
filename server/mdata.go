package main

import (
	"path/filepath"
	"text/template"
	"time"

	"github.com/google/uuid"
)

type sdata struct {
	Value string
	Unit  string
	Utime int
}

type PathHost struct {
	Port string
	Host string
}

type Result struct {
	Data string
}

func (s *sdata) getUniqueId() string {
	return uuid.New().String()
}

func (s *sdata) expirationDate() time.Duration {
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
	secret, err := filepath.Abs("../ui/secret.html")
	if err != nil {
		return nil, err
	}
	return template.ParseFiles(secret)
}
