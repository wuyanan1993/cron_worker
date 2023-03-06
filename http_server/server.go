package http_server

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

func InitHttpServer() {
	// TODO: 需要启动一个http server，提供接口来接受新的巡检任务
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
	c.JSON(200, map[string]int{"count": 100})
}

func NewTask(c *gin.Context) {
	c.JSON(200, map[string]string{"msg": "done"})
}
