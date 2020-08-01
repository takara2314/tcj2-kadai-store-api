package main

import (
	"net"
	"net/http/fcgi"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	// 毎時指定した時間(分)にdevoirsから課題一覧を取得
	go getRegularly(configData.UpdateTimes)

	// GETリミットが無制限でなければ、リクエストマネージャーを呼び出す
	if configData.GETLimit != -1 {
		go requestManager()
	}

	r := gin.Default()

	r.GET("/", homeRequestFunc)
	r.GET("/get", getRequestFunc)
	r.GET("/version", versionRequestFunc)

	// FastCGIとして動かす場合
	if configData.FCGI {
		l, err := net.Listen("tcp", ":"+strconv.Itoa(configData.Port))
		if err != nil {
			panic(err)
		}
		if err := fcgi.Serve(l, r); err != nil {
			panic(err)
		}
	} else {
		r.Run(":" + strconv.Itoa(configData.Port))
	}
}

// homeRequestFunc は/アクセスされたときの処理
func homeRequestFunc(c *gin.Context) {
	c.String(404, "Please add get or version path.")
}

// getRequestFunc は/getアクセスされたときの処理
func getRequestFunc(c *gin.Context) {
	// ヘッダーAuthorizationを取得
	authHeader := c.Request.Header.Get("Authorization")
	// HTTPレスポンスステータスコードを取得
	var statusCode int = isProvide(authHeader)

	switch statusCode {
	case 200:
		// 正常
		// パラメータdueに提出期限の指定を入れることで、返される課題一覧を調整
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

	case 401:
		// 認証エラー
		c.String(401, "401 Unauthorized")

	case 429:
		// リクエストが多すぎてAPI制限にかかっているならば
		c.String(429, "429 Too Many Requests")
	}
}

// versionRequestFunc は/versionアクセスされたときの処理
func versionRequestFunc(c *gin.Context) {
	c.String(200, "TCJ2 Kadai Store API - v0.2.0 pre1")
}

// isProvide はこのヘッダーをもとにトークンを取得し、HTTPレスポンスステータスコードを返す
func isProvide(authHeader string) int {
	var token string

	// Bearerトークンであるかないか
	if strings.HasPrefix(authHeader, "Bearer ") {
		token = strings.TrimLeft(authHeader, "Bearer ")
	} else {
		return 401
	}

	// 許可されたトークンでなければ
	if !tokenCheck(token) {
		return 401
	}

	// GETリミットが無制限でなければ、提供するかどうかを調べる
	if configData.GETLimit != -1 {
		if tokenLimit[token] == 0 {
			return 429
		}
		// 利用可能回数をデクリメント
		tokenLimit[token]--
	}

	// 何も異常がなければ200
	return 200
}

// tokenCheck はAuthorization(Header)のトークンと一致すればtrueを返す関数
func tokenCheck(hToken string) bool {
	for _, allowedToken := range tokenData.AllowedTokens {
		if hToken == allowedToken {
			return true
		}
	}
	return false
}
