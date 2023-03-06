package tasks

import (
	"fmt"
	"time"
)

var TaskChannel = make(chan Task, 10000)
var concurrencyChannel = make(chan struct{}, 2)

type Task interface {
	Run(chan struct{})
	GetTimeoutSecond() int
}

type ReplicationLagChecker struct {
	TaskName      string
	TaskType      string
	TimeoutSecond int
}

// Run TODO: 需要加上超时控制，最大执行时间可以从 checker 属性TimeoutSecond获取
func (r *ReplicationLagChecker) Run(conC chan struct{}) {
	defer func() {
		// 运行结束之后给并发的channel里写入一个struct，相当于释放ticket，以供其他任务使用
		conC <- struct{}{}
	}()
	time.Sleep(time.Second)
	fmt.Println("正在执行ReplicationLagChecker!")
}

func (r *ReplicationLagChecker) GetTimeoutSecond() int {
	if r.TimeoutSecond != 0 {
		return r.TimeoutSecond
	}
	return 10
}

type BinVersionChecker struct {
	TaskName      string
	TaskType      string
	TimeoutSecond int
}

func (b *BinVersionChecker) Run(conC chan struct{}) {
	defer func() {
		// 运行结束之后给并发的channel里写入一个struct，相当于释放ticket，以供其他任务使用
		conC <- struct{}{}
	}()
	time.Sleep(time.Second)
	fmt.Println("正在执行BinVersionChecker!")
}

func (b *BinVersionChecker) GetTimeoutSecond() int {
	if b.TimeoutSecond != 0 {
		return b.TimeoutSecond
	}
	return 10
}

type ShardTopologyChecker struct {
	TaskName      string
	TaskType      string
	TimeoutSecond int
}

func (t *ShardTopologyChecker) Run(conC chan struct{}) {
	defer func() {
		// 运行结束之后给并发的channel里写入一个struct，相当于释放ticket，以供其他任务使用
		conC <- struct{}{}
	}()
	time.Sleep(time.Second)
	fmt.Println("正在执行ShardTopologyChecker!")
}

func (t *ShardTopologyChecker) GetTimeoutSecond() int {
	if t.TimeoutSecond != 0 {
		return t.TimeoutSecond
	}
	return 10
}

func GetAllCronTaskList() []Task {
	// 初始化任务列表, 并写入到local channel 和 remote database
	taskList := []Task{
		&ShardTopologyChecker{TaskName: "shardTopology", TimeoutSecond: 5},
		&ReplicationLagChecker{TaskName: "replicationLag", TimeoutSecond: 5},
		&BinVersionChecker{TaskName: "binVersion", TimeoutSecond: 5},
	}
	// TODO: 写入到 remote database 中
	return taskList
}

func CurrentTaskChannelCount() {

}

// StartTaskConsumer 启动任务消费
func StartTaskConsumer() {
	// 初始化任务 channel
	// 初始化任务生成器
	taskProducerTicker := time.NewTicker(3 * time.Second)
	go func() {
		for range taskProducerTicker.C {
			fmt.Println("模拟巡检任务生成中。。。")
			//根据taskList循环写入到task channel里边
			// 控制task channel 的写入上限 为总大小的80%，超过的话就不往里边添加任务了
			if len(TaskChannel) > 8000 {
				continue
			}
			taskList := GetAllCronTaskList()
			for _, v := range taskList {
				TaskChannel <- v
			}
		}
	}()
}

func StartConcurrencyControl() {
	// 初始化并发控制channel，从这个channel 中拿到消息才能进行任务执行
	for i := 0; i < cap(concurrencyChannel); i++ {
		concurrencyChannel <- struct{}{}
	}
}

func StartTaskProducer() {
	// 初始化任务执行器
	for {
		_ = <-concurrencyChannel
		t := <-TaskChannel
		// 此处拿到票据才能执行，需要控制并发度
		go t.Run(concurrencyChannel)
	}
}
