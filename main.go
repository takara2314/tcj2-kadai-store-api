package main

import (
	"fmt"
	"net"
	"net/http/fcgi"
	"os/exec"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type scheduleList struct {
	Date      time.Time         `json:"data"`
	Homeworks []HomeworkStructs `json:"homeworks"`
}

type HomeworkStructs struct {
	Comment string    `json:"comment"`
	Due     time.Time `json:"due"`
	Subject string    `json:"subject"`
}

func main() {
	var err error
	r := gin.Default()

	r.GET("/", homeRequestFunc)
	r.GET("/get_old", getRequestFuncO)
	r.GET("/get", getRequestFunc)

	l, err := net.Listen("tcp", ":2314")
	if err != nil {
		panic(err)
	}
	if err := fcgi.Serve(l, r); err != nil {
		panic(err)
	}
}

func homeRequestFunc(c *gin.Context) {
	c.String(200, "Please use get or status parameters.")
}

func getRequestFunc(c *gin.Context) {
	var command *exec.Cmd = exec.Command("npm", "start")
	command.Dir = "../../devoirs"

	out, err := command.Output()
	if err != nil {
		panic(err)
	}

	var outputList []string = strings.Split(string(out), "\n")

	fmt.Println("List length:", len(outputList))

	for _, str := range outputList {
		if strings.HasPrefix(str, "・") {
			devidedList := strings.Split(str, "\t")

			for _, str2 := range devidedList {
				fmt.Println(str2)
			}

			fmt.Print("\n")
		}
	}

	c.String(200, "Finished task")
}

func getRequestFuncO(c *gin.Context) {
	userBrowser := c.Request.Header.Get("User-Agent")
	fmt.Println(userBrowser)

	returnJSON := scheduleList{
		Date: time.Now(),
		Homeworks: []HomeworkStructs{
			{
				Comment: "プリントをやってください",
				Due:     time.Now(),
				Subject: "基礎数学４",
			},
			{
				Comment: "ここしてください",
				Due:     time.Now(),
				Subject: "電気電子基礎",
			},
			{
				Comment: "感想を800字以内に書いてください",
				Due:     time.Now(),
				Subject: "一般基礎教育２",
			},
			{
				Comment: "Do the handout.",
				Due:     time.Now(),
				Subject: "Sonzaishinai English 2",
			},
		},
	}
	c.JSON(200, returnJSON)
}
