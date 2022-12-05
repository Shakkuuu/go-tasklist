package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.StaticFS("/views", http.Dir("views"))

	router.GET("/", func(ctx *gin.Context) {
		ctx.Redirect(302, "/views/index.html")
	})

	router.GET("/hello", func(ctx *gin.Context) {
		name := ctx.Query("name")
		ctx.Header("Content-Type", "text/html; charset=UTF-8")
		ctx.String(200, "<h1>Hello, "+name+"</h1>")
	})

	err := router.Run("127.0.0.1:8888")
	if err != nil {
		log.Fatal("サーバー起動に失敗 ", err)
	}

}
