package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"net/http"
)

type Person struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

var db *gorm.DB

func init() {
	var err error
	db, err = gorm.Open(sqlite.Open("test.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Включение отладочного режима
	})
	if err != nil {
		log.Fatal("Не удалось подключиться к базе данных:", err)
	}

	// Обработка ошибок при миграции
	if err := db.AutoMigrate(&Person{}); err != nil {
		log.Fatal("Не удалось выполнить миграцию:", err)
	}
}

func main() {
	r := gin.Default()
	r.GET("/people", getPeople)
	r.POST("/person", createPerson)
	r.PUT("/person/:id", updatePerson)
	r.DELETE("/person/:id", deletePerson)

	if err := r.Run(":8080"); err != nil {
		log.Fatal("Ошибка запуска сервера:", err)
	}
}

func createPerson(c *gin.Context) {
	var person Person
	if err := c.ShouldBindJSON(&person); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Create(&person)
	c.JSON(http.StatusOK, person)
}

func getPeople(c *gin.Context) {
	var people []Person
	db.Find(&people)
	c.JSON(http.StatusOK, people)
}

func updatePerson(c *gin.Context) {
	var person Person
	id := c.Param("id")
	if err := db.First(&person, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Человек не найден"})
		return
	}
	if err := c.ShouldBindJSON(&person); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.Save(&person)
	c.JSON(http.StatusOK, person)
}

func deletePerson(c *gin.Context) {
	var person Person
	id := c.Param("id")
	if err := db.First(&person, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Человек не найден"})
		return
	}
	db.Delete(&person)
	c.JSON(http.StatusOK, gin.H{"message": "Человек удален"})
}
