package controllers

import (
	"context"
	"encoding/json"
	"net/http"

	errors "chat-app.vi/errors"
	db "chat-app.vi/rdb"
	services "chat-app.vi/services"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
)

func ContactsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	err := db.POOL.HandleFunc(func(conn *pgxpool.Conn) error {
		userService := services.NewUserService(conn)
		loggedUserId, err := getLoggedUserId(conn, r)
		if err != nil {
			sendResp(w, "", http.StatusUnauthorized)
			return err
		}

		users, err := userService.GetContacts(loggedUserId)
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

func AddContactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	err := db.POOL.HandleFunc(func(conn *pgxpool.Conn) error {
		params := mux.Vars(r)
		userId := params["id"]
		userService := services.NewUserService(conn)

		loggedUserId, err := getLoggedUserId(conn, r)
		if err != nil {
			return err
		}
		if userId == loggedUserId {
			sendResp(w, "", http.StatusBadRequest)
			return err
		}
		tx, err := conn.Begin(context.Background())
		if err != nil {
			return err
		}

		err = userService.AddContact(loggedUserId, userId)
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

func RemoveContactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	err := db.POOL.HandleFunc(func(conn *pgxpool.Conn) error {
		params := mux.Vars(r)
		userId := params["id"]
		userService := services.NewUserService(conn)

		loggedUserId, err := getLoggedUserId(conn, r)
		if err != nil {
			return err
		}
		if userId == loggedUserId {
			return &errors.BadRequestError{Err: err}
		}
		tx, err := conn.Begin(context.Background())
		if err != nil {
			return err
		}

		err = userService.RemoveContact(loggedUserId, userId)
		if err != nil {
			tx.Rollback(context.Background())
			return &errors.BadRequestError{Err: err}
		}

		tx.Commit(context.Background())
		sendResp(w, "", http.StatusNoContent)

		return nil
	})

	SendErrorRespIfHas(w, err)
}
