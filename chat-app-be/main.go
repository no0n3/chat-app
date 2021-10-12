package main

import (
	"net/http"

	"chat-app.vi/controllers"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	initWSHub()
	initMessageBroker()

	router := mux.NewRouter()

	router.HandleFunc("/ws", WsHandler)
	router.HandleFunc("/api/user/{id}", controllers.ProfileHandler).Methods("GET")
	router.HandleFunc("/api/login", controllers.LoginHandler).Methods("POST")
	router.HandleFunc("/api/change-profile-image", controllers.ChangeProfileImageHandler).Methods("POST")
	router.HandleFunc("/api/sign-up", controllers.SignupHandler).Methods("POST")
	router.HandleFunc("/api/upload-profile-image", controllers.UploadProfileImageHandler).Methods("POST")
	router.HandleFunc("/api/upload", controllers.UploadHandler).Methods("POST")
	router.HandleFunc("/api/contacts", controllers.ContactsHandler).Methods("GET")
	router.HandleFunc("/api/find", controllers.FindUsersHandler).Methods("GET")
	router.HandleFunc("/api/chat", controllers.ChatsHandler).Methods("GET")
	router.HandleFunc("/api/chat/{id}", controllers.ChatHandler).Methods("GET")
	router.HandleFunc("/api/chat/{id}/messages", controllers.ChatMessageHandler).Methods("GET")
	router.HandleFunc("/api/media/{id}", controllers.GetMediaHandler).Methods("GET")
	router.HandleFunc("/api/user/{id}/add-contact", controllers.AddContactHandler).Methods("POST")
	router.HandleFunc("/api/user/{id}/remove-contact", controllers.RemoveContactHandler).Methods("POST")

	http.Handle("/", router)
	http.ListenAndServe(":1122", handlers.CORS(
		handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization", "X-Auth-Token"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}),
		handlers.AllowedOrigins([]string{"*"}),
	)(router))
}
