package services

import (
	"context"

	"chat-app.vi/models"
	"chat-app.vi/utils"
	"github.com/jackc/pgx/v4/pgxpool"
	uuid "github.com/satori/go.uuid"
)

type ChatService struct {
	conn *pgxpool.Conn
}

func NewChatService(conn *pgxpool.Conn) *ChatService {
	return &ChatService{conn: conn}
}

func (cs *ChatService) GetChats(userId string) ([]*models.ChatMember, error) {
	rows, err := cs.conn.Query(context.Background(), "select_chats", userId)
	if err != nil {
		return []*models.ChatMember{}, err
	}

	var chatIds []string
	for rows.Next() {
		var chatId string
		rows.Scan(&chatId)
		chatIds = append(chatIds, chatId)
	}

	rows.Close()

	rows, err = cs.conn.Query(context.Background(), "select_chat_members_chat_id", chatIds)
	if err != nil {
		return []*models.ChatMember{}, err
	}

	chatMembersMap := map[string][]string{}
	for rows.Next() {
		var chatId string
		var userId string
		rows.Scan(&chatId, &userId)
		chatMembersMap[chatId] = append(chatMembersMap[chatId], userId)
	}

	rows.Close()

	memberIds := []string{}
	for key := range chatMembersMap {
		if chatMembersMap[key][0] == userId {
			memberIds = append(memberIds, chatMembersMap[key][1])
		} else {
			memberIds = append(memberIds, chatMembersMap[key][0])
		}
	}

	rows.Close()

	rows, err = cs.conn.Query(context.Background(), "select_chat_users", memberIds)
	if err != nil {
		return []*models.ChatMember{}, err
	}

	chatUsers := []models.ChatMember{}
	for rows.Next() {
		var user models.ChatMember
		var imageId *string
		rows.Scan(&user.Id, &user.Name, &user.Username, &imageId, &user.CreatedAt)
		if imageId == nil {
			user.Image = GetDefaultMediaPath()
		} else {
			user.Image = GetMediaPath(*imageId)
		}
		chatUsers = append(chatUsers, user)
	}

	rows.Close()

	result := []*models.ChatMember{}
	for chatId := range chatMembersMap {
		var targetUserId string
		if chatMembersMap[chatId][0] == userId {
			targetUserId = chatMembersMap[chatId][1]
		} else {
			targetUserId = chatMembersMap[chatId][0]
		}
		member := findUser(chatUsers, targetUserId)
		if member == nil {
			continue
		}

		member.ChatId = chatId

		result = append(result, member)
	}

	targetChatIds := []string{}
	for _, member := range result {
		targetChatIds = append(targetChatIds, member.ChatId)
	}

	r, err := getLastMessageIdsForChats(cs.conn, targetChatIds)
	if err != nil {
		return []*models.ChatMember{}, err
	}
	messageIds := []string{}
	for chatId := range r {
		if r[chatId] == "" {
			continue
		}
		messageIds = append(messageIds, r[chatId])
	}
	messages, err := getMessagesByIds(cs.conn, messageIds)
	if err != nil {
		return []*models.ChatMember{}, err
	}

	for _, member := range result {
		messageId := r[member.ChatId]
		if messageId == "" {
			continue
		}
		member.LastMessage = messages[messageId]
	}

	return result, nil
}

func getLastMessageIdsForChats(conn *pgxpool.Conn, chatIds []string) (map[string]string, error) {
	rows, err := conn.Query(context.Background(), "select_chat_membersx", chatIds)
	if err != nil {
		return map[string]string{}, err
	}

	result := map[string]string{}
	for rows.Next() {
		var chatId string
		var lastMessageId string
		rows.Scan(&chatId, &lastMessageId)
		result[chatId] = lastMessageId
	}

	rows.Close()

	return result, nil
}

func getMessagesByIds(conn *pgxpool.Conn, messageIds []string) (map[string]models.ChatMessage, error) {
	if len(messageIds) == 0 {
		return map[string]models.ChatMessage{}, nil
	}
	rows, err := conn.Query(context.Background(), "select_chat_membersxx", messageIds)
	if err != nil {
		return map[string]models.ChatMessage{}, err
	}

	result := map[string]models.ChatMessage{}
	for rows.Next() {
		var msg models.ChatMessage
		rows.Scan(&msg.Id, &msg.ChatId, &msg.UserId, &msg.Message)
		result[msg.Id] = msg
	}

	rows.Close()

	return result, nil
}

func findUser(users []models.ChatMember, userId string) *models.ChatMember {
	for _, user := range users {
		if user.Id == userId {
			return &user
		}
	}

	return nil
}

func (cs *ChatService) GetChat(chatId, loggedUserId string) ([]models.ChatMember, error) {
	memberIds, err := GetChatMembersForChat(cs.conn, chatId)
	if err != nil {
		return []models.ChatMember{}, err
	}

	rows, err := cs.conn.Query(context.Background(), "select_users_by_id", memberIds)
	if err != nil {
		return []models.ChatMember{}, err
	}

	chatUsers := []models.ChatMember{}
	for rows.Next() {
		var user models.ChatMember
		var imageId *string
		rows.Scan(&user.Id, &user.Name, &user.Username, &imageId, &user.CreatedAt)
		if imageId == nil {
			user.Image = GetDefaultMediaPath()
		} else {
			user.Image = GetMediaPath(*imageId)
		}
		chatUsers = append(chatUsers, user)
	}

	rows.Close()

	result := []models.ChatMember{}
	for _, user := range chatUsers {
		if user.Id == loggedUserId {
			continue
		}

		result = append(result, user)
	}

	return result, nil
}

func GetChatMembersForChat(conn *pgxpool.Conn, chatId string) ([]string, error) {
	rows, err := conn.Query(context.Background(), "select_chat_members_by_chat", chatId)
	if err != nil {
		return []string{}, err
	}

	var memberIds []string
	for rows.Next() {
		var memberId string
		rows.Scan(&memberId)
		memberIds = append(memberIds, memberId)
	}

	rows.Close()

	return memberIds, nil
}

func (cs *ChatService) GetMessages(chatId, loggedUserId string) ([]models.ChatMessage, error) {
	rows, err := cs.conn.Query(context.Background(), "select_chat_messages_by_chat", chatId)
	if err != nil {
		return []models.ChatMessage{}, err
	}

	messages := []models.ChatMessage{}
	for rows.Next() {
		var message models.ChatMessage
		rows.Scan(&message.Id, &message.UserId, &message.Message, &message.MediasCount, &message.CreatedAt)
		messages = append(messages, message)
	}

	messageIds := []string{}
	for _, message := range messages {
		messageIds = append(messageIds, message.Id)
	}

	medias, err := getMediasForMessages(cs.conn, messageIds)
	if err != nil {
		return []models.ChatMessage{}, err
	}

	t := []*models.ChatMessage{}
	for i := 0; i < len(messages); i += 1 {
		t = append(t, &messages[i])
	}

	for _, message := range t {
		mediaPaths := []string{}
		for _, mediaId := range medias[message.Id] {
			mediaPaths = append(mediaPaths, GetMediaPath(mediaId))
		}

		message.Medias = mediaPaths
	}

	rows.Close()

	return messages, nil
}

func getMediasForMessages(conn *pgxpool.Conn, messageIds []string) (map[string][]string, error) {
	if len(messageIds) == 0 {
		return map[string][]string{}, nil
	}

	rows, err := conn.Query(context.Background(), "select_chat_message_medias", messageIds)
	if err != nil {
		return map[string][]string{}, err
	}

	result := map[string][]string{}
	for rows.Next() {
		var messageId string
		var mediaId string
		rows.Scan(&messageId, &mediaId)
		result[messageId] = append(result[messageId], mediaId)
	}

	rows.Close()

	return result, nil
}

func (cs *ChatService) AddMessage(chatId string, message string, mediaIds []string, userId string) (models.ChatMessage, error) {
	messageId := uuid.NewV4().String()
	createdAt := utils.CurrentTimeMs()
	mediasCount := len(mediaIds)

	_, err := cs.conn.Exec(
		context.Background(),
		"add_chat_message",
		messageId,
		userId,
		chatId,
		message,
		mediasCount,
		createdAt,
	)
	if err != nil {
		return models.ChatMessage{}, err
	}

	for _, mediaId := range mediaIds {
		_, err = cs.conn.Exec(
			context.Background(),
			"add_chat_message_media_rel",
			messageId,
			mediaId,
		)
		if err != nil {
			return models.ChatMessage{}, err
		}
	}

	chatMmessage := models.ChatMessage{
		Id:          messageId,
		UserId:      userId,
		ChatId:      chatId,
		Message:     message,
		MediasCount: mediasCount,
		Medias:      MediaIdsToPaths(mediaIds),
		CreatedAt:   createdAt,
	}

	_, err = cs.conn.Exec(
		context.Background(),
		"update_last_message_id_for_chat",
		messageId,
		chatId,
	)
	if err != nil {
		return models.ChatMessage{}, err
	}

	return chatMmessage, nil
}
