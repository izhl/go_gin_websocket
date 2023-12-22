package WebSocketHandler

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	// "time"
)

// websocket 升级并跨域
var (
	upgrade = &websocket.Upgrader{
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

// WebSocketBase TODO:服务基本函数
func WebSocketBase(c *gin.Context) {
	var (
		err error
		conn *websocket.Conn
		ws *wsConn
	)
	if conn, err = upgrade.Upgrade(c.Writer,c.Request,nil); err != nil{
		return
	}
	if ws, err = InitWebSocket(conn); err != nil{
		return
	}
	// 使得inChan和outChan耦合起来
	for {
		var data []byte
		if data, err = ws.InChanRead(); err != nil{
			goto ERR
		}
		if err = ws.OutChanWrite(data); err != nil{
			goto ERR
		}
	}
ERR:
	ws.CloseConn()
}