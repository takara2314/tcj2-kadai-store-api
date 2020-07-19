package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	// 全ての課題情報が格納されているインスタンス (タイムゾーン: UTC)
	homeworksData ResponseJSON
	// 全ての課題情報が格納されているインスタンス (タイムゾーン: UTC)
	homeworksDataOnlyFuture ResponseJSON
	// 全ての課題情報が格納されているインスタンス (タイムゾーン: JST)
	homeworksDataJST ResponseJSON
	// 全ての課題情報が格納されているインスタンス (タイムゾーン: JST)
	homeworksDataOnlyFutureJST ResponseJSON
	// 許可されたトークンリスト
	allowedTokens []string
	// DiscordBotのセッション
	dg *discordgo.Session
	// アラームする人のDiscordID
	adminDiscordID string
	// 重要なアラームを送信したかどうか
	discordAlarmed bool = false
)

// ResponseJSON は返すJSONの元の構造体
type ResponseJSON struct {
	Acquisition time.Time        `json:"acquisition"`
	Homeworks   []HomeworkStruct `json:"homeworks"`
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
	var fileData []byte
	var err error

	// トークンリストは前のディレクトリの中のtokenファイルに書いてある
	fileData, err = ioutil.ReadFile("../tcj2-kadai-store-api.token")
	if err != nil {
		log.Fatal(err)
	}
	allowedTokens = strings.Split(string(fileData), "\n")

	// Discordのトークンリストは前のディレクトリの中のtokenファイルに書いてある
	fileData, err = ioutil.ReadFile("../tcj2-kadai-store-api_admin-discord-ID.token")
	if err != nil {
		log.Fatal(err)
	}
	adminDiscordID = strings.TrimRight(string(fileData), "\n")

	// Discordのトークンリストは前のディレクトリの中のtokenファイルに書いてある
	fileData, err = ioutil.ReadFile("../tcj2-kadai-store-api_discord-alarm.token")
	if err != nil {
		log.Fatal(err)
	}
	var discordToken string = strings.TrimRight(string(fileData), "\n")

	// DiscordBot起動
	go discordInit(discordToken)
}

// discordInit はDiscordBotを準備するための関数
func discordInit(dToken string) {
	var err error

	dg, err = discordgo.New("Bot " + dToken)
	if err != nil {
		panic(fmt.Sprint("Discordセッション作成にエラーが発生しました:", err))
	}

	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		panic(fmt.Sprint("接続エラーが発生しました:", err))
	}

	// Discordボットを稼働
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	defer dg.Close()
}
