package main

import (
	"strings"
	"time"
)

// homeworkStructer はDeviorsで出力された課題データを構造体に収納する関数
func homeworkStructer(oList []string) {
	var nowSubject string
	var nowHomeworkInfo []string
	var homeworkSlice []HomeworkStruct
	var subjectName string
	var omittedName string
	var dueTime time.Time

	for _, str := range oList {
		// - (prefix): 教科名
		// ・(prefix): 課題
		if strings.HasPrefix(str, "- ") {
			// 次以降表示される課題の教科名
			nowSubject = strings.TrimLeft(str, "- ")
		} else {
			// 課題の情報 (名前、ID、期限)
			nowHomeworkInfo = strings.Split(strings.TrimLeft(str, "・"), "\t")

			// 課題の教科名 (シラバス版)
			subjectName = subjectFinder(nowSubject, "syllabus")
			// 課題の教科名 (省略版)
			omittedName = subjectFinder(nowSubject, "omitted")
			// 課題の期限 (time.Time型)
			dueTime, _ = time.Parse("2006-01-02T15:04:05Z", nowHomeworkInfo[2])

			homeworkSlice = append(homeworkSlice, HomeworkStruct{
				Subject: subjectName,
				Omitted: omittedName,
				Name:    nowHomeworkInfo[0],
				ID:      nowHomeworkInfo[1],
				Due:     dueTime,
			})
		}
	}

	// 課題スライスを最後に課題スライス(総合)に入れる
	homeworksData.Homeworks = homeworkSlice
}

// subjectFinder はTeamsチーム名(before)とリンクする教科名を探し、利用しやすい名前にする関数
func subjectFinder(bSubjectName string, convTo string) string {
	// syllabus… シラバス名 ([024]J2_基礎数学4（2020） → 基礎数学４)
	// omitted…  省略名     ([024]J2_基礎数学4（2020） → 数学)
	if convTo == "omitted" {
		for _, subjectName := range omittedSubjectNames {
			if subjectName == bSubjectName {
				return bSubjectName
			}
		}
	} else {
		for _, subjectName := range syllabusSubjectNames {
			if subjectName == bSubjectName {
				return bSubjectName
			}
		}
	}

	// 教科名が見つからなかった場合
	return "教科"
}
