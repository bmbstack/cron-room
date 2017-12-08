package main

import (
	"github.com/urfave/cli"
	"github.com/labstack/gommon/color"
	"github.com/onestack/cron-room/job"
	. "github.com/onestack/cron-room/helper"
	"fmt"
	"os"
	//_ "net/http/pprof"
	//"net/http"
	"time"
)

func Commands() []cli.Command {
	return []cli.Command{
		// 输出小组排行榜
		{
			Name:  "update-room",
			Usage: "Update room",
			Action: func(c *cli.Context) {
				job.RunUpdateRoom()
			},
		},
	}
}

func main() {

	// http://localhost:9999/debug/pprof/
	//go func() {
	//	http.ListenAndServe("0.0.0.0:9999", nil)
	//}()

	app := cli.NewApp()
	app.Name = "cron-room"
	app.Usage = "A cron-room application powered by cron framework"
	app.UsageText = "go run main.go command"
	app.Author = "wangmingjob"
	app.Email = "wangmingjob@icloud.com"
	app.Version = "0.0.1"
	app.Commands = Commands()

	fmt.Println(fmt.Sprintf("%s%s%s%s",
		color.White(Line),
		color.Bold(color.Green("Running")),
		color.Bold(color.Yellow("["+time.Now().Format(DateFullLayout))+"]"),
		color.White(Line)))
	fmt.Println(color.Bold(color.White("包含以下任务:")))
	for key, command := range app.Commands {
		fmt.Println(fmt.Sprintf("任务%d：%s	%s", key+1, command.Name, command.Usage))
	}
	app.Run(os.Args)
}
