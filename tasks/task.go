package tasks

import (
	"cron_worker/store"
	"fmt"
	"time"
)

var TaskChannel = make(chan Task, 10000)
var concurrencyChannel = make(chan struct{}, 1)

type Task interface {
	Run(chan struct{})
	GetTimeoutSecond() time.Duration
}

type Checker struct {
	TaskName      string        `json:"task_name"`
	TaskType      string        `json:"task_type"`
	TimeoutSecond time.Duration `json:"timeout_second"`
}

type ReplicationLagChecker struct {
	Checker
}

func (r *ReplicationLagChecker) Run(conC chan struct{}) {
	defer func() {
		// 运行结束之后给并发的channel里写入一个struct，相当于释放ticket，以供其他任务使用
		conC <- struct{}{}
	}()
	time.Sleep(4 * time.Second)
	fmt.Println("正在执行ReplicationLagChecker!")
}

func (r *ReplicationLagChecker) GetTimeoutSecond() time.Duration {
	if r.TimeoutSecond != 0 {
		return r.TimeoutSecond
	}
	return 10
}

type BinVersionChecker struct {
	Checker
}

func (b *BinVersionChecker) Run(conC chan struct{}) {
	defer func() {
		// 运行结束之后给并发的channel里写入一个struct，相当于释放ticket，以供其他任务使用
		conC <- struct{}{}
	}()
	timeoutT := time.NewTicker(b.GetTimeoutSecond() * time.Second)
	finishedFlag := make(chan struct{})
	go func() {
		fmt.Println("正在执行BinVersionChecker!")
		time.Sleep(11 * time.Second)
	}()
	select {
	case <-timeoutT.C:
		fmt.Println("执行任务超时")
		return
	case <-finishedFlag:
		return
	}

}

func (b *BinVersionChecker) GetTimeoutSecond() time.Duration {
	if b.TimeoutSecond != 0 {
		return b.TimeoutSecond
	}
	return 10
}

type ShardTopologyChecker struct {
	Checker
}

func (t *ShardTopologyChecker) Run(conC chan struct{}) {
	defer func() {
		// 运行结束之后给并发的channel里写入一个struct，相当于释放ticket，以供其他任务使用
		conC <- struct{}{}
	}()
	time.Sleep(2 * time.Second)
	fmt.Println("正在执行 ShardTopologyChecker!")
}

func (t *ShardTopologyChecker) GetTimeoutSecond() time.Duration {
	if t.TimeoutSecond != 0 {
		return t.TimeoutSecond
	}
	return 10
}

func GetAllCronTaskList() []Task {
	// 初始化任务列表, 并写入到local channel 和 remote database
	taskList := []Task{
		&ShardTopologyChecker{Checker: Checker{TaskName: "shardTopology", TimeoutSecond: 5}},
		//&ReplicationLagChecker{TaskName: "replicationLag", TimeoutSecond: 5},
		//&BinVersionChecker{TaskName: "binVersion", TimeoutSecond: 5},
	}
	// TODO: 写入到 remote database 中
	var backend store.RemoteStore
	backend = &store.MongodbStore{ConnectionString: "localhost:8080"}
	err := backend.AddTask(taskList)
	if err != nil {
		fmt.Println("写入到远程database中失败: ", err)
	}

	return taskList
}

// StartTaskProducer 启动定时巡检任务生成
func StartTaskProducer() {
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

func StartTaskConsumer() {
	// 初始化任务执行器
	for {
		_ = <-concurrencyChannel
		t := <-TaskChannel
		// 此处拿到票据才能执行，需要控制并发度
		go t.Run(concurrencyChannel)
	}
}

func CreateChecker(c *Checker) Task {
	//&ShardTopologyChecker{Checker: Checker{TaskName: "shardTopology"}, TimeoutSecond: 5},
	//&ReplicationLagChecker{TaskName: "replicationLag", TimeoutSecond: 5},
	//&BinVersionChecker{TaskName: "binVersion", TimeoutSecond: 5},
	fmt.Println("add task, task_name: ", c.TaskName)
	switch c.TaskName {
	case "binVersion":
		return &BinVersionChecker{Checker: *c}
	case "replicationLag":
		return &ReplicationLagChecker{Checker: *c}
	case "shardTopology":
		return &ShardTopologyChecker{Checker: *c}
	default:
		return nil
	}

}
