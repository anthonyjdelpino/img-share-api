package imgShareAPI

import (
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/gin-gonic/gin"
)

func init() {
	route := gin.Default()
	route.GET("/", imgShareAPIFunc)
	route.GET("/images", getAllImages)
	route.GET("/images/:id", getSpecificImage)
	route.POST("/images", uploadImage)
	route.DELETE("images/:id", deleteImage)

	functions.HTTP("imgShareAPIFunc", route.Handler().ServeHTTP)
}

func imgShareAPIFunc(c *gin.Context) {
	// provide an api guide here?
	c.String(http.StatusOK, "img-share-api")
}

func getAllImages(c *gin.Context) {
	c.String(http.StatusOK, "list all")
}

func getSpecificImage(c *gin.Context) {
	c.String(http.StatusOK, "list specific")
}

func uploadImage(c *gin.Context) {
	c.String(http.StatusOK, "upload image")
}

func deleteImage(c *gin.Context) {
	c.String(http.StatusOK, "delete image")
}
