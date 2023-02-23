package starry

import (
	"Evo/config"
	"Evo/db"
	"Evo/model"
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"time"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

type Client struct {
	conn *websocket.Conn
	data chan []byte
}

// handle 客户端连接
func (c *Client) handle(id uint) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		_ = c.conn.Close()
	}()

	// 发送初始化数据
	initMsg, _ := json.Marshal(&unityData{
		Type: INIT,
		Data: getInfo(),
	})
	c.data <- initMsg // 发送初始的数据

	for {
		select {
		case msg, ok := <-c.data:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			_, _ = w.Write(msg)
			if err := w.Close(); err != nil {
				log.Printf("Failed to close %v\n", err)
				delete(hub.clients, id)
				close(c.data)
				return
			}

		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				delete(hub.clients, id)
				close(c.data)
				return
			}
		}
	}
}

func getInfo() Info {
	var teams []model.Team
	var unityTeams []Team
	db.DB.Model(&model.Team{}).Order("score DESC").Find(&teams)
	for rank, team := range teams {
		unityTeams = append(unityTeams, Team{
			Id:    int(team.ID),
			Name:  team.Name,
			Rank:  rank + 1,
			Img:   team.Logo,
			Score: int(team.Score),
		})
	}
	var info Info

	info.Teams = unityTeams
	info.Title = config.GAME_NAME
	info.Round = int(config.ROUND_NOW)
	info.Time = config.GetRoundRemainTime() // TODO 本轮剩余时间

	return info
}
