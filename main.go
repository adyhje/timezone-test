package main

import (
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"time"
)

type Database struct {
	*sql.DB
}

var DB Database

func Connect() {
	connectionString := "root:123456@tcp(localhost:3306)/testdb?charset=utf8&parseTime=True&loc=Europe%2FBerlin"
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		panic(err.Error())
	}
	if err = db.Ping(); err != nil {
		panic("could not ping database")
	}
	log.Println("Successfully connected to database")
	DB = Database{db}
}

type TimeEntry struct {
	Timestamp time.Time `json:"timestamp"`
}

func getTimesHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := DB.Query("SELECT zeitstempel FROM time_example")
	if err != nil {
		http.Error(w, "Query failed", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var entries []TimeEntry
	for rows.Next() {
		var t time.Time
		if err := rows.Scan(&t); err != nil {
			http.Error(w, "Scan failed", http.StatusInternalServerError)
			return
		}
		entries = append(entries, TimeEntry{Timestamp: t.UTC()})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(entries)
}

func postTimeHandler(w http.ResponseWriter, r *http.Request) {
	var entry TimeEntry
	if err := json.NewDecoder(r.Body).Decode(&entry); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	_, err := DB.Exec("INSERT INTO time_example (zeitstempel) VALUES (?)", entry.Timestamp)
	if err != nil {
		http.Error(w, "Insert failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func setupRoutes() {
	http.HandleFunc("/times", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getTimesHandler(w, r)
		} else if r.Method == http.MethodPost {
			postTimeHandler(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func main() {
	Connect()

	_, _ = DB.Exec("CREATE TABLE IF NOT EXISTS time_example ( zeitstempel DATETIME );")

	setupRoutes()

	log.Println("Server running on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
