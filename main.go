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
	// 初始化并发控制channel
	concurrencyChannel := make(chan struct{}, 2)
	for i := 0; i < cap(concurrencyChannel); i++ {
		concurrencyChannel <- struct{}{}
	}
	// 初始化任务列表，并且暴漏http 接口，可以通过api 来进行列表修改
	taskList := tasks.GetCronTaskList()
	// 初始化任务生成器
	taskProducerTicker := time.NewTicker(3 * time.Second)
	go func() {
		for range taskProducerTicker.C {
			fmt.Println("模拟巡检任务生成中。。。")
			//根据taskList循环写入到task channel里边
			// TODO: 是否控制写入的速度，如果消费积累了很多，可以暂时不往队列里放任务，加一个判断
			for _, v := range taskList {
				taskChannel <- v
			}
		}
	}()
	// 初始化任务执行器
	for {
		_ = <-concurrencyChannel
		task := <-taskChannel
		// 此处应该拿到 票据才能执行，需要控制并发度
		go task.Run(concurrencyChannel)
	}

}
