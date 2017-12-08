package job

import (
	"time"
	. "github.com/onestack/cron-room/helper"
	_ "github.com/onestack/cron-room/models"
	"github.com/onestack/cron-room/models"
	"runtime"
	"github.com/bmbstack/gron"
	"github.com/bmbstack/gron/xtime"
	"fmt"
)

const (
	LogFileUpdateRoom string = "update-room"
)

type UpdateRoom struct {
	*BaseJob
}

func (this UpdateRoom) Run() {
	models.MigrateAll()

	this.WriteInfo("[房源更新]", "开始任务")

	//房源: ziroom，xiangyu，hizhu，anjuke，qfang，fangduoduo
	tagRomListHizhu := "RoomListHizhu"
	requestHizhuHouseList := CreateRequestEntity(
		tagRomListHizhu,
		"http://m.hizhu.com/Home/House/houselist.html",
		map[string]string{
			"city_code":   "001001",
			"pageno":      "1",
			"limit":       "1",
			"sort":        "-1",
			"region_id":   "",
			"plate_id":    "",
			"money_max":   "999999",
			"money_min":   "0",
			"logicSort":   "0",
			"line_id":     "0",
			"stand_id":    "0",
			"key":         "0",
			"key_self":    "0",
			"type_no":     "0",
			"search_id":   "0",
			"latitude":    "0",
			"longitude":   "0",
			"distance":    "0",
			"update_time": "0",
		},
		func(json string) {
			writeInfo(fmt.Sprintf("[DataCallback] %s", tagRomListHizhu), "Get done")
		},
	)
	AddRequestEntity(requestHizhuHouseList)

	this.WriteInfo("[房源更新]", "结束任务")
}

func RunUpdateRoom() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	cron := gron.New()

	// Add the job
	// 北京时间 = 世界时间 + 8 (At 是世界时间格林尼治标准时间)
	// 北京时间00:01  世界时间前一天的16:01
	cron.Add(gron.Every(2*xtime.Second), UpdateRoom{BaseJob: CreateJob(LogFileUpdateRoom)})

	cron.Start()
	defer cron.Stop()

	StartCrawl()

	select {
	case <-time.After(JobEndTime):
		return
	}
}
