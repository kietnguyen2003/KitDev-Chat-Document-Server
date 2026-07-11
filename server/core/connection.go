package core

import (
	"fmt"
	"server/internal/infracstructure"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB(url string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		fmt.Println("Connect fail with err: ", err.Error())
		return nil, fmt.Errorf("Fail to connect")
	}
	err = AutoMigrate(db)
	if err != nil {
		fmt.Println("Fail to Migration", err)
		return nil, fmt.Errorf("Fail to migration")
	}
	return db, nil
}

func AutoMigrate(db *gorm.DB) error {
	if err := db.AutoMigrate(
		&infracstructure.GormUser{},
		&infracstructure.GormCategory{},
		&infracstructure.GormStorage{},
		&infracstructure.GormDocument{},
	); err != nil {
		return err
	}
	return nil
}
