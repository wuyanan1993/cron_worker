package tasks

import (
	"fmt"
	"time"
)

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
