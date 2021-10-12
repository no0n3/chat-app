package services

import (
	"context"

	"chat-app.vi/models"
	utils "chat-app.vi/utils"
	"github.com/jackc/pgx/v4/pgxpool"
	uuid "github.com/satori/go.uuid"
)

type AuthService struct {
	conn *pgxpool.Conn
}

func NewAuthService(conn *pgxpool.Conn) *AuthService {
	return &AuthService{conn: conn}
}

func (as *AuthService) GetLoggedUserId(token string) (string, error) {
	rows, err := as.conn.Query(context.Background(), "select_logged_user_id", token)
	if err != nil {
		return "", err
	}

	if !rows.Next() {
		return "", nil
	}

	var userId string
	rows.Scan(&userId)

	rows.Close()

	return userId, nil
}

func (as *AuthService) createAuthToken(userId string) (string, string, error) {
	token := uuid.NewV4().String()
	_, err := as.conn.Exec(
		context.Background(),
		"create_token",
		userId,
		token,
		utils.CurrentTimeMs(),
	)
	if err != nil {
		return "", "", err
	}

	return userId, token, nil
}

func (as *AuthService) Login(email, password string) (string, string, error) {
	rows, err := as.conn.Query(context.Background(), "select_login_info", email)
	if err != nil {
		return "", "", err
	}

	if !rows.Next() {
		return "", "", nil
	}

	var userId string
	var passwordHash string
	rows.Scan(&userId, &passwordHash)

	rows.Close()

	if !utils.CheckPasswordHash(password, passwordHash) {
		return "", "", nil
	}

	return as.createAuthToken(userId)
}

func (as *AuthService) SignUp(payload models.SignupPayload) (string, string, error) {
	userId := uuid.NewV4().String()
	hashedPassword, err := utils.HashPassword(payload.Password)
	if err != nil {
		return "", "", err
	}

	_, err = as.conn.Exec(
		context.Background(),
		"create_user",
		userId,
		payload.Name,
		payload.Username,
		payload.Email,
		hashedPassword,
		utils.CurrentTimeMs(),
	)
	if err != nil {
		return "", "", err
	}

	return as.createAuthToken(userId)
}
