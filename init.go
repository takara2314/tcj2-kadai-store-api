package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"gopkg.in/yaml.v2"
)

var (
	// config.yaml, token.yaml で設定した情報
	configData configYaml
	tokenData  tokenYaml
	// 各トークンの残りのGETリクエスト可能回数を記録
	tokenLimit map[string]int = make(map[string]int, 0)
	// 全ての課題情報 (タイムゾーン: UTC)
	homeworksData ResponseJSON
	// 全ての課題情報 (タイムゾーン: UTC) (期限: 未来のみ)
	homeworksDataOnlyFuture ResponseJSON
	// 全ての課題情報 (タイムゾーン: Asia/Tokyo)
	homeworksDataJST ResponseJSON
	// 全ての課題情報 (タイムゾーン: Asia/Tokyo) (期限: 未来のみ)
	homeworksDataOnlyFutureJST ResponseJSON
	// Discordボットのセッション
	dg *discordgo.Session
	// アラームする人のDiscordID
	adminDiscordID string
	// アラームを送信したかどうか
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

// config.yaml のデータを格納する構造体
type configYaml struct {
	UpdateTimes []int              `yaml:"update-times"`
	GETLimit    int                `yaml:"get-limit"`
	Subjects    configYamlSubjects `yaml:"subjects"`
	Discord     configYamlDiscord  `yaml:"discord"`
}
type configYamlSubjects struct {
	Teams    []string `yaml:"teams"`
	Syllabus []string `yaml:"syllabus"`
	Omitted  []string `yaml:"omitted"`
}
type configYamlDiscord struct {
	Alarm         bool   `yaml:"alarm"`
	AdminID       string `yaml:"admin-id"`
	MessageFormat string `yaml:"message-format"`
	CommandPrefix string `yaml:"command-prefix"`
}

// token.yaml のデータを格納する構造体
type tokenYaml struct {
	AllowedTokens []string `yaml:"allowed-tokens"`
	DiscordToken  string   `yaml:"discord-token"`
}

func init() {
	var fileData []byte
	var err error

	// config.yaml を読み込む
	if !isFileExist("config.yaml") {
		log.Fatalln("エラー: config.yaml が見つかりません")
		log.Fatalln("APIの基本的な設定を書くファイルですので、ファイルが存在しないと起動できません。")
		panic("起動に失敗しました。")
	}
	fileData, err = ioutil.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}

	// config.yaml からデータを取得
	err = yaml.Unmarshal(fileData, &configData)
	if err != nil {
		panic(err)
	}

	// token.yaml を読み込む
	if !isFileExist("token.yaml") {
		log.Fatalln("エラー: token.yaml が見つかりません")
		log.Fatalln("APIを利用されるのに必要なトークンを書くファイルですので、ファイルが存在しないと起動できません。")
		panic("起動に失敗しました。")
	}
	fileData, err = ioutil.ReadFile("token.yaml")
	if err != nil {
		panic(err)
	}

	// token.yaml からデータを取得
	err = yaml.Unmarshal(fileData, &tokenData)
	if err != nil {
		panic(err)
	}

	// Discordを起動する設定になっていれば
	if configData.Discord.Alarm {
		if tokenData.DiscordToken == "" {
			log.Fatalln("エラー: Discordボットのトークンが設定されていません。")
			log.Fatalln("Discord報告機能をオンにするには、報告するボットのトークンが必要です。")
			panic("起動に失敗しました。")
		}
		if configData.Discord.AdminID == "" {
			log.Fatalln("エラー: Discordボットに報告してもらうユーザーのIDが設定されていません。")
			log.Fatalln("Discord報告機能をオンにするには、報告してもらうユーザーのIDが必要です。")
			log.Fatalln("もしユーザーIDを確認できない場合は、Discordの設定の「テーマ」の「詳細設定」からユーザーIDを確認できる設定にできます。")
			panic("起動に失敗しました。")
		}
		// Discordボットを起動
		go discordInit()
	}
}

// discordInit はDiscordボットを準備するための関数
func discordInit() {
	var err error

	dg, err = discordgo.New("Bot " + tokenData.DiscordToken)
	if err != nil {
		log.Fatalln("エラー: Discordボットの起動に失敗しました。")
		log.Fatalln(err)
		panic("起動に失敗しました。")
	}

	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		log.Fatalln("エラー: Discordボットの起動に失敗しました。")
		log.Fatalln(err)
		panic("起動に失敗しました。")
	}

	fmt.Println("Discordボットを起動しました。")
	// Discordボットを稼働
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	defer dg.Close()
}

// isFileExist は特定のファイルが見つかればtrueを返す関数
func isFileExist(fileName string) bool {
	// ファイル情報を取得できなかったら
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return false
	}
	return true
}
