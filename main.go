package main

import (
	"database/sql"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
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

	e := echo.New()

	e.GET("/hoge", func(c echo.Context) error {
		return c.String(http.StatusOK, "hoge")
	})

	e.GET("/", func(c echo.Context) error {
		c.String(http.StatusOK, "Hello")

		rows, err := db.Query("select * from shuho_user;")
		if err != nil {
			log.Fatal("db select error")
		}

		for rows.Next() {
			var id string
			var value string
			err := rows.Scan(&id, &value)
			if err != nil {
				log.Fatal("scan error: %v", err)
			}
			jsonMap := map[string]string{
				"id":    id,
				"value": value,
			}
			return c.JSON(http.StatusOK, jsonMap)
		}
		return nil
	})

	log.Fatal(e.Start(":" + port))
}
