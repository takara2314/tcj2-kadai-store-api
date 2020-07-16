package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// discordAlarm はエラーが発生したときにDiscordのDMで僕に報告する関数
func discordAlarm(greetingM string, errorM error) {
	fmt.Println("dAの関数までは入ったよ！")
	// // 始めの挨拶(概要説明)とエラー内容を入れる
	dmGreetingM = greetingM
	dmErrorM = fmt.Sprint(errorM)
	// // DiscordDMを送る許可を与える
	// isDiscordAlarm = true

	fmt.Println("今からDMに送るね！")
	// 拡張的な宝箱#9220(226453185613660160)のDMに挨拶とエラー含めた内容を送る
	dmChannel, err := dg.UserChannelCreate("226453185613660160")
	if err != nil {
		panic(err)
	}
	dg.ChannelMessageSend(dmChannel.ID,
		fmt.Sprintf("%s\n```\n%s\n```", dmGreetingM, dmErrorM))
}

// messageCreate はDiscordBotで投稿をする関数
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println("たからん大好き♡")
	// 応答確認用メッセージ
	if m.Content == "::mechaTakaran ping" {
		s.ChannelMessageSend(m.ChannelID, "ボットは正常に稼働しています。")
	}

	// if isDiscordAlarm {
	// 	fmt.Println("今からDMに送るね！")
	// 	isDiscordAlarm = false
	// 	// 拡張的な宝箱#9220(226453185613660160)のDMに挨拶とエラー含めた内容を送る
	// 	dmChannel, err := s.UserChannelCreate("226453185613660160")
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	s.ChannelMessageSend(dmChannel.ID,
	// 		fmt.Sprintf("%s\n```\n%s\n```", dmGreetingM, dmErrorM))
	// }
}
