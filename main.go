package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

type Task struct {
	gorm.Model
	Content string `form:"content" binding:"required"`
}

func gormConnect() *gorm.DB {
	DBMS := "mysql"
	USER := "Shakku"
	PASS := "12345678"
	DBNAME := "task"
	// MySQLだと文字コードの問題で"?parseTime=true"を末尾につける必要がある
	CONNECT := USER + ":" + PASS + "@/" + DBNAME + "?parseTime=true"
	db, err := gorm.Open(DBMS, CONNECT)

	if err != nil {
		panic(err.Error())
	}
	return db
}

func dbInit() {
	db := gormConnect()

	// コネクション解放解放
	defer db.Close()
	db.AutoMigrate(&Task{}) //構造体に基づいてテーブルを作成
}

func dbInsert(content string) {
	db := gormConnect()

	defer db.Close()
	// Insert処理
	db.Create(&Task{Content: content})
}

func dbUpdate(id int, taskText string) {
	db := gormConnect()
	var task Task
	db.First(&task, id)
	task.Content = taskText
	db.Save(&task)
	db.Close()
}

// 全件取得
func dbGetAll() []Task {
	db := gormConnect()

	defer db.Close()
	var tasks []Task
	// FindでDB名を指定して取得した後、orderで登録順に並び替え
	db.Order("created_at desc").Find(&tasks)
	return tasks
}

// DB一つ取得
func dbGetOne(id int) Task {
	db := gormConnect()
	var task Task
	db.First(&task, id)
	db.Close()
	return task
}

// DB削除
func dbDelete(id int) {
	db := gormConnect()
	var task Task
	db.First(&task, id)
	db.Delete(&task)
	db.Close()
}

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("views/*.html")

	dbInit()

	//task一覧
	r.GET("/", func(c *gin.Context) {
		tasks := dbGetAll()
		c.HTML(200, "index.html", gin.H{"tasks": tasks})
	})

	r.GET("/aaa", func(c *gin.Context) {
		c.HTML(200, "result.html", nil)
	})

	//task登録
	r.POST("/new", func(c *gin.Context) {
		var form Task

		if err := c.Bind(&form); err != nil {
			tasks := dbGetAll()
			c.HTML(http.StatusBadRequest, "index.html", gin.H{"tasks": tasks, "err": err})
			c.Abort()
		} else {
			content := c.PostForm("content")
			dbInsert(content)
			c.Redirect(302, "/")
		}
	})

	//投稿詳細
	r.GET("/detail/:id", func(c *gin.Context) {
		n := c.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		task := dbGetOne(id)
		c.HTML(200, "detail.html", gin.H{"task": task})
	})

	//更新
	r.POST("/update/:id", func(c *gin.Context) {
		n := c.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("ERROR")
		}
		task := c.PostForm("task")
		dbUpdate(id, task)
		c.Redirect(302, "/")
	})

	//削除確認
	r.GET("/delete_check/:id", func(c *gin.Context) {
		n := c.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("ERROR")
		}
		task := dbGetOne(id)
		c.HTML(200, "delete.html", gin.H{"task": task})
	})

	//削除
	r.POST("/delete/:id", func(c *gin.Context) {
		n := c.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("ERROR")
		}
		dbDelete(id)
		c.Redirect(302, "/")

	})

	r.Run()
}
