package handler

import (
	"JWT_auth/internal/models"
	"net/http"
	"os"
	"strings"

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
		c.JSON(http.StatusUnauthorized, gin.H{"error": "credential error"})
		c.Abort()
		return
	}
	if role != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "credential error"})
		c.Abort()
		return
	}
	uuidID, err := uuid.FromString(id)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "credential error"})
		c.Abort()
		return
	}
	//check admin in db
	roleFromDB, err := h.service.CheckUser(uuidID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "credential error"})
		c.Abort()
		return
	}
	if roleFromDB != role {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "credential error"})
		c.Abort()
		return
	}
}

//add admin
func (h *Handler) AddAddmin(c *gin.Context) {
	var admin models.Admin
	//parse request
	if err := c.ShouldBindJSON(&admin); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if admin.Code != os.Getenv("ADMINCODE") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "not valid data"})
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

//Registration
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
	if (len(user.Password) < 7) || (len(user.Password) > 20) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password length will be from 7 to 15 simbols"})
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

// @Summary SignUp
// @Tags auth
// @Descriotion create account
// @ID signIn
// @Accept json
// @Produce json
// @Param input body models.User true "account info"
// @Success 200 {object} models.user

func (h *Handler) SignIn(c *gin.Context) {
	var user models.User
	//parse request
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//check user in db
	user.Password = h.service.Auth.HashingPassword(user.Password)
	id, role, err := h.service.Repository.GetUser(&user)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	//create tokens
	t, rt, err := h.service.Auth.GenerateTokenPair(id, role)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"access token": t, "refresh token": rt})
}

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
		c.JSON(http.StatusUnauthorized, gin.H{"error": "not valid refresh token"})
		return
	}
	//if validate is ok -> create new tokens
	t, rt, err := h.service.Auth.GenerateTokenPair(id, role)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	//sent response
	c.JSON(http.StatusOK, gin.H{"access token": t, "refresh token": rt})

}
