package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"

	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Todo struct {
	id         string
	date       string
	text       string
	identifier string
}

func main() {
	var (
		id         int
		text       string
		identifier string
		date       []uint8
	)

	router := gin.Default()

	db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/react_todo")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "DELETE"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	}))

	api := router.Group("/api")
	{
		api.GET("/getTodos", func(ctx *gin.Context) {
			var ip = ctx.ClientIP()

			res, err := db.Query("SELECT `id`, `text`, `identifier`, `date` FROM `todo_items` WHERE `identifier` = '" + ip + "'")

			if err != nil {
				log.Fatal(err)
			}

			defer res.Close()

			if err != nil {
				log.Fatal(err)
				ctx.JSON(500, gin.H{"error": err})
				return
			} else {
				rows := make([]map[string]string, 0)

				for res.Next() {
					err := res.Scan(&id, &text, &identifier, &date)

					if err != nil {
						log.Fatal(err)
						ctx.JSON(500, gin.H{"error": err})
					} else {
						rows = append(rows, map[string]string{
							"id":         strconv.Itoa(id),
							"text":       text,
							"date":       strconv.Itoa(int(date[1])),
							"identifier": identifier,
						})
					}
				}

				ctx.JSON(200, gin.H{"todos": rows})
			}
		})

		api.POST("/addTodo", func(ctx *gin.Context) {
			var params = ctx.Request.URL.Query()

			if params != nil && params["text"] != nil {
				var text = params["text"][0]
				var ip = ctx.ClientIP()

				db.ExecContext(ctx, "INSERT INTO `todo_items` (`text`, `date`, `identifier`) VALUES ('"+text+"', CURRENT_TIMESTAMP(), '"+ip+"')")

				ctx.JSON(200, gin.H{})
			} else {
				ctx.JSON(400, gin.H{"msg": "Invalid parameters"})
			}
		})

		api.DELETE("/removeTodo", func(ctx *gin.Context) {
			var params = ctx.Request.URL.Query()

			if params != nil && params["id"] != nil {
				var id = params["id"][0]
				var ip = ctx.ClientIP()

				db.ExecContext(ctx, "DELETE FROM `todo_items` WHERE `id` = "+id+" AND `identifier` = '"+ip+"'")

				ctx.JSON(200, gin.H{})
			} else {
				ctx.JSON(400, gin.H{"msg": "Invalid parameters"})
			}
		})
	}

	router.NoRoute(func(ctx *gin.Context) { ctx.JSON(http.StatusNotFound, gin.H{}) })

	router.Run(":8080")
}
