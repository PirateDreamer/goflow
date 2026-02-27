package controller

import (
	"goerp-api/internal/application/service"
	"goerp-api/internal/domain/derrors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	userSvc *service.UserService
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type SendCodeRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type LoginEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required,len=6"`
}

func NewUserController(userSvc *service.UserService) *UserController {
	return &UserController{userSvc: userSvc}
}

// Register godoc
// @Summary Register a new user
// @Description register by username, email and password
// @Tags users
// @Accept  json
// @Produce  json
// @Param user body RegisterRequest true "User registration info"
// @Success 201 {object} entity.User
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/register [post]
func (ctrl *UserController) Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := ctrl.userSvc.Register(c.Request.Context(), req.Username, req.Email, req.Password)
	if err != nil {
		ctrl.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, user)
}

// Login godoc
// @Summary Login by username
// @Description login by username and password
// @Tags users
// @Accept  json
// @Produce  json
// @Param login body LoginRequest true "Login credentials"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 412 {object} map[string]string
// @Router /users/login [post]
func (ctrl *UserController) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := ctrl.userSvc.Login(c.Request.Context(), req.Username, req.Password)
	if err != nil {
		ctrl.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "login success",
		"user":    user,
	})
}

// GetUser godoc
// @Summary Get user by ID
// @Description get user detail by id
// @Tags users
// @Accept  json
// @Produce  json
// @Param id path int true "User ID"
// @Success 200 {object} entity.User
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /users/{id} [get]
func (ctrl *UserController) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	user, err := ctrl.userSvc.GetUser(c.Request.Context(), uint(id))
	if err != nil {
		ctrl.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

// SendEmailCode godoc
// @Summary Send verification code to email
// @Description send 6-digit code to email address
// @Tags users
// @Accept  json
// @Produce  json
// @Param email body SendCodeRequest true "Email address"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/send-code [post]
func (ctrl *UserController) SendEmailCode(c *gin.Context) {
	var req SendCodeRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := ctrl.userSvc.SendEmailVerificationCode(c.Request.Context(), req.Email); err != nil {
		ctrl.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "verification code sent"})
}

// LoginByEmail godoc
// @Summary Login by email verification code
// @Description login using email and 6-digit code
// @Tags users
// @Accept  json
// @Produce  json
// @Param login body LoginEmailRequest true "Email and code"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 411 {object} map[string]string
// @Router /users/login-email [post]
func (ctrl *UserController) LoginByEmail(c *gin.Context) {
	var req LoginEmailRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := ctrl.userSvc.LoginByEmailCode(c.Request.Context(), req.Email, req.Code)
	if err != nil {
		ctrl.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "login success",
		"user":    user,
	})
}

func (ctrl *UserController) handleError(c *gin.Context, err error) {
	dErr := derrors.FromError(err)
	status := http.StatusInternalServerError

	switch dErr.Code {
	case derrors.ErrInvalidParam.Code:
		status = http.StatusBadRequest
	case derrors.ErrUserNotFound.Code:
		status = http.StatusNotFound
	case derrors.ErrInvalidCredentials.Code, derrors.ErrVerificationExpired.Code, derrors.ErrInvalidVerification.Code:
		status = http.StatusUnauthorized
	}

	c.JSON(status, dErr)
}
