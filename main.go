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
	Name string
	Desc string
	Action func(string)
}

var Cmds []Cmd

func main() {
	bing := bing.NewBing()
	bing.Init()

	Cmds = append(Cmds, Cmd{
		"now",
		"当天壁纸",
		func(params string) {
			bing.Now()
		},
	})
	Cmds = append(Cmds, Cmd{
		"prev",
		"前一天壁纸",
		func(params string) {
			bing.Prev()
		},
	})
	Cmds = append(Cmds, Cmd{
		"next",
		"后一天壁纸",
		func(params string) {
			bing.Next()
		},
	})
	Cmds = append(Cmds, Cmd{
		"interval",
		"间隔时间切换(5m：5分钟)",
		func(params string) {
			d, err := time.ParseDuration(params)
			if err != nil {
				log.Println(err)
			} else {
				bing.Interval(d)
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

	fmt.Println("please input command, usage:")
	for _, cmd := range Cmds {
		fmt.Println(cmd.Name, "-", cmd.Desc)
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

		if (cmd == "quit" || cmd == "exit") {
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

func ReadLine()string {
	reader := bufio.NewReader(os.Stdin)
	data, _, _ := reader.ReadLine()
	return string(data)
}