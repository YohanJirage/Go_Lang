package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	_ "github.com/lib/pq"
)

type Person struct {
	Name string
	Age  int
}

func main() {
	// Open a connection to the PostgreSQL database
	db, err := sql.Open("postgres", "user=postgres password=root dbname=demo sslmode=disable")
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	defer db.Close()

	// Register handler functions for different routes
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		rootHandler(w, r, db)
	})
	http.HandleFunc("/getAlldata", func(w http.ResponseWriter, r *http.Request) {
		getAllDataHandler(w, r, db)
	})
	http.HandleFunc("/getName", func(w http.ResponseWriter, r *http.Request) {
		getNameHandler(w, r, db)
	})

	// Start the HTTP server on port 8080
	fmt.Println("Server listening on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting HTTP server:", err)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Write response headers
	w.Header().Set("Content-Type", "application/json")
	p := Person{"Yohan", 12}

	// Marshal the Person object into JSON
	jsonBytes, err := json.Marshal(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the JSON response
	w.Write(jsonBytes)
}

func getAllDataHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Write response headers
	w.Header().Set("Content-Type", "application/json")

	// Query the database for all data
	rows, err := db.Query("SELECT name, age FROM people")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Iterate through the rows and populate the data slice
	var data []Person
	for rows.Next() {
		var p Person
		err := rows.Scan(&p.Name, &p.Age)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data = append(data, p)
	}

	// Marshal the data slice into JSON
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the JSON response
	w.Write(jsonBytes)
}

func getNameHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Write response headers
	w.Header().Set("Content-Type", "text/plain")

	// Query the database for the name
	var name string
	err := db.QueryRow("SELECT name FROM people WHERE id = $1", 2).Scan(&name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the name as plain text response
	fmt.Fprintln(w, name)
}
