package starry

import (
	"Evo/util"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

func ServeWebsocket(c *gin.Context) {
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}).Upgrade(c.Writer, c.Request, nil)

	// 建立连接失败，响应500
	if err != nil {
		log.Println(err.Error())
		c.JSON(500, nil)
	}
	client := &Client{
		conn: conn,
		data: make(chan []byte, 256),
	}
	hub.mutex.Lock()
	temp := hub.id
	hub.clients[hub.id] = client
	hub.id++
	hub.mutex.Unlock()
	// 建立连接
	go client.handle(temp)
}

func Attack(c *gin.Context) {
	type attackForm struct {
		From uint
		To   uint
	}
	attack := attackForm{}
	if err := c.ShouldBind(&attack); err != nil {
		util.Fail(c, "发送失败", nil)
		return
	}
	if attack.From == attack.To {
		util.Fail(c, "参数错误", nil)
		return
	}
	sendAttack(attack.From, attack.To)
	util.Success(c, "success", nil)
}

func Rank(c *gin.Context) {
	sendRank()
	util.Success(c, "success", nil)
}

func Status(c *gin.Context) {
	var status struct {
		Id     uint   `binding:"required"`
		Status string `binding:"required"`
	}
	if err := c.ShouldBind(&status); err != nil {
		util.Fail(c, "发送失败", nil)
		return
	}
	if status.Status != "down" && status.Status != "attacked" {
		util.Fail(c, "请选择正确的状态", nil)
		return
	}
	sendStatus(status.Id, status.Status)
	util.Success(c, "success", nil)
}

func Round(c *gin.Context) {
	var round struct {
		Round uint `binding:"required"`
	}
	if err := c.ShouldBind(&round); err != nil {
		util.Fail(c, "发送失败", nil)
		return
	}
	sendRound(round.Round)
	util.Success(c, "success", nil)
}

func Time(c *gin.Context) {
	var time struct {
		Time float64 `binding:"required"`
	}
	if err := c.ShouldBind(&time); err != nil {
		util.Fail(c, "发送失败", nil)
		return
	}
	sendTime(time.Time)
	util.Success(c, "success", nil)
}

func Clear(c *gin.Context) {
	var clear struct {
		Id uint `binding:"required"`
	}
	if err := c.ShouldBind(&clear); err != nil {
		util.Fail(c, "发送失败", nil)
		return
	}
	sendClear(clear.Id)
	util.Success(c, "success", nil)
}

func ClearAll(c *gin.Context) {
	sendClearAll()
	util.Success(c, "success", nil)
}
