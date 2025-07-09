package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"socialNetworkOtus/internal/actions/getuserbyid"
	"socialNetworkOtus/internal/actions/login"
	"socialNetworkOtus/internal/actions/register"
	"socialNetworkOtus/internal/api"
	"socialNetworkOtus/internal/repository"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	UserRepo *repository.UserRepository

	JWTSecret string

	RegisterService    *register.Service
	LoginService       *login.Service
	GetUserByIDService *getuserbyid.Service
}

func NewUserHandler(userRepo *repository.UserRepository) *UserHandler {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "dev_secret"
	}
	return &UserHandler{
		UserRepo: userRepo,

		JWTSecret: secret,

		RegisterService:    register.NewService(userRepo),
		LoginService:       login.NewService(userRepo),
		GetUserByIDService: getuserbyid.NewService(userRepo),
	}
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

	userID, err := h.RegisterService.RegisterUser(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

	token, err := h.LoginService.LoginUser(context.Background(), *req.Id, *req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *UserHandler) GetUserGetId(c *gin.Context, id api.UserId) {
	user, err := h.GetUserByIDService.GetUserByID(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
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
