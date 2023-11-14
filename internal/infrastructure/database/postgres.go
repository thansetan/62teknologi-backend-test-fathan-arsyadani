package database

import (
	"fmt"

	"github.com/thansetan/62teknologi-backend-test-fathan-arsyadani/internal/domain"
	"github.com/thansetan/62teknologi-backend-test-fathan-arsyadani/internal/utils/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func New(conf config.DB) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=Asia/Jakarta", conf.Host, conf.User, conf.Password, conf.Name, conf.Port)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(domain.Business{}, domain.Category{}, domain.Transaction{}, domain.Location{})
	if err != nil {
		return nil, err
	}

	// populate the category table
	var count int64
	db.Model(&domain.Category{}).Count(&count)
	if count == 0 {
		err = populateTable(db, categories)
		if err != nil {
			return nil, err
		}
	}

	// populate the transaction table
	db.Model(&domain.Transaction{}).Count(&count)
	if count == 0 {
		err = populateTable(db, transactions)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}
