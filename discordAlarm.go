package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// discordAlarm はエラーが発生したときにDiscordのDMで僕に報告する関数
func discordAlarm(description string, coping string, errorM error) {
	// 指定したユーザーのDMに挨拶とエラー含めた内容を送る
	dmChannel, err := dg.UserChannelCreate(configData.Discord.AdminID)
	if err != nil {
		panic(err)
	}
	dg.ChannelMessageSend(dmChannel.ID,
		fmt.Sprintf(configData.Discord.MessageFormat, description, coping, errorM))
}

// messageCreate はDiscordボットで投稿をする関数
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// 応答確認用メッセージ
	if m.Content == configData.Discord.CommandPrefix+" ping" {
		s.ChannelMessageSend(m.ChannelID, "ボットは正常に稼働しています。")
	}

	// バージョン確認用メッセージ
	if m.Content == configData.Discord.CommandPrefix+" version" {
		s.ChannelMessageSend(m.ChannelID, "TCJ2 Kadai Store API - v0.2.0 pre3")
	}

	// 強制停止メッセージ
	if m.Content == configData.Discord.CommandPrefix+" stop" {
		if m.Author.ID == configData.Discord.AdminID {
			s.ChannelMessageSend(m.ChannelID, "ボットを強制終了させます。")
			panic("管理者による強制終了")
		} else {
			s.ChannelMessageSend(m.ChannelID, "あなたはそのコマンドを実行することができません。")
		}
	}
}
