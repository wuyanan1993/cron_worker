package main

import (
	"cron_worker/http_server"
	"cron_worker/tasks"
)

func main() {
	InitTask()
}

func InitTask() {
	// 初始化任务
	//ticker := time.NewTicker(3 * time.Second)
	//for range ticker.C {
	//	fmt.Println("重复任务执行中")
	//}
	http_server.InitHttpServer()
	tasks.StartTaskConsumer()
	tasks.StartConcurrencyControl()
	tasks.StartTaskProducer()

}
