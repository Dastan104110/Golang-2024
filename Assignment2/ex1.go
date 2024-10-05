package main

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

func createTable(db *sql.DB) {
	createTableSQL := `
    CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        name TEXT NOT NULL,
        age INT NOT NULL
    );`
	_, err := db.Exec(createTableSQL)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}
}

func insertUser(db *sql.DB, name string, age int) {
	insertSQL := `INSERT INTO users (name, age) VALUES ($1, $2)`
	_, err := db.Exec(insertSQL, name, age)
	if err != nil {
		log.Fatalf("Error inserting data: %v", err)
	}
}

func queryUsers(db *sql.DB) {
	rows, err := db.Query("SELECT id, name, age FROM users")
	if err != nil {
		log.Fatalf("Error querying data: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		var age int
		err := rows.Scan(&id, &name, &age)
		if err != nil {
			log.Fatalf("Error reading data: %v", err)
		}
		log.Printf("%d %s %d\n", id, name, age)
	}
}

func main() {
	connStr := "user=dastan dbname=goproject password=123104110115118 host=localhost sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	createTable(db)
	insertUser(db, "Dastan", 21)
	insertUser(db, "Erlan", 23)
	insertUser(db, "Aspandiyar", 20)
	queryUsers(db)
}
