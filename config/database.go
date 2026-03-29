package config

import (
	"errandify/models"
	"fmt"
	"log"
	"os"

	"golang.org/x/crypto/bcrypt"
	"github.com/joho/godotenv"
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
	// 1. Load file .env
	// Kita taruh warning saja (bukan fatal), karena saat di production (seperti Docker/server), 
	// env biasanya disuntikkan langsung dari sistem, bukan dari file .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file, falling back to system environment variables")
	}

	// 2. Ambil nilai dari environment variables
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// 3. Format DSN
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		host, user, password, dbName, port,
	)

	// 4. Buka Koneksi
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database PostgreSQL: %v", err)
	}

	fmt.Println("connection to database PostgreSQL success!")
	return db
}

// func DatabaseConnection() *gorm.DB {
// 	// 1. Format Data Source Name (DSN) khusus untuk PostgreSQL
// 	dsn := fmt.Sprintf(
// 		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
// 		host, user, password, dbName, port,
// 	)

// 	// 2. Membuka koneksi ke database menggunakan GORM
// 	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
// 	if err != nil {
// 		// Menghentikan aplikasi (fatal) jika database tidak bisa diakses
// 		log.Fatalf("failed to connect to database PostgreSQL: %v", err)
// 	}

// 	fmt.Println("connection to database PostgreSQL success!")
// 	return db
// }

func CreateOwnerAccount (db *gorm.DB) {
	//make owner account with email and password, role is owner
	hashedPasswordBytes, _ := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	owner := models.User{
		Role:     "owner",
		Name:     "Admin",
		Password: string(hashedPasswordBytes),
		Email:    "owner@go.id",
	}

	//check if owner account already exists, if not create it
	if db.Where("email=?", owner.Email).FirstOrCreate(&owner) == nil {
		log.Printf("Owner account created with email: %s", owner.Email)
		db.Create(&owner)
	} else {
		log.Printf("Owner account already exists with email: %s", owner.Email)
	}

}