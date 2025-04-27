package models

type User struct {
	ID       int
	username string
	is_admin bool
}

func NewUser() *User {
	return &User{}
}

type Article struct {
	ID         int
	title      string
	is_deleted bool
}

func NewArticle() *Article {
	return &Article{}
}

type ArticleRevision struct {
	ID         int
	user_id    int
	user_ip    string
	is_deleted bool
	title      string
	content    string
}

func NewArticleRevision() *ArticleRevision {
	return &ArticleRevision{}
}
