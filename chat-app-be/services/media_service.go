package services

import (
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"chat-app.vi/models"
	"chat-app.vi/utils"
	"github.com/jackc/pgx/v4/pgxpool"
	uuid "github.com/satori/go.uuid"
)

type MediaService struct {
	conn *pgxpool.Conn
}

func NewMediaService(conn *pgxpool.Conn) *MediaService {
	return &MediaService{conn: conn}
}

func (ms *MediaService) GetMediaMetadata(mediaId string) (*models.MediaMetadata, error) {
	rows, err := ms.conn.Query(
		context.Background(),
		"get_media_metadata",
		mediaId,
	)
	if err != nil {
		return nil, err
	}
	if !rows.Next() {
		return nil, nil
	}
	var result models.MediaMetadata
	rows.Scan(&result.Id, &result.MimeType, &result.CreatedAt)

	return &result, nil
}

func (ms *MediaService) CreateMediaMetadata(mimeType string) (string, error) {
	mediaId := uuid.NewV4().String()
	_, err := ms.conn.Exec(
		context.Background(),
		"create_media_metadata",
		mediaId,
		mimeType,
		utils.CurrentTimeMs(),
	)
	if err != nil {
		return "", err
	}

	return mediaId, nil
}

func (ms *MediaService) UploadMedia(file multipart.File) (string, error) {
	if _, err := os.Stat("./uploads"); os.IsNotExist(err) {
		if err != nil {
			return "", err
		}
		err = os.Mkdir("./uploads", 0777)
		if err != nil {
			return "", err
		}
	}

	mimeType, err := DetectType(file)
	if err != nil || strings.Split(mimeType, "/")[0] != "image" {
		return "", err
	}

	mediaId, err := ms.CreateMediaMetadata(mimeType)
	if err != nil {
		return "", err
	}

	f, err := os.OpenFile("./uploads/"+mediaId+"."+strings.Split(mimeType, "/")[1], os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return "", err
	}
	defer f.Close()
	_, err = io.Copy(f, file)

	return mediaId, err
}

func DetectType(file multipart.File) (string, error) {
	buff := make([]byte, 512)
	_, err := file.Read(buff)
	if err != nil {
		return "", err
	}

	mimeType := http.DetectContentType(buff)

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		return "", err
	}

	return mimeType, nil
}

func MediaIdsToPaths(mediaIds []string) []string {
	mediaPaths := []string{}
	for _, mediaId := range mediaIds {
		mediaPaths = append(mediaPaths, GetMediaPath(mediaId))
	}

	return mediaPaths
}

func GetMediaPath(mediaId string) string {
	return os.Getenv("CHAT_APP_ENDPOINT") + "/api/media/" + mediaId
}

func GetDefaultMediaPath() string {
	return "/default-profile-pic.jpg"
}
