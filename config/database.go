package config

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const (
	host     = "localhost"
	port     = "5432"
	user     = "admin"
	password = "admin"
	dbName   = "mydb"
)

func DatabaseConnection() *gorm.DB {
	// 1. Format Data Source Name (DSN) khusus untuk PostgreSQL
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		host, user, password, dbName, port,
	)

	// 2. Membuka koneksi ke database menggunakan GORM
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		// Menghentikan aplikasi (fatal) jika database tidak bisa diakses
		log.Fatalf("Gagal terhubung ke database PostgreSQL: %v", err)
	}

	fmt.Println("Koneksi ke database PostgreSQL berhasil!")
	return db
}
