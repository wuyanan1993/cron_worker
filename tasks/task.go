package tasks

import (
	"fmt"
	"time"
)

type Task interface {
	Run(chan struct{})
}

type ReplicationLagChecker struct {
	TaskName string
	TaskType string
}

func (t *ReplicationLagChecker) Run(conC chan struct{}) {
	defer func() {
		// 运行结束之后给并发的channel里写入一个struct，相当于释放ticket，以供其他任务使用
		conC <- struct{}{}
	}()
	time.Sleep(time.Second)
	fmt.Println("正在执行ReplicationLagChecker!")
}

type BinVersionChecker struct {
	TaskName string
	TaskType string
}

func (t *BinVersionChecker) Run(conC chan struct{}) {
	defer func() {
		// 运行结束之后给并发的channel里写入一个struct，相当于释放ticket，以供其他任务使用
		conC <- struct{}{}
	}()
	time.Sleep(time.Second)
	fmt.Println("正在执行BinVersionChecker!")
}

type ShardTopologyChecker struct {
	TaskName string
	TaskType string
}

func (t *ShardTopologyChecker) Run(conC chan struct{}) {
	defer func() {
		// 运行结束之后给并发的channel里写入一个struct，相当于释放ticket，以供其他任务使用
		conC <- struct{}{}
	}()
	time.Sleep(time.Second)
	fmt.Println("正在执行ShardTopologyChecker!")
}

func GetCronTaskList() []Task {
	taskList := []Task{
		&ShardTopologyChecker{TaskName: "shardTopology"},
		&ReplicationLagChecker{TaskName: "replicationLag"},
		&BinVersionChecker{TaskName: "binVersion"},
	}
	return taskList
}
