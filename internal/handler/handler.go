package handler

import (
	_ "JWT_auth/docs"
	"JWT_auth/internal/models"
	"JWT_auth/internal/service"
	"math"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// @title Internet-shop API
// @version 1.0
// @description API Server for catalog of internet-shop

// @host localhost:8000

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

const userRole = "user"

type Handler struct {
	service *service.Service
	logger  *logrus.Logger
}

func NewHandler(service *service.Service, logger *logrus.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

//init routers =====================================================================
func (h *Handler) Init() *gin.Engine {

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Logger())

	//authorization routing
	auth := router.Group("/auth")
	{
		auth.POST("/signUp", h.AuthMiddleware, h.IsAdminMiddleware, h.SignUp)
		auth.POST("/signIn", h.SignIn)
		auth.POST("/update", h.TokenRefreshing)
		auth.POST("/admin", h.AddAddmin)
	}

	catalog := router.Group("/catalog").Use(h.AuthMiddleware)
	{
		//add category
		catalog.POST("/category", h.IsAdminMiddleware, h.AddCategory)
		//add product
		catalog.POST("/product", h.IsAdminMiddleware, h.AddProduct)
		//change products visible in catalog
		catalog.PUT("/product/change", h.IsAdminMiddleware, h.ChangeVisible)
		//get all catalog
		catalog.GET("/", h.GetCatalog)
		//search
		catalog.GET("/search", h.Search)
	}

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Not allowed request"})
	})

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}

//adding new category ============================================================
func (h *Handler) AddCategory(c *gin.Context) {
	//bindig request
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		h.logger.Error(err)
		c.Status(http.StatusBadRequest)
		return
	}
	if err := h.service.Repository.AddCategory(&category); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

}

// adding new product ==============================================================
func (h *Handler) AddProduct(c *gin.Context) {
	//bindig request
	var product models.ProductDTO
	if err := c.ShouldBindJSON(&product); err != nil {
		h.logger.Error(err)
		c.Status(http.StatusBadRequest)
		return
	}
	//Round float to 2 decimal places
	product.Price = math.Round(product.Price*100) / 100
	product.Weight = math.Round(product.Weight*100) / 100
	product.Valume = math.Round(product.Valume*100) / 100
	if err := h.service.Repository.AddProduct(&product); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
}

//changing product visible in catalog ================================================
func (h *Handler) ChangeVisible(c *gin.Context) {
	//bindig request
	var visible models.Visible
	if err := c.ShouldBindJSON(&visible); err != nil {
		h.logger.Error(err)
		c.Status(http.StatusBadRequest)
		return
	}
	if err := h.service.Repository.ChangeVisible(&visible); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

}

//show all catalog ======================================================================
func (h *Handler) GetCatalog(c *gin.Context) {
	catalog, err := h.service.Repository.GetCatalog()
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, catalog)
}

//shearch product in catalog =============================================================
func (h *Handler) Search(c *gin.Context) {
	category := c.Query("category")
	productName := c.Query("product")

	if category == "" && productName == "" {
		c.Status(http.StatusBadRequest)
		return
	}
	if category == "" {
		result, err := h.service.Repository.GetByAllCategories(productName)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, result)
		return
	} else {
		result, err := h.service.Repository.GetByCategory(productName, category)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.JSON(http.StatusOK, result)
	}

}
