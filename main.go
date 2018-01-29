package main

import (
	"database/sql"
	"github.com/ipfans/echo-session"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
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
	e.Use(middleware.CORS())

	store, err := session.NewRedisStore(32, "tcp", "localhost:6379", "", []byte("secret"))
	if err != nil {
		log.Fatal("redis error")
	}
	e.Use(session.Sessions("GSESSION", store))

	e.POST("/users", func(c echo.Context) error {
		// FIMXE:クソ実装しています
		u := new(user)
		if err := c.Bind(u); err != nil {
			log.Print(err)
		}
		log.Println("username: " + u.ID)
		log.Println("password: " + u.Password)
		db.Query("insert into shuho_user values('" + u.ID + "', '" + u.Password + "')")
		return c.String(http.StatusOK, "Add user!!")
	})

	e.GET("/users", func(c echo.Context) error {
		rows, err := db.Query("select * from shuho_user;")
		if err != nil {
			log.Fatal("db select error")
		}

		arrayUsers := []map[string]string{}
		for rows.Next() {
			var id string
			var password string
			err := rows.Scan(&id, &password)
			if err != nil {
				log.Fatal("scan error: %v", err)
			}
			jsonMap := map[string]string{
				"id":    id,
				"password": password,
			}
			arrayUsers = append(arrayUsers, jsonMap)
		}
		return c.JSON(http.StatusOK, arrayUsers)
	})

	e.POST("/session", func(c echo.Context) error {
		// TODO: 実装するよ
		return nil
	})

	log.Fatal(e.Start(":" + port))
}
