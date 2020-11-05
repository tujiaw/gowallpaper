package main

import (
	"bufio"
	"fmt"
	"gowallpaper/bing"
	"log"
	"os"
	"strings"
	"time"
)

type Cmd struct {
	Name   string
	Desc   string
	Action func(string)
}

func main() {
	newBing := bing.NewBing()

	var Cmds []Cmd
	Cmds = append(Cmds, Cmd{
		"day",
		"每天更新壁纸",
		func(params string) {
			newBing.Day()
		},
	})
	Cmds = append(Cmds, Cmd{
		"now",
		"设置当天壁纸",
		func(params string) {
			newBing.Now()
		},
	})
	Cmds = append(Cmds, Cmd{
		"prev",
		"设置前一天壁纸",
		func(params string) {
			newBing.Prev()
		},
	})
	Cmds = append(Cmds, Cmd{
		"next",
		"设置后一天壁纸",
		func(params string) {
			newBing.Next()
		},
	})
	Cmds = append(Cmds, Cmd{
		"rand",
		"间隔随机切换壁纸（如每分钟切换壁纸：rand 1m）",
		func(params string) {
			d, err := time.ParseDuration(params)
			if err != nil {
				log.Println("time.ParseDuration", err)
			} else {
				newBing.Rand(d)
			}
		},
	})
	Cmds = append(Cmds, Cmd{
		"quit",
		"退出",
		func(params string) {
			os.Exit(0)
		},
	})

	fmt.Println("设置微软必应壁纸，输入如下命令：")
	for _, cmd := range Cmds {
		fmt.Println(cmd.Name, "\t-", cmd.Desc)
	}

	for {
		fmt.Print("# ")
		line := ReadLine()
		line = strings.Trim(line, " ")
		strList := strings.Split(line, " ")
		if len(strList) == 0 {
			continue
		}

		cmd := strList[0]
		var params string
		if len(strList) > 1 {
			params = strList[1]
		}

		if cmd == "quit" || cmd == "exit" {
			os.Exit(0)
			return
		}

		for _, v := range Cmds {
			if v.Name == cmd {
				v.Action(params)
				break
			}
		}
	}
}

func ReadLine() string {
	reader := bufio.NewReader(os.Stdin)
	data, _, _ := reader.ReadLine()
	return string(data)
}
