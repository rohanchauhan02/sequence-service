package database

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/rohanchauhan02/sequence-service/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type (
	Postgress interface {
		InitClient(ctx context.Context) (*gorm.DB, error)
	}

	database struct {
		SharedConfig config.ImmutableConfig
	}
)

var (
	once sync.Once
	db   *gorm.DB
	err  error
)

var gormOpen = gorm.Open

func NewPostgressClient(conf config.ImmutableConfig) Postgress {
	return &database{
		SharedConfig: conf,
	}
}

func (d *database) InitClient(ctx context.Context) (*gorm.DB, error) {
	once.Do(func() {
		log.Info("Initializing PostgreSQL connection...")

		dbConfig := d.SharedConfig.GetDBConf()

		connectionString := fmt.Sprintf(
			"user=%s password=%s dbname=%s host=%s sslmode=%s",
			dbConfig.User,
			dbConfig.Password,
			dbConfig.Name,
			dbConfig.Host,
			dbConfig.SSLMode,
		)

		// Retry mechanism for transient failures
		maxRetries := 3
		for i := range maxRetries {
			db, err = gormOpen(postgres.New(postgres.Config{
				DSN:                  connectionString,
				PreferSimpleProtocol: true,
			}), &gorm.Config{
				DisableAutomaticPing: false,
				PrepareStmt:          true,
			})

			if err == nil {
				break
			}

			log.Errorf("Postgres connection attempt %d failed: %v", i+1, err)
			time.Sleep(2 * time.Second)
		}

		if err != nil {
			log.Errorf("PostgreSQL connection failed after retries: %v", err)
			return
		}

		sqlDB, sqlErr := db.DB()
		if sqlErr != nil {
			log.Errorf("Failed to get underlying sql.DB: %v", sqlErr)
			err = sqlErr
			return
		}

		// Set secure and optimized connection pool settings
		sqlDB.SetMaxIdleConns(10)
		sqlDB.SetMaxOpenConns(50)
		// Reduce risk of stale connections
		sqlDB.SetConnMaxLifetime(30 * time.Minute)

		log.Info("🚀 Successfully connected to PostgreSQL!")
	})

	return db, err
}
