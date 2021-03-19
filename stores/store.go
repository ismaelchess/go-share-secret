package stores

import (
	"time"
)

type Store interface {
	Save(key string, data string, expires time.Duration) error
	Load(key string) (string, error)
	Delete(key string) error
}
