package models

type MediaMetadata struct {
	Id        string
	MimeType  string
	CreatedAt int64
}

type ProfileContact struct {
	Id       string
	Name     string
	IsOnline bool
}

type UserItem struct {
	Id          string
	Name        string
	Username    string
	Image       string
	Description string
	IsContact   bool
	CreatedAt   int64
}

type ChatMessage struct {
	Id          string
	UserId      string
	ChatId      string
	Message     string
	MediasCount int
	Medias      []string
	CreatedAt   int64
}

type ChatMember struct {
	Id          string
	Name        string
	Username    string
	Image       string
	ChatId      string
	LastMessage ChatMessage
	CreatedAt   int64
}

type UserProfile struct {
	Id          string
	Name        string
	Username    string
	Description string
	Image       string
	IsContact   bool
	CreatedAt   int64
}
