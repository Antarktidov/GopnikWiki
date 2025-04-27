package models

type User struct {
	ID                int
	Username          string
	IsAdmin           bool
	EncryptedPassword string
	Password          string
}

func NewUser() *User {
	return &User{}
}

type Article struct {
	ID        int
	Title     string
	IsDeleted bool
}

func NewArticle() *Article {
	return &Article{}
}

type ArticleRevision struct {
	ID        int
	UserID    int
	ArticleID int
	UserIP    string
	IsDeleted bool
	Title     string
	Content   string
}

func NewArticleRevision() *ArticleRevision {
	return &ArticleRevision{}
}
