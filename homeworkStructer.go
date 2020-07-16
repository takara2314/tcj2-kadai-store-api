package main

import (
	"fmt"
	"strings"
	"time"
)

// homeworkStructer はDeviorsで出力された課題データを構造体に収納する関数
func homeworkStructer(oList []string) {
	var homeworkInfo []string
	var homeworkSlice []HomeworkStruct
	var elementsNo int
	var dueTime time.Time
	var checkLock bool = false

	for _, str := range oList {
		// - (prefix): 教科名
		// ・(prefix): 課題
		if strings.HasPrefix(str, "- ") {
			// 次以降表示される課題の教科名が、有効なものであれば教科番号(要素数)を返す
			elementsNo = subjectFinder(strings.TrimLeft(str, "- "))
			// 見つからなかった場合
			if elementsNo == -1 {
				checkLock = true
				continue
			} else {
				checkLock = false
			}
		} else if strings.HasPrefix(str, "・") && !checkLock {
			// 課題の情報 (名前、ID、期限)
			homeworkInfo = strings.Split(strings.TrimLeft(str, "・"), "\t")

			// 課題の期限 (time.Time型)
			// ついでに時刻データのタイムゾーンをUTCからJSTに変更
			dueTime, _ = time.Parse("2006-01-02T15:04:05Z", homeworkInfo[2])
			dueTime = timeDiffConv(dueTime)

			homeworkSlice = append(homeworkSlice, HomeworkStruct{
				Subject: syllabusSubjectNames[elementsNo],
				Omitted: omittedSubjectNames[elementsNo],
				Name:    homeworkInfo[0],
				ID:      homeworkInfo[1],
				Due:     dueTime,
			})

			fmt.Println(
				syllabusSubjectNames[elementsNo],
				omittedSubjectNames[elementsNo],
				homeworkInfo[0],
				homeworkInfo[1],
				dueTime,
			)
		}
	}

	// 課題スライスを最後に課題スライス(総合)に入れる
	homeworksData.Homeworks = homeworkSlice
}

// subjectFinder はTeamsチーム名(before)とリンクする教科名を探し、教科番号(要素数)を返す関数
func subjectFinder(bSubjectName string) int {
	for i, subjectName := range teamsSubjectNames {
		if subjectName == bSubjectName {
			return i
		}
	}
	// 教科名が見つからなかった場合
	return -1
}

// timeDiffConv は時差変換をして返す関数
func timeDiffConv(tTime time.Time) (rTime time.Time) {
	// よりUTCらしくする
	rTime = tTime.UTC()

	// UTC → JST
	var jst *time.Location = time.FixedZone("Asia/Tokyo", 9*60*60)
	rTime = rTime.In(jst)

	return
}
