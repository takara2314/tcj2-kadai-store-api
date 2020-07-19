package main

import (
	"strings"
	"time"
)

// homeworkStructer はDeviorsで出力された課題データを構造体に収納する関数
func homeworkStructer(oList []string) {
	var homeworkInfo []string
	var homeworkSlice []HomeworkStruct
	var homeworkSliceOnlyFuture []HomeworkStruct
	var homeworkSliceJST []HomeworkStruct
	var homeworkSliceOnlyFutureJST []HomeworkStruct
	var elementsNo int
	var dueTime time.Time
	var dueTimeJST time.Time
	var checkLock bool = false

	var subjectName, omittedName string

	for _, str := range oList {
		// - (prefix): 教科名
		// ・(prefix): 課題
		if strings.HasPrefix(str, "- ") {
			// 次以降表示される課題の教科名が、有効なものであれば教科番号(要素数)を返す
			elementsNo = subjectFinder(strings.TrimLeft(str, "- "), "teamsName")
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

			// 課題の情報 (シラバス表記の教科名、省略された教科名)
			subjectName = syllabusSubjectNames[elementsNo]
			omittedName = omittedSubjectNames[elementsNo]

			// 省略された教科名に"/"が入っていたら、正しい教科名は"/"を境にしたどちらかになる
			if strings.Contains(omittedName, "/") {
				for i, setOmittedName := range strings.Split(omittedName, "/") {
					if strings.Contains(homeworkInfo[0], setOmittedName) {
						subjectName = strings.Split(subjectName, "/")[i]
						omittedName = strings.Split(omittedName, "/")[i]
					}
				}
			}

			// 課題の期限 (time.Time型)
			// ついでに時刻データのタイムゾーンをUTCからJSTに変更
			dueTime, _ = time.Parse("2006-01-02T15:04:05Z", homeworkInfo[2])
			dueTimeJST = timeDiffConv(dueTime)

			homeworkSlice = append(homeworkSlice, HomeworkStruct{
				Subject: subjectName,
				Omitted: omittedName,
				Name:    homeworkInfo[0],
				ID:      homeworkInfo[1],
				Due:     dueTime,
			})
			homeworkSliceJST = append(homeworkSliceJST, HomeworkStruct{
				Subject: subjectName,
				Omitted: omittedName,
				Name:    homeworkInfo[0],
				ID:      homeworkInfo[1],
				Due:     dueTimeJST,
			})

			// 提出期限が現在時刻より後の場合
			if dueTime.After(time.Now()) {
				homeworkSliceOnlyFuture = append(homeworkSliceOnlyFuture, HomeworkStruct{
					Subject: subjectName,
					Omitted: omittedName,
					Name:    homeworkInfo[0],
					ID:      homeworkInfo[1],
					Due:     dueTime,
				})
				homeworkSliceOnlyFutureJST = append(homeworkSliceOnlyFutureJST, HomeworkStruct{
					Subject: subjectName,
					Omitted: omittedName,
					Name:    homeworkInfo[0],
					ID:      homeworkInfo[1],
					Due:     dueTimeJST,
				})
			}
		}
	}

	// devoirsから取得した時刻を課題スライス(総合)に入れる
	homeworksData.Acquisition = time.Now()
	homeworksDataOnlyFuture.Acquisition = time.Now()
	// 課題スライスを最後に課題スライス(総合)に入れる
	homeworksData.Homeworks = homeworkSlice
	homeworksDataOnlyFuture.Homeworks = homeworkSliceOnlyFuture
}

// subjectFinder は指定したタイプの教科名とリンクする教科番号(要素数)を返す関数
func subjectFinder(bSubjectName string, beforeType string) int {
	switch beforeType {
	case "teamsName":
		for i, subjectName := range teamsSubjectNames {
			if subjectName == bSubjectName {
				return i
			}
		}
	case "syllabusName":
		for i, subjectName := range syllabusSubjectNames {
			if subjectName == bSubjectName {
				return i
			}
		}
	case "omittedName":
		for i, subjectName := range omittedSubjectNames {
			if subjectName == bSubjectName {
				return i
			}
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
