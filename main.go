package main

import (
	"database/sql"
	"github.com/labstack/echo"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

type (
	user struct {
		ID string `json:"id"`
		Password string `json:"password"`
	}
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

	e.POST("/users", func(c echo.Context) error {
		// FIXME:クソ実装しています
		u := new(user)
		if err := c.Bind(u); err != nil {
			log.Print(err)
		}
		log.Println("username: " + u.ID)
		log.Println("password: " + u.Password)
		db.Query("insert into shuho_user values('" + u.ID + "', '" + u.Password + "')")
		return c.String(http.StatusOK, "Add user!!")
	})

	e.GET("/", func(c echo.Context) error {
		c.String(http.StatusOK, "Hello")

		rows, err := db.Query("select * from shuho_user;")
		if err != nil {
			log.Fatal("db select error")
		}

		arrayUsers := []map[string]string{}
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
			arrayUsers = append(arrayUsers, jsonMap)
		}
		return c.JSON(http.StatusOK, arrayUsers)
	})

	log.Fatal(e.Start(":" + port))
}
