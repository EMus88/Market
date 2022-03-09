package handler

import (
	"net/http"
	"os"
	"strings"

	"github.com/EMus88/Market/internal/models"

	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
)

func (h *Handler) AuthMiddleware(c *gin.Context) {
	//read header
	authHeader := strings.Split(c.GetHeader("Authorization"), " ")
	if len(authHeader) != 2 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		c.Abort()
		return
	}
	bearerToken := authHeader[1]
	//validate token
	if _, _, err := h.service.ValidateToken(bearerToken, "access"); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		c.Abort()
		return
	}
	c.Next()
}

func (h *Handler) IsAdminMiddleware(c *gin.Context) {
	//getting a token
	authHeader := strings.Split(c.GetHeader("Authorization"), " ")
	if len(authHeader) != 2 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthenticated"})
		c.Abort()
		return
	}
	bearerToken := authHeader[1]
	//getting claims from token
	id, role, err := h.service.ValidateToken(bearerToken, "access")
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "credential error"})
		c.Abort()
		return
	}
	if role != "admin" {
		c.JSON(http.StatusConflict, gin.H{"error": "credential error"})
		c.Abort()
		return
	}
	uuidID, err := uuid.FromString(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		c.Abort()
		return
	}
	//check admin in db
	roleFromDB, err := h.service.CheckUser(uuidID)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "credential error"})
		c.Abort()
		return
	}
	if roleFromDB != role {
		c.JSON(http.StatusConflict, gin.H{"error": "credential error"})
		c.Abort()
		return
	}
}

//add admin
// @Summary Add administrator
// @Tags auth
// @Descriotion create user with adminostrator credentials
// @Accept json
// @Produce json
// @Param input body models.Admin true "account info"
// @Success 200 {object} models.User
// @Failure 400 {string} json "{"error":"Not allowed request"}"
// @Failure 500 {string} json "{"error":"Internal server error"}"
// @Router /auth/admin [post]
func (h *Handler) AddAddmin(c *gin.Context) {
	var admin models.Admin
	//parse request
	if err := c.ShouldBindJSON(&admin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not allowed request"})
		return
	}
	if admin.Code != os.Getenv("ADMINCODE") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not allowed request"})
		return
	}
	var user models.User

	user.Username = admin.Username
	user.Password = admin.Password
	user.Phone = admin.Phone
	user.Role = "admin"

	//save in db
	if err := h.service.Auth.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// @Summary Registration
// @Security ApiKeyAuth
// @Tags auth
// @Descriotion registration new user
// @Accept json
// @Produce json
// @Param input body models.User true "account info"
// @Success 200 {object} models.User
// @Failure 400 {string} json "{"error":"Not allowed request"}"
// @Failure 409 {string} json "{"error":"credential error"}"
// @Failure 411 {string} json "{"error":"Not allowed lengths of data"}"
// @Failure 500 {string} json "{"error":"Internal server error"}"
// @Router /auth/signUp [post]
func (h *Handler) SignUp(c *gin.Context) {
	var user models.User
	//parse request
	if err := c.ShouldBindJSON(&user); err != nil {
		h.logger.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not allowed request"})
		return
	}
	//validation request
	if ok, _ := govalidator.ValidateStruct(user); !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not allowed request"})
		return
	}
	if (len(user.Password) < 7) || (len(user.Password) > 50) || (len(user.Username) > 50) {
		c.JSON(http.StatusLengthRequired, gin.H{"error": "Not allowed lengths of data"})
		return
	}
	user.Role = userRole
	//save in db
	if err := h.service.Auth.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// @Summary Authorizaton
// @Tags auth
// @Descriotion authorization user in system
// @Accept json
// @Produce json
// @Param input body models.User true "account info"
// @Success 200 {string} json "{"access token":"...","refresh token":"..."}"
// @Failure 400 {string} json "{"error":"Not allowed request"}"
// @Failure 401 {string} json "{"error":"User not found"}"
// @Failure 500 {string} json "{"error":"Internal server error"}"
// @Router /auth/signIn [post]
func (h *Handler) SignIn(c *gin.Context) {
	var user models.User
	//parse request
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not allowed request"})
		return
	}
	//check user in db
	user.Password = h.service.Auth.HashingPassword(user.Password)
	id, role, err := h.service.Repository.GetUser(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		return
	}
	//create tokens
	t, rt, err := h.service.Auth.GenerateTokenPair(id, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"access token": t, "refresh token": rt})
}

// @Summary Update tokens
// @Tags auth
// @Descriotion refresing tokens
// @Accept json
// @Produce json
// @Param input body models.UpdateRequest true "refresh token"
// @Success 200 {string} json "{"access token":"...","refresh token":"..."}"
// @Failure 400 {string} json "{"error":"Not allowed request"}"
// @Failure 401 {string} json "{"error":"Not valid refresh token"}"
// @Failure 500 {string} json "{"error":"Internal server error"}"
// @Router /auth/update [post]
func (h *Handler) TokenRefreshing(c *gin.Context) {
	var request models.UpdateRequest
	//read refresh token
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Not allowed request"})
		return
	}
	//validate token
	id, role, err := h.service.ValidateToken(request.RefreshToken, "refresh")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Not valid refresh token"})
		return
	}
	//if validate is ok -> create new tokens
	t, rt, err := h.service.Auth.GenerateTokenPair(id, role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	//sent response
	c.JSON(http.StatusOK, gin.H{"access token": t, "refresh token": rt})

}
