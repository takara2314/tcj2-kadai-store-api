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
	// 全ての課題情報が格納されているインスタンス (タイムゾーン: Asia/Tokyo)
	homeworksDataJST ResponseJSON
	// 全ての課題情報が格納されているインスタンス (タイムゾーン: Asia/Tokyo)
	homeworksDataOnlyFutureJST ResponseJSON
	// 許可されたトークンリスト
	allowedTokens []string
	// Discord通知をONにするか
	isDiscordAlarm bool
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
	var fileName string
	var err error

	// トークンリストは前のディレクトリの中のtokenファイルに書いてある
	fileName = "../kadai-store-api.token"
	fileData, err = ioutil.ReadFile(fileName)

	if err != nil {
		log.Fatal(err)
	}
	allowedTokens = strings.Split(string(fileData), "\n")

	// API管理者のDiscordIDは前のディレクトリの中のidファイルに書いてある
	fileName = "../kadai-store-api_admin-discord-ID.id"
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		isDiscordAlarm = false
	} else {
		isDiscordAlarm = true
		fileData, err = ioutil.ReadFile(fileName)

		if err != nil {
			log.Fatal(err)
		}
		adminDiscordID = strings.TrimRight(string(fileData), "\n")
	}

	// DiscordBotのトークンは前のディレクトリの中のtokenファイルに書いてある
	fileName = "../kadai-store-api_discord-alarm.token"
	if _, err := os.Stat(fileName); !os.IsNotExist(err) {
		fileData, err = ioutil.ReadFile(fileName)
		if err != nil {
			log.Fatal(err)
		}
		var discordToken string = strings.TrimRight(string(fileData), "\n")

		if isDiscordAlarm {
			// DiscordBot起動
			go discordInit(discordToken)
		} else {
			fmt.Println("DiscordBotの起動に失敗しました: 報告するユーザーのIDが書かれたファイルが必要です。")
			fmt.Println("詳しくはREADME.mdをご確認ください。")
		}
	} else if isDiscordAlarm {
		fmt.Println("DiscordBotの起動に失敗しました: トークンが書かれたファイルが必要です。")
		fmt.Println("詳しくはREADME.mdをご確認ください。")
	}
}

// discordInit はDiscordBotを準備するための関数
func discordInit(dToken string) {
	var err error

	dg, err = discordgo.New("Bot " + dToken)
	if err != nil {
		panic(fmt.Sprint("DiscordBotの起動に失敗しました:", err))
	}

	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		panic(fmt.Sprint("DiscordBotの起動に失敗しました:", err))
	}

	fmt.Println("DiscordBotを起動します…")
	// Discordボットを稼働
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	defer dg.Close()
}
