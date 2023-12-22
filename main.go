package main

import (
	"go_gin_websocket/WebSocketHandler"

	"github.com/gin-gonic/gin"
)

func main() {
	// gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/ws", WebSocketHandler.WebSocketBase)
	//r.POST("/to_ws", to_ws)
	r.Run(":9501")
	//err := r.Run(":9501")
	//if err != nil {
	//	return
	//}
}

// 协程调用，待完善
//func to_ws(c *gin.Context) {
//	// 创建在 goroutine 中使用的副本
//	cCp := c.Copy()
//	// 验证ip来源
//	//ip := cCp.Request.Header.Get("X-Forward-For")
//	//if ip != WebSocketHandler.REMOTE_ADDR {
//	//	return
//	//}
//	// 声明接收的变量
//	var to_ws_data ToWsData
//	// 将request的body中的数据，自动按照json格式解析到结构体
//	if err := cCp.ShouldBindJSON(&to_ws_data); err != nil {
//		// 返回错误信息
//		// gin.H封装了生成json数据的工具
//		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//	source := to_ws_data.Source
//	cid := to_ws_data.Cid
//	if source != "php" {
//		return
//	}
//	go WebSocketHandler.TwSetData(cid, to_ws_data)
//	c.JSON(200, gin.H{
//		"message": "OK",
//	})
//}
