package main

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type User struct {
	ID      uint `gorm:"primaryKey"`
	Name    string
	Age     int
	Profile Profile `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;foreignKey:UserID"`
}

type Profile struct {
	ID                uint `gorm:"primaryKey"`
	UserID            uint `gorm:"index;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Bio               string
	ProfilePictureURL string
}

func main() {
	dsn := "host=localhost user=dastan password=123104110115118 dbname=goproject port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal("Failed to get the generic database object:", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(0)

	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatal("Failed to migrate User table:", err)
	}

	err = db.AutoMigrate(&Profile{})
	if err != nil {
		log.Fatal("Failed to migrate Profile table:", err)
	}

	createUsersWithProfiles(db)

	log.Println("Users and profiles before deletion:")
	getUsersWithProfiles(db)

	deleteUserWithProfile(db, 1)

	log.Println("Users and profiles after deletion:")
	getUsersWithProfiles(db)
}

func createUsersWithProfiles(db *gorm.DB) {
	err := db.Transaction(func(tx *gorm.DB) error {
		users := []User{
			{
				Name: "Dastan",
				Age:  21,
				Profile: Profile{
					Bio:               "Software developer",
					ProfilePictureURL: "https://example.com/john.jpg",
				},
			},
			{
				Name: "Erlan",
				Age:  23,
				Profile: Profile{
					Bio:               "Graphic designer",
					ProfilePictureURL: "https://example.com/jane.jpg",
				},
			},
		}
		if err := tx.Create(&users).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		log.Println("Failed to create users with profiles:", err)
	} else {
		log.Println("Users with profiles created successfully")
	}
}

func getUsersWithProfiles(db *gorm.DB) {
	var users []User
	db.Preload("Profile").Find(&users)
	for _, user := range users {
		log.Printf("User: %s, Age: %d, Profile Bio: %s, Profile Picture: %s\n",
			user.Name, user.Age, user.Profile.Bio, user.Profile.ProfilePictureURL)
	}
}

func updateProfile(db *gorm.DB, userID uint, newBio string) {
	var profile Profile
	db.First(&profile, "user_id = ?", userID)
	profile.Bio = newBio
	db.Save(&profile)
	log.Printf("Profile updated for user ID %d\n", userID)
}

func deleteUserWithProfile(db *gorm.DB, userID uint) {
	db.Delete(&User{}, userID)
	log.Printf("User with ID %d and associated profile deleted\n", userID)
}

//
