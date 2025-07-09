package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"socialNetworkOtus/internal/api"
	"socialNetworkOtus/internal/repository"
	"socialNetworkOtus/internal/utils"

	"github.com/doug-martin/goqu/v9"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type UserHandler struct {
	UserRepo  *repository.UserRepository
	JWTSecret string
}

func NewUserHandler(userRepo *repository.UserRepository) *UserHandler {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev_secret"
	}
	return &UserHandler{UserRepo: userRepo, JWTSecret: secret}
}

func (h *UserHandler) PostUserRegister(c *gin.Context) {
	var req api.PostUserRegisterJSONBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}
	if req.Password == nil || req.FirstName == nil || req.SecondName == nil || req.Birthdate == nil || req.City == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing required fields"})
		return
	}
	passwordHash, err := utils.HashPassword(*req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password: " + err.Error()})
		return
	}
	userID, err := h.UserRepo.CreateUser(context.Background(), &req, passwordHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user_id": userID})
}

func (h *UserHandler) PostLogin(c *gin.Context) {
	fmt.Println("Lets start")
	var req api.PostLoginJSONBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
		return
	}
	if req.Id == nil || req.Password == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing id or password"})
		return
	}
	user, err := h.UserRepo.GetUserByID(context.Background(), *req.Id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found: " + err.Error()})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	var passwordHash string
	db := h.UserRepo.DB()
	ds := db.From("users").Select("password_hash").Where(goqu.Ex{"id": *req.Id})
	found, err := ds.ScanValContext(context.Background(), &passwordHash)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get password hash: " + err.Error()})
		return
	}
	if !found {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	if !utils.CheckPasswordHash(*req.Password, passwordHash) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid password"})
		return
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": *req.Id,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString([]byte(h.JWTSecret))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func (h *UserHandler) GetUserGetId(c *gin.Context, id api.UserId) {
	user, err := h.UserRepo.GetUserByID(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found: " + err.Error()})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetDialogUserIdList(c *gin.Context, userId api.UserId) {
	c.JSON(501, gin.H{"error": "not implemented"})
}
func (h *UserHandler) PostDialogUserIdSend(c *gin.Context, userId api.UserId) {
	c.JSON(501, gin.H{"error": "not implemented"})
}
func (h *UserHandler) PutFriendDeleteUserId(c *gin.Context, userId api.UserId) {
	c.JSON(501, gin.H{"error": "not implemented"})
}
func (h *UserHandler) PutFriendSetUserId(c *gin.Context, userId api.UserId) {
	c.JSON(501, gin.H{"error": "not implemented"})
}
func (h *UserHandler) PostPostCreate(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}
func (h *UserHandler) PutPostDeleteId(c *gin.Context, id api.PostId) {
	c.JSON(501, gin.H{"error": "not implemented"})
}
func (h *UserHandler) GetPostFeed(c *gin.Context, params api.GetPostFeedParams) {
	c.JSON(501, gin.H{"error": "not implemented"})
}
func (h *UserHandler) GetPostGetId(c *gin.Context, id api.PostId) {
	c.JSON(501, gin.H{"error": "not implemented"})
}
func (h *UserHandler) PutPostUpdate(c *gin.Context) {
	c.JSON(501, gin.H{"error": "not implemented"})
}
func (h *UserHandler) GetUserSearch(c *gin.Context, params api.GetUserSearchParams) {
	c.JSON(501, gin.H{"error": "not implemented"})
}
