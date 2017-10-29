package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"html"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q\n", html.EscapeString(r.URL.Path))

		db, err := sql.Open("postgres", "user=testrole password=testrole dbname=test sslmode=disable")
		fmt.Printf("%v\n", db)
		if err != nil {
			fmt.Println("failed to connect to db")
		}
		defer db.Close()

		rows, err := db.Query("select * from test;")
		fmt.Printf("%v\n", rows)
		if err != nil {
			fmt.Println("db select error")
		}

		for rows.Next() {
			var id int
			var value int
			err := rows.Scan(&id, &value)
			if err != nil {
				log.Fatal("scan error: %v", err)
			}
			fmt.Printf("id: %d, value: %d\n", id, value)
			fmt.Fprintf(w, "id: %d, value: %d\n", id, value)
		}
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
