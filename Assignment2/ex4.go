package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

type User struct {
	ID   int
	Name string
	Age  int
}

func main() {
	connStr := "host=localhost port=5432 user=dastan password=123104110115118 dbname=goproject sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Создание таблицы
	CreateTable(db)

	// Вставка пользователей
	insertUsers(db)

	// Запрос пользователей без фильтрации и с пагинацией
	queryUsers(db, 0, 1, 10) // Все пользователи, 1 страница, 10 пользователей на странице

	// Обновление пользователя
	updateUser(db, 1, "Dastan Updated", 22)

	// Запрос пользователей после обновления
	queryUsers(db, 0, 1, 10) // Все пользователи, 1 страница, 10 пользователей на странице

	// Удаление пользователя
	deleteUser(db, 2)

	// Запрос пользователей после удаления
	queryUsers(db, 0, 1, 10) // Все пользователи, 1 страница, 10 пользователей на странице

	// Запрос пользователей с фильтрацией по возрасту
	queryUsers(db, 21, 1, 10) // Фильтр по возрасту, 1 страница, 10 пользователей на странице
}

// Функция создания таблицы
func CreateTable(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		name VARCHAR(100) UNIQUE NOT NULL,
		age INT NOT NULL
	);`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Ошибка создания таблицы: %v", err)
	} else {
		fmt.Println("Таблица users успешно создана.")
	}
}

// Функция вставки пользователей
func insertUsers(db *sql.DB) {
	users := []User{
		{Name: "Dastan", Age: 21},
		{Name: "Erlan", Age: 23},
		{Name: "Aspandiyar", Age: 20},
	}

	for _, user := range users {
		_, err := db.Exec("INSERT INTO users (name, age) VALUES ($1, $2)", user.Name, user.Age)
		if err != nil {
			log.Printf("Ошибка вставки пользователя %s: %v", user.Name, err)
		}
	}
	fmt.Println("Пользователи успешно вставлены.")
}

// Функция запроса пользователей с фильтрацией и пагинацией
func queryUsers(db *sql.DB, ageFilter int, page int, pageSize int) {
	var query string
	var args []interface{}
	if ageFilter > 0 {
		query = `SELECT id, name, age FROM users WHERE age >= $1 LIMIT $2 OFFSET $3;`
		args = append(args, ageFilter, pageSize, (page-1)*pageSize)
	} else {
		query = `SELECT id, name, age FROM users LIMIT $1 OFFSET $2;`
		args = append(args, pageSize, (page-1)*pageSize)
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		log.Fatalf("Ошибка выборки пользователей: %v", err)
	}
	defer rows.Close()

	fmt.Println("Пользователи:")
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Name, &user.Age); err != nil {
			log.Fatalf("Ошибка сканирования строки: %v", err)
		}
		fmt.Printf("ID: %d, Name: %s, Age: %d\n", user.ID, user.Name, user.Age)
	}
}

// Функция обновления пользователя
func updateUser(db *sql.DB, id int, name string, age int) {
	_, err := db.Exec("UPDATE users SET name = $1, age = $2 WHERE id = $3", name, age, id)
	if err != nil {
		log.Fatalf("Ошибка обновления пользователя: %v", err)
	}
	fmt.Println("Пользователь успешно обновлен.")
}

// Функция удаления пользователя
func deleteUser(db *sql.DB, id int) {
	_, err := db.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		log.Fatalf("Ошибка удаления пользователя: %v", err)
	}
	fmt.Println("Пользователь успешно удален.")
}
