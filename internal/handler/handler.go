package handler

import (
	"JWT_auth/internal/model"
	"JWT_auth/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

const userRole = "user"

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
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
	//add category
	router.POST("/category", h.AuthMiddleware, h.IsAdminMiddleware, h.AddCategory)

	product := router.Group("/product").Use(h.AuthMiddleware, h.IsAdminMiddleware)
	{
		//add product
		product.POST("/", h.AddProduct)
		//change products visible in catalog
		product.PUT("/change", h.ChangeVisible)
	}

	catalog := router.Group("/catalog").Use(h.AuthMiddleware)
	{ //get all catalog
		catalog.GET("/", h.GetCatalog)
		//search
		catalog.GET("/search", h.Search)
	}

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusBadRequest, gin.H{"Error": "Not allowed request"})
	})
	return router
}

//adding new category ============================================================
func (h *Handler) AddCategory(c *gin.Context) {
	//bindig request
	var category model.Category
	if err := c.ShouldBindJSON(&category); err != nil {
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
	var product model.ProductDTO
	if err := c.ShouldBindJSON(&product); err != nil {
		c.Status(http.StatusBadRequest)
		return
	}
	if err := h.service.Repository.AddProduct(&product); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

}

//changing product visible in catalog ================================================
func (h *Handler) ChangeVisible(c *gin.Context) {
	//bindig request
	var visible model.Visible
	if err := c.ShouldBindJSON(&visible); err != nil {
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
