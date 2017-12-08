package job

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"sync/atomic"
	"github.com/levigross/grequests"
	"math/rand"
	"fmt"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	. "github.com/onestack/cron-room/helper"
	"go.uber.org/zap"
	"runtime"
	"strconv"
)

var (
	crawledUrlsNum uint64

	crawlLimit uint64 			= 100
	numThreads uint64 			= 100
	crawlDelaySeconds uint64 	= 0

	crawlWaitGroup 	sync.WaitGroup
	requestQueue 	= make(chan *RequestEntity, 10)
	responseQueue 	= make(chan *ResponseEntity, 15)
	shutdownNotify 	= make(chan struct{})
	statusRequest 	= make(chan os.Signal, 1)

	logger *zap.Logger
)

func init() {
	logger = CreateLogger("crawler")
}

type RequestEntity struct {
	Tag            string
	URL            string
	RequestOptions *grequests.RequestOptions
	DataCallback   func(json string)
}

type ResponseEntity struct {
	Tag  string
	Resp *grequests.Response
}

func status() {
	for {
		select {
		case <-statusRequest:
			crawledNum := strconv.FormatUint(atomic.LoadUint64(&crawledUrlsNum), 10)
			log.Println("爬虫状态", fmt.Sprintf("请求队列长度:%d, 响应队列长度:%d, 爬取API数量:%s", len(requestQueue), len(responseQueue), crawledNum))
			writeInfo("[爬虫状态]", fmt.Sprintf("请求队列长度:%d, 响应队列长度:%d, 爬取API数量:%s", len(requestQueue), len(responseQueue), crawledNum))
		case <-shutdownNotify:
			log.Println("爬虫状态：接收到关闭信号")
			writeInfo("[爬虫状态]", "接收到关闭信号")
			signal.Stop(statusRequest)
			return
		}
	}
}

func spiders() {
	crawlWaitGroup.Add(1)
	for {
		select {
		case requestEntity := <-requestQueue:
			time.Sleep(time.Second * time.Duration(crawlDelaySeconds))
			resp, err := grequests.Get(requestEntity.URL, requestEntity.RequestOptions)
			if err != nil {
				writeInfo("[响应结果]", fmt.Sprintf("Tag: %s, Error: %s", requestEntity.Tag, err.Error()))
				continue
			}

			if !resp.Ok {
				writeInfo("[响应结果]", fmt.Sprintf("Tag: %s, Error: Http response status code is not 200", requestEntity.Tag))
				continue
			}

			if requestEntity.DataCallback != nil {
				requestEntity.DataCallback(resp.String())
			}

			// Crawl data OK
			writeInfo("[响应正常]", fmt.Sprintf("爬取成功数量:%s", strconv.FormatUint(atomic.LoadUint64(&crawledUrlsNum), 10)))
			atomic.AddUint64(&crawledUrlsNum, 1)
			responseQueue <- CreateResponseEntity(requestEntity.Tag, resp)
		case <-shutdownNotify:
			crawlWaitGroup.Done()
			return

		}
	}
}

func eaters() {
	crawlWaitGroup.Add(1)
	for {
		select {
		case <-shutdownNotify:
			crawlWaitGroup.Done()
			return
		case responseEntity := <-responseQueue:
			resp := responseEntity.Resp
			if atomic.LoadUint64(&crawledUrlsNum) == crawlLimit {
				log.Println("爬虫状态： 关闭爬虫，爬取成功数量:", atomic.LoadUint64(&crawledUrlsNum))
				writeInfo("[爬虫状态]", fmt.Sprintf("关闭爬虫，爬取成功数量:%s", strconv.FormatUint(atomic.LoadUint64(&crawledUrlsNum), 10)))
				resp.Close()
				close(shutdownNotify)
				continue
			}
			resp.Close()

		}

	}
}

func CreateRequestEntity(tag string, requestURL string, params map[string]string, dataCallback func(json string)) *RequestEntity {
	return &RequestEntity{
		Tag: tag,
		URL: requestURL,
		RequestOptions: &grequests.RequestOptions{
			UserAgent: getUserAgent(),
			Params:    params,
		},
		DataCallback: dataCallback,
	}
}

func CreateResponseEntity(tag string, resp *grequests.Response) *ResponseEntity {
	return &ResponseEntity{
		Tag:  tag,
		Resp: resp,
	}
}

func CreateLogger(name string) *zap.Logger {
	filename := fmt.Sprintf("%s%s.log", GetEnv().LogFilePathPrefix, name)
	ws := zapcore.AddSync(&lumberjack.Logger{
		Filename:   filename,
		MaxSize:    500, // megabytes
		MaxBackups: 3, // backup
		MaxAge:     30, // days
	})
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = ""
	coreInstance := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		ws,
		zap.InfoLevel,
	)
	return zap.New(coreInstance)
}

func writeInfo(msg string, data string) {
	logger.Info(msg,
		zap.String(LogKeySource, "[爬虫]"),
		zap.String(LogKeyTime, time.Now().Format(DateFullLayout)),
		zap.String(LogKeyData, data),
		zap.Int(LogKeyGoroutineNum, runtime.NumGoroutine()),
	)
}

func getUserAgent() string {
	var userAgent = [...]string{
		"Mozilla/5.0 (compatible, MSIE 10.0, Windows NT, DigExt)",
		"Mozilla/4.0 (compatible, MSIE 7.0, Windows NT 5.1, 360SE)",
		"Mozilla/4.0 (compatible, MSIE 8.0, Windows NT 6.0, Trident/4.0)",
		"Mozilla/5.0 (compatible, MSIE 9.0, Windows NT 6.1, Trident/5.0,",
		"Opera/9.80 (Windows NT 6.1, U, en) Presto/2.8.131 Version/11.11",
		"Mozilla/4.0 (compatible, MSIE 7.0, Windows NT 5.1, TencentTraveler 4.0)",
		"Mozilla/5.0 (Windows, U, Windows NT 6.1, en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
		"Mozilla/5.0 (Macintosh, Intel Mac OS X 10_7_0) AppleWebKit/535.11 (KHTML, like Gecko) Chrome/17.0.963.56 Safari/535.11",
		"Mozilla/5.0 (Macintosh, U, Intel Mac OS X 10_6_8, en-us) AppleWebKit/534.50 (KHTML, like Gecko) Version/5.1 Safari/534.50",
		"Mozilla/5.0 (Linux, U, Android 3.0, en-us, Xoom Build/HRI39) AppleWebKit/534.13 (KHTML, like Gecko) Version/4.0 Safari/534.13",
		"Mozilla/5.0 (iPad, U, CPU OS 4_3_3 like Mac OS X, en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8J2 Safari/6533.18.5",
		"Mozilla/4.0 (compatible, MSIE 7.0, Windows NT 5.1, Trident/4.0, SE 2.X MetaSr 1.0, SE 2.X MetaSr 1.0, .NET CLR 2.0.50727, SE 2.X MetaSr 1.0)",
		"Mozilla/5.0 (iPhone, U, CPU iPhone OS 4_3_3 like Mac OS X, en-us) AppleWebKit/533.17.9 (KHTML, like Gecko) Version/5.0.2 Mobile/8J2 Safari/6533.18.5",
		"MQQBrowser/26 Mozilla/5.0 (Linux, U, Android 2.3.7, zh-cn, MB200 Build/GRJ22, CyanogenMod-7) AppleWebKit/533.1 (KHTML, like Gecko) Version/4.0 Mobile Safari/533.1",
	}
	var r = rand.New(rand.NewSource(time.Now().UnixNano()))
	return userAgent[r.Intn(len(userAgent))]
}

func AddRequestEntity(entity *RequestEntity) {
	requestQueue <- entity
}

func StartCrawl() {
	signal.Notify(statusRequest, syscall.SIGHUP)
	go status()
	go eaters()
	for i := uint64(0); i < numThreads; i++ {
		go spiders()
	}

	<-shutdownNotify
	crawlWaitGroup.Wait()
}
