package main

import "time"

// requestManager は10分ごとに各トークンの残りのGETリクエスト可能回数をリセットする関数
func requestManager() {
	// 残り回数をリセット
	for _, token := range tokenData.AllowedTokens {
		tokenLimit[token] = configData.GETLimit
	}

	for {
		var nowMinute int = time.Now().Minute()

		// 指定した時間になったら残り回数をリセット
		if containsInt([]int{0, 10, 20, 30, 40, 50}, nowMinute) {
			for _, token := range tokenData.AllowedTokens {
				tokenLimit[token] = configData.GETLimit
			}
			time.Sleep(1 * time.Minute)
		}
	}
}
