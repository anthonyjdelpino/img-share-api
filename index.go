package imgShareAPI

import (
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/gin-gonic/gin"
)

func init() {
	route := gin.Default()
	route.GET("/", imgShareAPIFunc)
	functions.HTTP("imgShareAPIFunc", route.Handler().ServeHTTP)
}

func imgShareAPIFunc(c *gin.Context) {
	c.String(http.StatusOK, "img-share-api")
}
