package databases

import (
	"gorm.io/gorm"
)

type DB interface {
	DB() *gorm.DB
	Close() error
}
