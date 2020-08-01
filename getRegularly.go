package main

import (
	"os/exec"
	"strings"
	"time"
)

// getRegularly は定期的にDevoirs経由でTeamsから課題一覧を取得する関数
func getRegularly(getTime []int) {
	for {
		var nowMinute int = time.Now().Minute()

		// 指定した時間になったら実行
		if containsInt(getTime, nowMinute) {
			var command *exec.Cmd = exec.Command("npm", "start")
			command.Dir = "../../devoirs"

			out, err := command.Output()
			if err != nil && configData.Discord.Alarm {
				// API管理者にDiscordでエラーを報告し、
				// プロセスを強制終了させる
				if !discordAlarmed {
					discordAlarmed = true
					discordAlarm("Devoirsで何かエラーが発生しました！", "RDP接続したデスクトップからDevoirsを手動で起動すると1時間だけ改善されるかもしれません。", err)
				}
				continue
			}
			discordAlarmed = false

			// 実行結果を1行ずつリストに入れる
			var outputData []string = strings.Split(string(out), "\n")

			for i, str := range outputData {
				// もし教科名の行が来たら
				if strings.HasPrefix(str, "-") {
					outputData = outputData[i:]
					break
				}
			}

			// 課題関係しか乗っていない実行結果を元に、データを構造体に収納する
			homeworkStructer(outputData)

			time.Sleep(1 * time.Minute)
		}
	}
}

// containsInt はint型のスライスから特定の整数があればtrueを返す関数
func containsInt(tSlice []int, tNum int) bool {
	for _, num := range tSlice {
		if tNum == num {
			return true
		}
	}
	return false
}
