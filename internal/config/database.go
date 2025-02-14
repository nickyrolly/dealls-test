package config

import (
	"log"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type DatabaseOption struct {
	Driver          string
	Host            string
	Port            int
	SSLMode         bool
	DBName          string
	Username        string
	Password        string
	MaxIdleConns    int
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
}

func NewDatabase(cfg DatabaseOption) *gorm.DB {
	if cfg.Host == "" {
		cfg.Host = "localhost"
	}

	if cfg.MaxIdleConns == 0 {
		cfg.MaxIdleConns = 5
	}

	if cfg.MaxOpenConns == 0 {
		cfg.MaxOpenConns = 10
	}

	if cfg.ConnMaxLifetime == 0 {
		cfg.ConnMaxLifetime = time.Hour
	}

	switch dbDriver := cfg.Driver; dbDriver {
	case "postgresql":
		//return newPostgreSQL(cfg)
		return nil
	case "pgx":
		//return newPostgreSQL(cfg)
		return nil
	default:
		return newSQLLite(cfg)
	}
}

func newSQLLite(config DatabaseOption) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(config.DBName), &gorm.Config{
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
	})
	if err != nil {
		log.Fatal("Failed to connect to database")
	}
	return db
}
