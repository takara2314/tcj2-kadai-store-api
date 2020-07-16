package main

import (
	"fmt"
	"net"
	"net/http/fcgi"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	var err error

	// 毎時指定した時間に課題一覧を取得
	go getRegularly([]int{0, 5, 10, 15, 20, 25, 30, 35, 40, 45, 50, 55})

	fmt.Println("ルーティン稼働できたよー！")

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

// homeRequestFunc は/アクセスされたときの処理
func homeRequestFunc(c *gin.Context) {
	c.String(404, "Please use get or status parameters.")
}

// getRequestFunc は/getアクセスされたときの処理
func getRequestFunc(c *gin.Context) {
	// ヘッダー「Authorization」を取得
	authHeader := c.Request.Header.Get("Authorization")
	// Bearerトークンであり、許可されたトークンであればJSONを返す (そうでなければ401)
	if strings.HasPrefix(authHeader, "Bearer ") && tokenCheck(strings.TrimLeft(authHeader, "Bearer ")) {
		c.JSON(200, homeworksData)
	} else {
		c.String(401, "401 Unauthorized")
	}
}

// tokenCheck はAuthorization(Header)のトークンと一致すればtrueを返す関数
func tokenCheck(hToken string) bool {
	for _, allowedToken := range allowedTokens {
		if hToken == allowedToken {
			return true
		}
	}
	return false
}
