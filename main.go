package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"html"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	dbURL := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("failed to connect to db")
	}
	defer db.Close()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q\n", html.EscapeString(r.URL.Path))

		rows, err := db.Query("select * from test;")
		if err != nil {
			log.Fatal("db select error")
		}

		for rows.Next() {
			var id int
			var value int
			err := rows.Scan(&id, &value)
			if err != nil {
				log.Fatal("scan error: %v", err)
			}
			fmt.Fprintf(w, "id: %d, value: %d\n", id, value)
		}
	})

	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		query := `select * from "user"`
		rows, err := db.Query(query)
		if err != nil {
			log.Fatalf("unable to select from db. err: %v", err);
		}

		for rows.Next() {
			var id int
			var name string
			err := rows.Scan(&id, &name)
			if err != nil {
				log.Fatal("scan error: %v", err)
			}
			fmt.Fprintf(w, "id: %d, name: %s\n", id, name)
		}
	})

	log.Fatal(http.ListenAndServe(":"+port, nil))
}
