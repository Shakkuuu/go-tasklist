package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.LoadHTMLGlob("views/*.html")

	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(200, "index.html", nil)
	})

	router.GET("/aaa", func(ctx *gin.Context) {
		ctx.HTML(200, "result.html", nil)
	})

	router.POST("/result", func(ctx *gin.Context) {
		bbb := ctx.PostForm("name")
		ctx.HTML(200, "index.html", gin.H{"bbb": bbb})
	})

	router.Run()
}
