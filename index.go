package imgShareAPI

/*TODO:
implement functions
structured error messages
authentication?
*/

import (
	"context"
	"log"
	"net/http"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/gin-gonic/gin"
)

type StorageUser struct {
	userClient *storage.Client
	filePath   string
}

const (
	projectName = "img-share-api-project"
	bucketName  = "img-share-api-func-bucket"
)

var user *StorageUser

func init() {
	route := gin.Default()
	route.GET("/", imgShareAPIFunc)
	route.GET("/images", getAllImages)
	route.GET("/images/:id", getSpecificImage)
	route.POST("/images", uploadImage)
	route.DELETE("images/:id", deleteImage)

	functions.HTTP("imgShareAPIFunc", route.Handler().ServeHTTP)

	client, err := storage.NewClient(context.Background())
	if err != nil {
		log.Fatalf("Failed to create a Google Cloud Storage Client: %v", err)
	}

	user = &StorageUser{
		userClient: client,
		filePath:   "/",
	}
}

func imgShareAPIFunc(c *gin.Context) {
	// provide an api guide here?
	// 	could use an html thing
	//	https://gin-gonic.com/docs/examples/html-rendering/
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
