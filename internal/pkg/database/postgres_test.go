package database

import (
	"context"
	"sync"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	mock_config "github.com/rohanchauhan02/sequence-service/files/mocks/config"
	"github.com/rohanchauhan02/sequence-service/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func resetSingleton() {
	db = nil
	err = nil
	once = sync.Once{}
}

func Test_InitClient(t *testing.T) {
	resetSingleton()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockConf := mock_config.NewMockImmutableConfig(ctrl)
	mockConf.EXPECT().GetDBConf().Return(config.DB{
		Host:     "localhost",
		Port:     "5432",
		User:     "postgres",
		Password: "password",
		Name:     "sequence_db",
		SSLMode:  "disable",
	}).Times(1)

	sqlDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}
	defer sqlDB.Close()

	origGormOpen := gormOpen
	gormOpen = func(dialector gorm.Dialector, opts ...gorm.Option) (*gorm.DB, error) {
		db, _ := gorm.Open(postgres.New(postgres.Config{
			Conn: sqlDB,
		}), &gorm.Config{})
		return db, nil
	}
	defer func() { gormOpen = origGormOpen }()

	client := NewPostgressClient(mockConf)
	ctx := context.Background()

	dbConn, err := client.InitClient(ctx)
	if err != nil {
		t.Errorf("InitClient() unexpected error: %v", err)
	}
	if dbConn == nil {
		t.Error("InitClient() returned nil db")
	}
}
