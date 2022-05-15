package main

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

var upGrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Msg 消息推送
func Msg(c *gin.Context) {
	//验证token
	token_ := c.Query("token")
	err := JWTWebsocket(token_)
	if err != nil {
		c.JSON(200, gin.H{
			"code": 2001,
			"msg":  err.Error(),
		})
		return
	}
	//升级get请求为webSocket协议
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()
	for {
		//读取ws中的数据
		mt, message, err := ws.ReadMessage()
		if err != nil {
			break
		}
		//fmt.Println(mt)
		if string(message) == "haha" {
			message = []byte("这是一个测试捏")
		}
		//写入ws数据
		err = ws.WriteMessage(mt, message)
		if err != nil {
			break
		}
		//time.Sleep(time.Second * 5)
		//mt := 1

		message = []byte("{data:'这是第二个测试'}")
		err = ws.WriteMessage(mt, message)
		if err != nil {
			break
		}
		//mt += 1
	}
}
