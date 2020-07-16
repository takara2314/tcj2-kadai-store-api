package main

import (
	"log"
	"os/exec"
	"strings"
	"time"
)

// getRegularly は定期的にDeviors経由でTeamsから課題一覧を取得する関数
func getRegularly(getTime []int) {
	for {
		var nowMinute int = time.Now().Minute()

		// 指定した時間になったら実行
		if containsInt(getTime, nowMinute) {
			var command *exec.Cmd = exec.Command("npm", "start")
			command.Dir = "../../devoirs"

			out, err := command.Output()
			if err != nil {
				log.Println("Devoirsで何かエラーが発生しました！")
				panic(err)
			} else {
				// fmt.Println("Devoirs側でエラーなんか起きなかったよ?")
			}

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
