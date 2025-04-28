package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"golang.org/x/crypto/bcrypt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/gorilla/sessions"

	"github.com/Antarktidov/GopnikWiki/models"
)

var store = sessions.NewCookieStore([]byte("секретный-ключ"))

var db *gorm.DB

func main() {
	initDB()

	router := gin.Default()

	router.GET("/api/articles", getArticles)
	router.GET("/api/article/:id", getArticle)
	router.GET("/api/debug/revisions", getRevisions)
	router.POST("/api/create_article", createArticle)
	router.POST("/api/edit_article/:id", editArticle)

	router.POST("/api/register", register)
	router.POST("/api/login", login)

	router.DELETE("/api/delete_article/:id", deleteArticle)

	router.Run("localhost:8080")
}

func initDB() {
	dsn := "host=localhost user=postgres password=kuzura dbname=gopnik_wik_dev port=5432 sslmode=disable"
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}

	log.Println("Успешное подключение к базе данных")
}

func getArticles(c *gin.Context) {
	var articles []models.Article
	db.Where("is_deleted = false").Find(&articles)

	c.IndentedJSON(http.StatusOK, &articles)
}

func getArticle(c *gin.Context) {
	var article models.Article
	var revision models.ArticleRevision
	id := c.Param("id")

	// Найти статью по ID
	if err := db.Where("is_deleted = false").First(&article, id).Error; err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "article not found"})
		return
	}

	// Найти последнюю ревизию статьи
	if err := db.Where("is_deleted = false").Where("article_id = ?", id).Last(&revision).Error; err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "revision not found"})
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{
		"article":  article,
		"revision": revision,
	})
}

func editArticle(c *gin.Context) {
	id := c.Param("id")

	var article models.Article
	var article_revision models.ArticleRevision

	// Найти статью по ID
	if err := db.Where("is_deleted = false").First(&article, id).Error; err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "article not found"})
		return
	}

	if err := c.ShouldBindJSON(&article_revision); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session, _ := store.Get(c.Request, "session-name")

	userIDValue, ok := session.Values["user_id"]
	if !ok {
		article_revision.UserID = 0
	} else {
		if userID, ok := userIDValue.(int); ok {
			article_revision.UserID = userID
		} else {
			article_revision.UserID = 0
		}
	}

	article_revision.ArticleID = article.ID

	//article_revision.UserID = 0
	article_revision.UserIP = "unknown"
	article_revision.IsDeleted = false
	article_revision.CreatedAt = time.Now()
	article_revision.UpdatedAt = time.Now()

	if err := db.Create(&article_revision).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "could not create article revision"})
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"message": "article and edited successfully"})
}

func createArticle(c *gin.Context) {

	var article models.Article
	var article_revision models.ArticleRevision

	if err := c.ShouldBindJSON(&article_revision); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	session, _ := store.Get(c.Request, "session-name")

	id, ok := session.Values["user_id"]
	if !ok {
		article_revision.UserID = 0
	} else {
		if userID, ok := id.(int); ok {
			article_revision.UserID = userID
		} else {
			article_revision.UserID = 0
		}
	}

	article.Title = article_revision.Title
	article.IsDeleted = false
	article.CreatedAt = time.Now()

	//fmt.Println(&article)

	if err := db.Create(&article).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "could not create article"})
		return
	}

	article_revision.ArticleID = article.ID

	//article_revision.UserID = 0
	article_revision.UserIP = "unknown"
	article_revision.IsDeleted = false
	article_revision.CreatedAt = time.Now()
	article_revision.UpdatedAt = time.Now()

	if err := db.Create(&article_revision).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "could not create article revision"})
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"message": "article and revision created successfully"})
}

func deleteArticle(c *gin.Context) {
	session, _ := store.Get(c.Request, "session-name")
	var user models.User

	id, ok := session.Values["user_id"]
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Unauthorized"})
		return
	} else {
		if userID, ok := id.(int); ok {
			if err := db.Model(&models.User{}).Omit("Password").Where("id = ? ", userID).Where("is_admin = true").First(&user).Error; err != nil {
				c.JSON(http.StatusForbidden, gin.H{"message": "Forbidden"})
				return
			}

			article_id := c.Param("id")

			var article models.Article

			db.Model(&article).Where("id = ?", article_id).Update("is_deleted", true)

		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "InternalServerError"})
			return
		}
	}
}

func getRevisions(c *gin.Context) {
	var revisions []models.ArticleRevision
	db.Where("is_deleted = false").Find(&revisions)

	c.IndentedJSON(http.StatusOK, &revisions)
}

func register(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	encryptedPassword, err := encryptString(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unknown error"})
		return
	}
	user.EncryptedPassword = encryptedPassword
	user.Password = ""

	if err := db.Model(&models.User{}).Omit("Password").Create(&user).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "could not finish registration"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "successfully registered"})
}

func login(c *gin.Context) {
	var userInput models.User
	var storedUser models.User

	if err := c.ShouldBindJSON(&userInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.Model(&models.User{}).Omit("Password").Where("username = ? ", userInput.Username).First(&storedUser).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid username or password"})
		return
	}

	if !storedUser.ComparePassword((userInput.Password)) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "invalid username or password"})
		return
	}

	// Создание новой сессии
	session, err := store.Get(c.Request, "session-name")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not create session"})
		return
	}

	// Сохранение ID пользователя в сессию
	session.Values["user_id"] = storedUser.ID
	if err := session.Save(c.Request, c.Writer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "could not save session"})
		return
	}
	// Успешный вход
	c.JSON(http.StatusOK, gin.H{"message": "login successful"})
}

func encryptString(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
