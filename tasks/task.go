package tasks

import "fmt"

type Task interface {
	Run()
}

type ReplicationLagChecker struct {
	TaskName string
	TaskType string
}

func (t *ReplicationLagChecker) Run() {
	fmt.Println("正在执行ReplicationLagChecker!")
}

type BinVersionChecker struct {
	TaskName string
	TaskType string
}

func (t *BinVersionChecker) Run() {
	fmt.Println("正在执行BinVersionChecker!")
}

type ShardTopologyChecker struct {
	TaskName string
	TaskType string
}

func (t *ShardTopologyChecker) Run() {
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
