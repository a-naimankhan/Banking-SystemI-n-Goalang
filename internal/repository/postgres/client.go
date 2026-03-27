package postgres

import (
	"BankingSystem/internal/domain"
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewClient() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s password=%s dbname=%s port=%s user=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_SSLMODE"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&domain.Account{}, &domain.SavingAccount{})
	return db, err
}
