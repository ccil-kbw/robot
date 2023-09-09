package main

import (
	"context"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type VideoUpload struct {
	gorm.Model
	FileName string
	Success  bool
}

type Database struct {
	client *gorm.DB
}

func setupDB() *Database {
	db, err := gorm.Open(sqlite.Open("kbw.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&VideoUpload{})

	return &Database{
		client: db,
	}
}

func (db *Database) insertVideo(ctx context.Context, filename string) {
	db.client.Create(&VideoUpload{FileName: filename, Success: false})
}

func (db *Database) updateVideo(ctx context.Context, filename string, success bool) {
	db.client.Model(&VideoUpload{}).Where("file_name = ?", filename).Update("success", success)
}
