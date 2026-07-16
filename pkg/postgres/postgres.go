// pkg/postgres/postgres.go
package postgres

import (
	"fmt"
	"dan-ai/pkg/config"

	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

// Connect opens a connection to PostgreSQL and returns a GORM DB client.
func Connect(cfg *config.Config) (*gorm.DB, error) {
	var dsn string
	if cfg.DB.URL != "" {
		dsn = cfg.DB.URL
	} else {
		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Jakarta",
			cfg.DB.Host,
			cfg.DB.User,
			cfg.DB.Password,
			cfg.DB.Name,
			cfg.DB.Port,
		)
	}

	gormCfg := &gorm.Config{}

	if cfg.App.Env == "production" {
		gormCfg.Logger = gormlogger.Default.LogMode(gormlogger.Error)
	} else {
		gormCfg.Logger = gormlogger.Default.LogMode(gormlogger.Info)
	}

	db, err := gorm.Open(gormpostgres.Open(dsn), gormCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	return db, nil
}
