package main

import (
	"net"
	"net/http/fcgi"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	// 全ての課題情報が格納されているインスタンス
	homeworksData ResponseJSON
)

// ResponseJSON は返すJSONの元の構造体
type ResponseJSON struct {
	Homeworks []HomeworkStruct `json:"homeworks"`
}

// HomeworkStruct は1つの課題情報を収納する構造体
type HomeworkStruct struct {
	Subject string    `json:"subject"`
	Omitted string    `json:"omitted"`
	Name    string    `json:"name"`
	ID      string    `json:"id"`
	Due     time.Time `json:"due"`
}

func main() {
	var err error

	// 毎時指定した時間に課題一覧を取得
	go getRegularly([]int{0, 10, 20, 30, 40, 50, 55})

	r := gin.Default()

	r.GET("/", homeRequestFunc)
	r.GET("/get", getRequestFunc)

	l, err := net.Listen("tcp", ":2314")
	if err != nil {
		panic(err)
	}
	if err := fcgi.Serve(l, r); err != nil {
		panic(err)
	}
}

func homeRequestFunc(c *gin.Context) {
	c.String(404, "Please use get or status parameters.")
}

func getRequestFunc(c *gin.Context) {
	if c.ClientIP() == "126.93.172.178" {
		c.JSON(200, homeworksData)
	} else {
		c.String(200, "誰だおめぇ")
	}
}

// func getRequestFuncO(c *gin.Context) {
// 	userBrowser := c.Request.Header.Get("User-Agent")
// 	fmt.Println(userBrowser)

// 	returnJSON := scheduleList{
// 		Date: time.Now(),
// 		Homeworks: []HomeworkStructs{
// 			{
// 				Comment: "プリントをやってください",
// 				Due:     time.Now(),
// 				Subject: "基礎数学４",
// 			},
// 			{
// 				Comment: "ここしてください",
// 				Due:     time.Now(),
// 				Subject: "電気電子基礎",
// 			},
// 			{
// 				Comment: "感想を800字以内に書いてください",
// 				Due:     time.Now(),
// 				Subject: "一般基礎教育２",
// 			},
// 			{
// 				Comment: "Do the handout.",
// 				Due:     time.Now(),
// 				Subject: "Sonzaishinai English 2",
// 			},
// 		},
// 	}
// 	c.JSON(200, returnJSON)
// }
