package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID                int    `json:"id"`
	Username          string `json:"username"`
	IsAdmin           bool   `json:"-"`
	EncryptedPassword string `json:"-"`
	Password          string `json:"password,omitempty"`
}

func NewUser() *User {
	return &User{}
}

func (u *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(password)) == nil
}

type Article struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	IsDeleted bool      `json:"-"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func NewArticle() *Article {
	return &Article{}
}

type ArticleRevision struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	ArticleID int       `json:"article_id"`
	UserIP    string    `json:"user_ip"`
	IsDeleted bool      `json:"-"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func NewArticleRevision() *ArticleRevision {
	return &ArticleRevision{}
}
