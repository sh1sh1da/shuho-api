package main

import (
	"github.com/garyburd/redigo/redis"
	"github.com/soveran/redisurl"
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
	redisURL := os.Getenv("REDISCLOUD_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("failed to connect to db")
	}
	defer db.Close()

	e := echo.New()
	e.Use(middleware.CORS())

	// redisさん起動
	pool := redis.NewPool(func() (redis.Conn, error) {
		return redisurl.ConnectToURL(redisURL)
	}, 30)
	store, err := session.NewRedisStoreWithPool(pool, []byte("secret"))
	if err != nil {
		log.Fatal("redis error")
	}

	// sessionあるか確認する
	e.Use(session.Sessions("GSESSION", store))
	e.GET("/session", func(c echo.Context) error {
		session := session.Default(c)
		v := session.Get("session")
		if v == nil { // not authorized
			c.JSON(200, map[string]interface{}{
				"authorized": false,
			})
		} else {  // authorized
			c.JSON(200, map[string]interface{}{
				"authorized": true,
			})
		}
		return nil
	})

	// session作る
	e.POST("/session", func(c echo.Context) error {
		session := session.Default(c)
		v := session.Get("session")
		if v != nil {
			log.Println("session already exists.")
			return c.JSON(200, map[string]interface{}{
				"authorized": true,
			})
		}

		u := new(user)
		if err := c.Bind(u); err != nil {
			log.Print(err)
		}
		rows, err := db.Query("SELECT password FROM shuho_user WHERE id = '" + u.ID + "'")
		if err != nil {
			log.Fatal("db select error.")
		}
		var password string
		for rows.Next() {
			err = rows.Scan(&password)
			if err != nil {
				log.Fatal("db result scan error.")
			}
		}
		if password != u.Password {
			log.Println("id or password is invalid.")
			return c.JSON(200, map[string]interface{}{
				"authorized": false,
			})
		}
		session.Set("session", true) // FIXME: 何セットしたらいいかわからん
		session.Save()
		return c.JSON(200, map[string]interface{}{
				"authorized": true,
		})
	})

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

	log.Fatal(e.Start(":" + port))
}
