package main

import (
	"cron_worker/tasks"
	"fmt"
	"time"
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
	// 初始化任务 channel
	taskChannel := make(chan tasks.Task, 10000)
	// 初始化并发控制channel，从这个channel 中拿到消息才能进行任务执行
	concurrencyChannel := make(chan struct{}, 2)
	for i := 0; i < cap(concurrencyChannel); i++ {
		concurrencyChannel <- struct{}{}
	}
	// TODO: 需要启动一个http server，提供接口来接受新的巡检任务
	// 初始化任务生成器
	taskProducerTicker := time.NewTicker(3 * time.Second)
	go func() {
		for range taskProducerTicker.C {
			fmt.Println("模拟巡检任务生成中。。。")
			//根据taskList循环写入到task channel里边
			// 控制task channel 的写入上限 为总大小的80%，超过的话就不往里边添加任务了
			if len(taskChannel) > 8000 {
				continue
			}
			taskList := tasks.GetAllCronTaskList()
			for _, v := range taskList {
				taskChannel <- v
			}
		}
	}()
	// 初始化任务执行器
	for {
		_ = <-concurrencyChannel
		task := <-taskChannel
		// 此处拿到票据才能执行，需要控制并发度
		go task.Run(concurrencyChannel)
	}

}
