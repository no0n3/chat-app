package controllers

import (
	"context"
	"encoding/json"
	"net/http"

	errors "chat-app.vi/errors"
	"chat-app.vi/models"
	db "chat-app.vi/rdb"
	services "chat-app.vi/services"
	"github.com/jackc/pgx/v4/pgxpool"
)

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	err := db.POOL.HandleFunc(func(conn *pgxpool.Conn) error {
		var payload models.LoginPayload
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			return err
		}
		tx, err := conn.Begin(context.Background())
		if err != nil {
			return err
		}

		authService := services.NewAuthService(conn)
		userId, token, err := authService.Login(payload.Email, payload.Password)
		if err != nil {
			tx.Rollback(context.Background())
			return err
		}

		if token == "" {
			return &errors.UnauthorizedError{Err: err}
		}

		result, err := json.Marshal(map[string]string{
			"token":  token,
			"userId": userId,
		})
		if err != nil {
			tx.Rollback(context.Background())
			return err
		}

		tx.Commit(context.Background())
		sendResp(w, string(result), http.StatusOK)

		return nil
	})

	SendErrorRespIfHas(w, err)
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	err := db.POOL.HandleFunc(func(conn *pgxpool.Conn) error {
		var payload models.SignupPayload
		err := json.NewDecoder(r.Body).Decode(&payload)
		if err != nil {
			return err
		}
		tx, err := conn.Begin(context.Background())
		if err != nil {
			return err
		}

		authService := services.NewAuthService(conn)
		userId, token, err := authService.SignUp(payload)

		if err != nil {
			tx.Rollback(context.Background())
			return err
		}

		if token == "" {
			return &errors.UnauthorizedError{Err: err}
		}

		result, err := json.Marshal(map[string]string{
			"token":  token,
			"userId": userId,
		})
		if err != nil {
			tx.Rollback(context.Background())
			return err
		}

		tx.Commit(context.Background())
		sendResp(w, string(result), http.StatusOK)

		return nil
	})

	SendErrorRespIfHas(w, err)
}
