package main

import (
	"time"

	"github.com/google/uuid"
)

type sdata struct {
	Value string
	Unit  string
	Utime int
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
