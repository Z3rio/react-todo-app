package main

import (
  "database/sql"
  _ "github.com/go-sql-driver/mysql"

  "strconv"
  "log"
  "net/http"
  
  "github.com/gin-gonic/gin"
)

func ToStringint(x int) string { 
  return strconv.Itoa(x)
}

func main() {
  router := gin.Default()

  db, err := sql.Open("mysql", "root@tcp(127.0.0.1:3306)/react_todo")
  defer db.Close()

  if err != nil {
      log.Fatal(err)
  }

  api := router.Group("/api")
  {
    api.GET("/getTodo", func(ctx *gin.Context) {
      var ip = ctx.ClientIP()

      log.Print(ip)

      res, err := db.Query("SELECT * FROM `todo_items` WHERE `ip` = '" + ip + "'")
      defer res.Close()

      if err != nil {
        log.Fatal(err)
      }

      ctx.JSON(200, gin.H{"response": res})
    })

    api.POST("/addTodo", func(ctx *gin.Context) {
      var params = ctx.Request.URL.Query()

      if params != nil && params["text"] != nil {
        var text = params["text"][0]

        db.ExecContext(ctx, "INSERT INTO `todo_items` (`text`, `date`) VALUES ('" + text + "', " + ToStringint(0) + ")")
        ctx.JSON(200, gin.H{})
      } else {
        ctx.JSON(400, gin.H{"msg": "Invalid parameters"})
      }
    })
  }

  router.NoRoute(func(ctx *gin.Context) { ctx.JSON(http.StatusNotFound, gin.H{}) })

  router.Run(":8080")
}
