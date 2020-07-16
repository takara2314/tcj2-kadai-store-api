package main

import (
	"io/ioutil"
	"log"
	"strings"
	"time"
)

var (
	// 全ての課題情報が格納されているインスタンス (タイムゾーン: JST)
	homeworksData ResponseJSON
	// 許可されたトークンリスト
	allowedTokens []string
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

func init() {
	// トークンリストは前のディレクトリの中のtokenファイルに書いてある
	fileData, err := ioutil.ReadFile("../tcj2-kadai-store-api.token")
	if err != nil {
		log.Fatal(err)
	}
	allowedTokens = strings.Split(string(fileData), "\n")
}
