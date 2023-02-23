package starry

import (
	"encoding/json"
	"log"
	"sync"
)

const (
	INIT      = "init"
	ATTACK    = "attack"
	RANK      = "rank"
	STATUS    = "status"
	ROUND     = "round"
	EGG       = "easterEgg"
	TIME      = "time"
	CLEAR     = "clear"
	CLEAR_ALL = "clearAll"
)

type Hub struct {
	id         uint // 不断增长的id
	clients    map[uint]*Client
	broadcast  chan []byte // 要广播的数据
	unregister chan uint   // 要去掉的客户端id
	mutex      sync.Mutex
}

var hub *Hub

// Init 初始化大屏server
func Init() {
	hub = newHub()
	go hub.run()
}

func newHub() *Hub {
	return &Hub{
		id:         0,
		clients:    make(map[uint]*Client),
		broadcast:  make(chan []byte),
		unregister: make(chan uint),
	}
}

func (h *Hub) run() {
	for {
		select {
		case clientId := <-h.unregister:
			h.mutex.Lock()
			close(h.clients[clientId].data)
			log.Println("close server 52")
			delete(h.clients, clientId)
			h.mutex.Unlock()
		case msg := <-h.broadcast:
			for clientId, client := range h.clients {
				select {
				case client.data <- msg: // 发送成功,感觉这里有死锁风险
				default: // 发送失败
					h.mutex.Lock()
					log.Println("close server 61")
					
					close(client.data)
					
					delete(h.clients, clientId)
					h.mutex.Unlock()
				}
			}
		}
	}
}

func (h *Hub) sendMessage(msgType string, data interface{}) {
	jsonData, _ := json.Marshal(&unityData{
		Type: msgType,
		Data: data,
	})
	h.broadcast <- jsonData
}

func NewRound() {
}

func SendStatus(team uint, status string) {
	sendStatus(team, status)
}
func SendAttack(from uint, to uint) {
	sendAttack(from, to)
}

func sendAttack(from uint, to uint) {
	hub.sendMessage(ATTACK, attack{
		from,
		to,
	})
}

func sendRank() {
	hub.sendMessage(RANK, rank{Teams: getInfo().Teams}) // TODO getInfo调得太频繁了
}

func sendRound(newRound uint) {
	hub.sendMessage(ROUND, round{
		Round: newRound,
	})
}

func sendStatus(team uint, statusStr string) {
	hub.sendMessage(STATUS, status{
		team,
		statusStr,
	})
}

func sendEasterEgg() {
	hub.sendMessage(EGG, nil)
}

// 发送回合剩余的时间
func sendTime(time float64) {
	hub.sendMessage(TIME, restTime{Time: time})
}

func sendClear(team uint) {
	hub.sendMessage(CLEAR, clearStatus{team})
}

// 清除所有状态
func sendClearAll() {
	hub.sendMessage(CLEAR_ALL, nil)
}
