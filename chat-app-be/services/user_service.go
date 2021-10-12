package services

import (
	"context"

	errors "chat-app.vi/errors"
	"chat-app.vi/models"
	"chat-app.vi/utils"
	"github.com/jackc/pgx/v4/pgxpool"
	uuid "github.com/satori/go.uuid"
)

type UserService struct {
	conn *pgxpool.Conn
}

func NewUserService(conn *pgxpool.Conn) *UserService {
	return &UserService{conn: conn}
}

func (us *UserService) GetUsers(loggedUserId string) ([]models.UserItem, error) {
	var users []models.UserItem

	rows, err := us.conn.Query(
		context.Background(),
		"SELECT id, name, username, description, profile_image_id, created_at FROM users",
	)
	if err != nil {
		return []models.UserItem{}, err
	}

	for rows.Next() {
		var user models.UserItem
		var imageId *string
		var description *string
		err := rows.Scan(&user.Id, &user.Name, &user.Username, &description, &imageId, &user.CreatedAt)
		if imageId == nil {
			user.Image = GetDefaultMediaPath()
		} else {
			user.Image = GetMediaPath(*imageId)
		}

		if description == nil {
			user.Description = ""
		} else {
			user.Description = *description
		}

		if err != nil {
			return []models.UserItem{}, err
		}

		users = append(users, user)
	}

	rows.Close()

	contacts, err := us.getContactsForUser(loggedUserId)
	if err != nil {
		return []models.UserItem{}, err
	}

	for i := 0; i < len(users); i += 1 {
		users[i].IsContact = contacts[users[i].Id]
	}

	return users, nil
}

func (us *UserService) GetContacts(userId string) ([]models.UserItem, error) {
	rows, err := us.conn.Query(context.Background(), "select_user_contacts", userId)
	if err != nil {
		return []models.UserItem{}, err
	}

	var users = []models.UserItem{}
	for rows.Next() {
		var user models.UserItem
		var imageId *string
		var description *string
		rows.Scan(&user.Id, &user.Name, &user.Username, &description, &imageId, &user.CreatedAt)
		if imageId == nil {
			user.Image = GetDefaultMediaPath()
		} else {
			user.Image = GetMediaPath(*imageId)
		}

		if description == nil {
			user.Description = ""
		} else {
			user.Description = *description
		}

		user.IsContact = true

		users = append(users, user)
	}
	rows.Close()

	return users, nil
}

func (us *UserService) GetUserById(userId, loggedUserId string) (models.UserProfile, error) {
	var user models.UserProfile
	r, err := us.conn.Query(context.Background(), "select_profile", userId)
	if err != nil {
		return models.UserProfile{}, err
	}

	if !r.Next() {
		return models.UserProfile{}, &errors.NotFountError{}
	}

	var imageId *string
	err = r.Scan(&user.Id, &user.Name, &user.Username, &imageId, &user.Description, &user.CreatedAt)
	if err != nil {
		return models.UserProfile{}, err
	}
	if imageId == nil {
		user.Image = GetDefaultMediaPath()
	} else {
		user.Image = GetMediaPath(*imageId)
	}
	r.Close()

	contactResult, err := us.conn.Query(context.Background(), "is_contact", loggedUserId, userId)
	if err != nil {
		return models.UserProfile{}, err
	}

	if contactResult.Next() {
		var s1 string
		var s2 string
		var s3 int64
		contactResult.Scan(&s1, &s2, &s3)
		user.IsContact = true
	}

	contactResult.Close()

	return user, nil
}

func (us *UserService) getContactsForUser(loggedUserId string) (map[string]bool, error) {
	rows, err := us.conn.Query(context.Background(), "select_contacts", loggedUserId)

	if err != nil {
		return map[string]bool{}, err
	}

	result := map[string]bool{}

	for rows.Next() {
		var adderId string
		var addedId string
		err := rows.Scan(&adderId, &addedId)
		if err != nil {
			return map[string]bool{}, err
		}

		var key string
		if loggedUserId == adderId {
			key = addedId
		} else {
			key = adderId
		}

		result[key] = true
	}

	rows.Close()

	return result, nil
}

func (us *UserService) AddContact(adderUserId, addedUserId string) error {
	if addedUserId == adderUserId {
		return nil
	}

	_, err := us.conn.Exec(context.Background(), "add_contact", adderUserId, addedUserId, utils.CurrentTimeMs())
	if err != nil {
		return err
	}

	rows, err := us.conn.Query(context.Background(), "select_common_chat_id", adderUserId, addedUserId)
	if err != nil {
		return err
	}

	if rows.Next() {
		return nil
	}

	chatId := uuid.NewV4()
	_, err = us.conn.Exec(context.Background(), "create_chat", chatId, utils.CurrentTimeMs())
	if err != nil {
		return err
	}

	_, err = us.conn.Exec(context.Background(), "create_chat_member", chatId, adderUserId, utils.CurrentTimeMs())
	if err != nil {
		return err
	}
	_, err = us.conn.Exec(context.Background(), "create_chat_member", chatId, addedUserId, utils.CurrentTimeMs())
	if err != nil {
		return err
	}

	return nil
}

func (us *UserService) RemoveContact(adderUserId, addedUserId string) error {
	_, err := us.conn.Exec(
		context.Background(),
		"remove_contact",
		adderUserId,
		addedUserId,
	)
	if err != nil {
		return err
	}

	return nil
}

func (us *UserService) ChangeProfileImage(userId, mediaId string) error {
	_, err := us.conn.Exec(context.Background(), "change_profile_image", mediaId, userId)
	if err != nil {
		return err
	}

	return nil
}
