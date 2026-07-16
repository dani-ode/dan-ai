// apps/api/bootstrap/database.go
package bootstrap

import (
	"dan-ai/pkg/config"
	"dan-ai/pkg/postgres"

	"gorm.io/gorm"
)

// NewDatabase connects to PostgreSQL database using the provided configuration.
func NewDatabase(cfg *config.Config) (*gorm.DB, error) {
	return postgres.Connect(cfg)
}
