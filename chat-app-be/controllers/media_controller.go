package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	errors "chat-app.vi/errors"
	db "chat-app.vi/rdb"
	services "chat-app.vi/services"
	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
)

func GetMediaHandler(w http.ResponseWriter, r *http.Request) {
	mediaId := mux.Vars(r)["id"]

	err := db.POOL.HandleFunc(func(conn *pgxpool.Conn) error {
		mediaService := services.NewMediaService(conn)
		media, err := mediaService.GetMediaMetadata(mediaId)
		if err != nil {
			return err
		}
		if media == nil {
			return &errors.NotFountError{Err: err}
		}
		f, err := os.Open("./uploads/" + media.Id + "." + strings.Split(media.MimeType, "/")[1])
		if err != nil {
			return err
		}

		for {
			bytes := make([]byte, 1012)
			r, err := f.Read(bytes)
			if r == 0 || err != nil {
				break
			}
			w.Write(bytes)
		}
		f.Close()

		return nil
	})

	SendErrorRespIfHas(w, err)
}

func UploadProfileImageHandler(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("image")
	if err != nil {
		fmt.Println(err)
		sendResp(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	err = db.POOL.HandleFunc(func(conn *pgxpool.Conn) error {
		loggedUserId, err := getLoggedUserId(conn, r)
		if err != nil {
			return err
		}
		tx, err := conn.Begin(context.Background())
		if err != nil {
			return err
		}

		userService := services.NewUserService(conn)
		mediaService := services.NewMediaService(conn)
		mediaId, err := mediaService.UploadMedia(file)
		if err != nil {
			tx.Rollback(context.Background())
			return err
		}
		err = userService.ChangeProfileImage(loggedUserId, mediaId)
		if err != nil {
			tx.Rollback(context.Background())
			return err
		}

		result, err := json.Marshal(map[string]string{
			"mediaPath": services.GetMediaPath(mediaId),
		})
		if err != nil {
			tx.Rollback(context.Background())
			return err
		}

		tx.Commit(context.Background())
		sendResp(w, string(result), http.StatusOK)
		fmt.Println(mediaId)

		return nil
	})

	SendErrorRespIfHas(w, err)
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	file, _, err := r.FormFile("image")
	if err != nil {
		fmt.Println(err)
		sendResp(w, "Inernal server error.", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	if _, err := os.Stat("./uploads"); os.IsNotExist(err) {
		err = os.Mkdir("./uploads", 0777)
		if err != nil {
			fmt.Println(err)
			sendResp(w, "Inernal server error.", http.StatusInternalServerError)
			return
		}
	}

	mimeType, err := services.DetectType(file)
	if err != nil || strings.Split(mimeType, "/")[0] != "image" {
		fmt.Println(err)
		sendResp(w, "Only images allowed.", http.StatusBadRequest)
		return
	}

	err = db.POOL.HandleFunc(func(conn *pgxpool.Conn) error {
		tx, err := conn.Begin(context.Background())
		if err != nil {
			return err
		}
		mediaService := services.NewMediaService(conn)
		mediaId, err := mediaService.CreateMediaMetadata(mimeType)
		if err != nil {
			tx.Rollback(context.Background())
			return err
		}

		f, err := os.OpenFile("./uploads/"+mediaId+"."+strings.Split(mimeType, "/")[1], os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			tx.Rollback(context.Background())
			return err
		}
		defer f.Close()

		_, err = io.Copy(f, file)
		if err != nil {
			tx.Rollback(context.Background())
			return err
		}

		result, err := json.Marshal(map[string]string{
			"mediaId":   mediaId,
			"mediaPath": services.GetMediaPath(mediaId),
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
