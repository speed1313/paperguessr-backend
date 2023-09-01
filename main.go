// echo server
package main

import (
	"fmt"
	"net/http"
	// echo
	"github.com/labstack/echo"

	"database/sql"
	// sqlite
	_ "github.com/mattn/go-sqlite3"
	// json
)

func main() {
	// init database
	db, err := sql.Open("sqlite3", "./test.db")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Close()

	// create table
	// point 0 to 5 number
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS SCORE (id INTEGER PRIMARY KEY AUTOINCREMENT, score INTEGER check(score >= 0 and score <= 5))")
	if err != nil {
		fmt.Println(err)
	}

	// create a new echo server

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!")
	})

	// fetch point info between 0 to 5 points
	e.GET("/points/:point", func(c echo.Context) error {
		point := c.Param("point")
		if point < "0" || point > "5" {
			return c.String(http.StatusBadRequest, "Invalid point")
		}
		// register score to db
		_, err = db.Exec("INSERT INTO SCORE (score) VALUES (?)", point)
		if err != nil {
			fmt.Println(err)
		}
		return c.String(http.StatusOK, point)
	})

	// serve point per user info
	e.GET("/points", func(c echo.Context) error {
		// get score histogram
		rows, err := db.Query("SELECT score, count(*) FROM SCORE GROUP BY score")
		if err != nil {
			fmt.Println(err)
		}
		defer rows.Close()
		// create response json
		type Score struct {
			Point string `json:"point"`
			Count int    `json:"count"`
		}
		var scores []Score
		for rows.Next() {
			var score Score
			err = rows.Scan(&score.Point, &score.Count)
			if err != nil {
				fmt.Println(err)
			}
			scores = append(scores, score)
		}
		// return json
		return c.JSON(http.StatusOK, scores)

	})

	// start server
	e.Logger.Fatal(e.Start(":1323"))

}
