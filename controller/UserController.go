package controller

import (
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/gin-gonic/gin"
	"github.com/vks/go-web-clean-architecture/domain"
	"github.com/vks/go-web-clean-architecture/service"
)

type userController struct {
	userService domain.UserService
	// fbService   service.FirebaseService
}
type UserController interface {
	GetUsers(c *gin.Context)
	AddUser(c *gin.Context)
}

//NewUserController: constructor, dependency injection from user service and firebase service
func NewUserController(s domain.UserService) UserController {
	return &userController{
		userService: s,
		// fbService:   f,
	}
}
func (u *userController) GetUsers(c *gin.Context) {
	users, err := u.userService.FindAll()
	if err != nil {
		sentry.CaptureException(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error while getting users"})
		return
	}
	c.JSON(http.StatusOK, users)
}
func (u *userController) AddUser(c *gin.Context) {
	var user domain.User
	if err := c.ShouldBindJSON(&user); err != nil {
		sentry.CaptureException(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err1 := (u.userService.Validate(&user)); err1 != nil {
		sentry.CaptureException(err1)
		c.JSON(http.StatusBadRequest, gin.H{"error": err1.Error()})
		return
	}
	if ageValidation := (u.userService.ValidateAge(&user)); ageValidation != true {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid DOB"})
		return
	}
	uid, err := u.fbService.CreateUser(user.Email, user.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Couldnot create user in firebase"})
		return
	}
	user.ID = uid
	u.userService.Create(&user)
	c.JSON(http.StatusOK, user)
}
