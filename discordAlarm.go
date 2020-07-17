package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// discordAlarm はエラーが発生したときにDiscordのDMで僕に報告する関数
func discordAlarm(greetingM string, errorM error) {
	// 指定したユーザーのDMに挨拶とエラー含めた内容を送る
	dmChannel, err := dg.UserChannelCreate(adminDiscordID)
	if err != nil {
		panic(err)
	}
	dg.ChannelMessageSend(dmChannel.ID,
		fmt.Sprintf("%s\n```\n%s\n```", greetingM, fmt.Sprint(errorM)))
}

// messageCreate はDiscordBotで投稿をする関数
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// 応答確認用メッセージ
	if m.Content == "::mechaTakaran ping" {
		s.ChannelMessageSend(m.ChannelID, "ボットは正常に稼働しています。")
	}
}
