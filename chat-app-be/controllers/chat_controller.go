package controllers

import (
	"encoding/json"
	"net/http"

	db "chat-app.vi/rdb"
	services "chat-app.vi/services"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
)

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	err := db.POOL.HandleFunc(func(conn *pgxpool.Conn) error {
		params := mux.Vars(r)
		chatId := params["id"]
		chatService := services.NewChatService(conn)
		loggedUserId, err := getLoggedUserId(conn, r)
		if err != nil {
			return err
		}

		user, err := chatService.GetChat(chatId, loggedUserId)
		if err != nil {
			return err
		}

		result, err := json.Marshal(user)
		if err != nil {
			return err
		}

		sendResp(w, string(result), http.StatusOK)

		return nil
	})

	SendErrorRespIfHas(w, err)
}

func ChatsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	err := db.POOL.HandleFunc(func(conn *pgxpool.Conn) error {
		loggedUserId, err := getLoggedUserId(conn, r)
		if err != nil {
			return err
		}
		chatService := services.NewChatService(conn)
		users, err := chatService.GetChats(loggedUserId)
		if err != nil {
			return err
		}

		result, err := json.Marshal(users)
		if err != nil {
			return err
		}

		sendResp(w, string(result), http.StatusOK)

		return nil
	})

	SendErrorRespIfHas(w, err)
}

func ChatMessageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	err := db.POOL.HandleFunc(func(conn *pgxpool.Conn) error {
		params := mux.Vars(r)
		chatId := params["id"]
		chatService := services.NewChatService(conn)
		loggedUserId, err := getLoggedUserId(conn, r)
		if err != nil {
			return err
		}
		messages, err := chatService.GetMessages(chatId, loggedUserId)
		if err != nil {
			return err
		}

		result, err := json.Marshal(messages)
		if err != nil {
			return err
		}

		sendResp(w, string(result), http.StatusOK)

		return nil
	})

	SendErrorRespIfHas(w, err)
}
