package controllers

import (
	"fmt"
	"net/http"

	errors "chat-app.vi/errors"
	services "chat-app.vi/services"
	"chat-app.vi/utils"
	"github.com/jackc/pgx/v4/pgxpool"
)

func SendErrorRespIfHas(w http.ResponseWriter, err error) {
	if err != nil {
		fmt.Println(err)
		switch err.(type) {
		case *errors.UnauthorizedError:
			sendResp(w, "Unauthorized", http.StatusUnauthorized)
		case *errors.NotFountError:
			sendResp(w, "Not found", http.StatusNotFound)
		case *errors.BadRequestError:
			sendResp(w, "Bad request", http.StatusBadRequest)
		default:
			sendResp(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}

func sendResp(w http.ResponseWriter, msg string, status int) {
	w.WriteHeader(status)
	fmt.Fprint(w, msg)
}

func getLoggedUserId(conn *pgxpool.Conn, r *http.Request) (string, error) {
	return getLoggedUserIdByToken(conn, utils.GetToken(r))
}

func getLoggedUserIdByToken(conn *pgxpool.Conn, token string) (string, error) {
	authService := services.NewAuthService(conn)
	loggedUserId, err := authService.GetLoggedUserId(token)
	if loggedUserId == "" || err != nil {
		return "", &errors.UnauthorizedError{}
	}

	return loggedUserId, nil
}
