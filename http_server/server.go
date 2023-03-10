package http_server

import (
	"cron_worker/tasks"
	"fmt"

	"github.com/gin-gonic/gin"
)

func InitHttpServer() {
	r := gin.Default()
	routeTask := r.Group("/task")
	{
		routeTask.GET("/currentCount", TaskCurrentCount)
		routeTask.PUT("/task", NewTask)
	}
	err := r.Run(":8080")
	if err != nil {
		fmt.Println("启动http 接口错误，程序将以巡检模式启动")
	}
}

func TaskCurrentCount(c *gin.Context) {
	// 获取当前待执行队列中的任务数量
	count := len(tasks.TaskChannel)
	c.JSON(200, map[string]int{"count": count})
}

func NewTask(c *gin.Context) {
	// 根据参数创建新的任务

	var checker tasks.Checker
	err := c.ShouldBindJSON(&checker)
	if err != nil {
		fmt.Println("get info from http err: ", err)
		c.JSON(400, map[string]interface{}{"msg": "err"})
	}
	newChecker := tasks.CreateChecker(&checker)
	tasks.TaskChannel <- newChecker
	c.JSON(200, map[string]interface{}{"msg": "done", "info": map[string]interface{}{"task_name": checker.TaskName, "task_type": checker.TaskType, "timeout_second": checker.TimeoutSecond}})
}
