package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type User struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"not null"`
	Age  int    `gorm:"not null"`
}

func main() {
	dsn := "host=localhost user=dastan dbname=goproject password=123104110115118 port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatalf("Error during auto migration: %v", err)
	}

	db.Create(&User{Name: "Dastan", Age: 21})
	db.Create(&User{Name: "Erlan", Age: 23})
	db.Create(&User{Name: "Aspandiyar", Age: 20})

	var users []User
	db.Find(&users)

	for _, user := range users {
		log.Printf("ID: %d, Name: %s, Age: %d\n", user.ID, user.Name, user.Age)
	}
}
