package main

import (
	"net"
	"net/http/fcgi"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	var err error

	// 毎時指定した時間に課題一覧を取得
	go getRegularly([]int{0, 2, 4, 6, 8, 10, 12, 14, 16, 18, 20, 22, 24, 26, 28, 30, 32, 34, 36, 38, 40, 42, 44, 46, 48, 50, 52, 54, 56, 58})

	r := gin.Default()

	r.GET("/", homeRequestFunc)
	r.GET("/get", getRequestFunc)
	r.GET("/version", versionRequestFunc)

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
	c.String(404, "Please add get or status path.")
}

// getRequestFunc は/getアクセスされたときの処理
func getRequestFunc(c *gin.Context) {
	// ヘッダー「Authorization」を取得
	authHeader := c.Request.Header.Get("Authorization")
	// Bearerトークンであり、許可されたトークンであればJSONを返す (そうでなければ401)
	if strings.HasPrefix(authHeader, "Bearer ") && tokenCheck(strings.TrimLeft(authHeader, "Bearer ")) {
		// URL変数due-targetに提出期限の指定を入れることで、返される課題一覧を調整
		if c.Query("due") == "future" {
			// タイムゾーンの指定
			if c.Query("timezone") == "Asia/Tokyo" {
				c.JSON(200, homeworksDataOnlyFutureJST)
			} else {
				c.JSON(200, homeworksDataOnlyFuture)
			}
		} else {
			// タイムゾーンの指定
			if c.Query("timezone") == "Asia/Tokyo" {
				c.JSON(200, homeworksDataJST)
			} else {
				c.JSON(200, homeworksData)
			}
		}
	} else {
		c.String(401, "401 Unauthorized")
	}
}

// versionRequestFunc は/versionアクセスされたときの処理
func versionRequestFunc(c *gin.Context) {
	c.String(200, "TCJ2 Kadai Store API - v0.1.1")
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
