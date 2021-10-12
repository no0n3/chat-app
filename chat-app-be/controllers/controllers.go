package controllers

import (
	"context"
	"encoding/json"
	"net/http"

	"chat-app.vi/models"
	db "chat-app.vi/rdb"
	services "chat-app.vi/services"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
)

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	err := db.POOL.HandleFunc(func(conn *pgxpool.Conn) error {
		params := mux.Vars(r)
		userId := params["id"]
		userService := services.NewUserService(conn)
		loggedUserId, err := getLoggedUserId(conn, r)
		if err != nil {
			return err
		}
		user, err := userService.GetUserById(userId, loggedUserId)
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

func ChangeProfileImageHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	err := db.POOL.HandleFunc(func(conn *pgxpool.Conn) error {
		var payload models.ChangeProfileImagePayload
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			return err
		}

		loggedUserId, err := getLoggedUserId(conn, r)
		if err != nil {
			return err
		}
		tx, err := conn.Begin(context.Background())
		if err != nil {
			return err
		}

		userService := services.NewUserService(conn)
		err = userService.ChangeProfileImage(loggedUserId, payload.MediaId)

		if err != nil {
			tx.Rollback(context.Background())
			return err
		}

		tx.Commit(context.Background())
		sendResp(w, "", http.StatusNoContent)

		return nil
	})

	SendErrorRespIfHas(w, err)
}

func FindUsersHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	err := db.POOL.HandleFunc(func(conn *pgxpool.Conn) error {
		loggedUserId, err := getLoggedUserId(conn, r)
		if err != nil {
			return err
		}

		userService := services.NewUserService(conn)
		users, err := userService.GetUsers(loggedUserId)
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
