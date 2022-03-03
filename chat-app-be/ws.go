package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"chat-app.vi/controllers"
	errors "chat-app.vi/errors"
	db "chat-app.vi/rdb"
	services "chat-app.vi/services"
	"github.com/gorilla/websocket"
	"github.com/jackc/pgx/v4/pgxpool"
)

type WsMessageData struct {
	userId  string
	message []byte
}

type WsChan struct {
	userId string
	ws     *websocket.Conn
}

type WsMsg struct {
	Type    string
	Payload map[string]interface{}
}

type WSClientHub struct {
	clients          map[string][]*websocket.Conn
	addClientChan    chan WsChan
	removeClientChan chan WsChan
	sendMessageChan  chan WsMessageData
}

var WS_HUB *WSClientHub

func newWSHub() *WSClientHub {
	hub := &WSClientHub{
		clients:          make(map[string][]*websocket.Conn),
		addClientChan:    make(chan WsChan),
		removeClientChan: make(chan WsChan),
		sendMessageChan:  make(chan WsMessageData),
	}

	return hub
}

func initWSHub() {
	WS_HUB = newWSHub()

	go wsHandler()
}

func wsHandler() {
	for {
		select {
		case data := <-WS_HUB.addClientChan:
			WS_HUB.addWSClient(data.userId, data.ws)
			addUserWsIP(data.userId)
			go listenForWSMessages(data.userId, data.ws)
		case data := <-WS_HUB.removeClientChan:
			WS_HUB.removeWSClient(data.userId, data.ws)
		case data := <-WS_HUB.sendMessageChan:
			WS_HUB.sendMessageToWsClient(data.message, data.userId)
		}
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WsHandler(w http.ResponseWriter, r *http.Request) {
	err := db.POOL.HandleFunc(func(conn *pgxpool.Conn) error {
		authService := services.NewAuthService(conn)
		loggedUserId, err := authService.GetLoggedUserId(r.FormValue("token"))
		if loggedUserId == "" || err != nil {
			return &errors.UnauthorizedError{Err: err}
		}

		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return err
		}

		MB.subscriptionChan <- loggedUserId
		WS_HUB.addClientChan <- WsChan{
			userId: loggedUserId,
			ws:     ws,
		}

		return nil
	})

	controllers.SendErrorRespIfHas(w, err)
}

func listenForWSMessages(userId string, ws *websocket.Conn) {
	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			WS_HUB.removeClientChan <- WsChan{
				userId: userId,
				ws:     ws,
			}

			return
		}

		payload := &WsMsg{}
		err = json.Unmarshal(msg, payload)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("WS MSG RECEIVED:", payload)

		err = db.POOL.HandleFunc(func(conn *pgxpool.Conn) error {
			if payload.Type == "msg" {
				tx, err := conn.Begin(context.Background())
				if err != nil {
					return err
				}

				err = handleCreateMessage(userId, conn, payload)
				if err != nil {
					tx.Rollback(context.Background())
					return err
				}

				tx.Commit(context.Background())
			} else if payload.Type == "ping" {
				result, _ := json.Marshal(WsMsg{Type: "pong"})

				MB.sendMessages(result, []string{userId})
			}

			return nil
		})

		if err != nil {
			fmt.Println(err)
		}
	}
}

func handleCreateMessage(userId string, conn *pgxpool.Conn, payload *WsMsg) error {
	chatService := services.NewChatService(conn)
	mediaIds := []string{}
	for _, mediaId := range payload.Payload["mediaIds"].([]interface{}) {
		mediaIds = append(mediaIds, mediaId.(string))
	}

	message, err := chatService.AddMessage(
		payload.Payload["chatId"].(string),
		payload.Payload["message"].(string),
		mediaIds,
		userId,
	)

	if err != nil {
		return err
	}

	result, _ := json.Marshal(WsMsg{
		Type: "msg",
		Payload: map[string]interface{}{
			"Id":          message.Id,
			"UserId":      message.UserId,
			"ChatId":      message.ChatId,
			"Message":     message.Message,
			"MediasCount": message.MediasCount,
			"Medias":      message.Medias,
			"CreatedAt":   message.CreatedAt,
		},
	})

	memberIds, err := services.GetChatMembersForChat(conn, payload.Payload["chatId"].(string))
	if err != nil {
		fmt.Println(err)
		return err
	}

	tcpPingNodes(result, memberIds)
	// MB.sendMessages(result, memberIds)

	return nil
}

func tcpPingNodes(message []byte, memberIds []string) {
	// fmt.Println("memberIds:", memberIds)
	for _, memberId := range memberIds {
		// fmt.Println("REDIS_CLI.GetUserWsIps(memberId):", memberId, ", ", REDIS_CLI.GetUserWsIps(memberId))
		for _, ip := range REDIS_CLI.GetUserWsIps(memberId) {
			localIp, _ := IsLocalIp(ip)
			// fmt.Println("===>", memberId, ", ", ip, ", ", localIp)
			if localIp {
				WS_HUB.sendMessageChan <- WsMessageData{
					userId:  memberId,
					message: message,
				}
			} else {
				tcpSendMessage(TcpData{
					UserId:  memberId,
					Message: message,
				}, ip)
			}
		}
	}
}

func addUserWsIP(userId string) error {
	ip, err := GetIp()
	if err != nil {
		return err
	}

	REDIS_CLI.UserIpNodeAdd(userId, ip)

	return nil
}

func (h *WSClientHub) sendMessageToWsClient(message []byte, userId string) {
	for i := 0; i < len(h.clients[userId]); i += 1 {
		if h.clients[userId][i] == nil {
			continue
		}

		h.clients[userId][i].WriteMessage(websocket.TextMessage, message)
	}
}

func (h *WSClientHub) addWSClient(userId string, ws *websocket.Conn) {
	newLen := len(h.clients[userId]) + 1
	conns := make([]*websocket.Conn, newLen)
	for i := 0; i < newLen-1; i += 1 {
		conns[i] = h.clients[userId][i]
	}
	conns[newLen-1] = ws
	h.clients[userId] = conns
}

func (h *WSClientHub) removeWSClient(userId string, ws *websocket.Conn) {
	index := 0
	for index < len(h.clients[userId]) {
		if h.clients[userId][index] == ws {
			break
		}
		index++
	}

	newClients := []*websocket.Conn{}
	for _, tws := range h.clients[userId] {
		if tws == ws {
			continue
		}

		newClients = append(newClients, tws)
	}

	h.clients[userId] = newClients
}
