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
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

type StorageUser struct {
	userClient   *storage.Client
	context      *context.Context
	bucketHandle *storage.BucketHandle
	filePath     string
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
	route.DELETE("/images/:id", deleteImage)

	functions.HTTP("imgShareAPIFunc", route.Handler().ServeHTTP)

	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithoutAuthentication())
	if err != nil {
		log.Fatalf("Failed to create a Google Cloud Storage Client: %v", err)
	}

	bucket := client.Bucket(bucketName)
	user = &StorageUser{
		userClient:   client,
		context:      &ctx,
		bucketHandle: bucket,
		filePath:     "images/",
	}
}

func imgShareAPIFunc(c *gin.Context) {
	// provide an api guide here?
	// 	could use an html thing
	//	https://gin-gonic.com/docs/examples/html-rendering/
	c.String(http.StatusOK, "img-share-api")
}

func getAllImages(c *gin.Context) {
	query := &storage.Query{Prefix: user.filePath}

	var imageNames []string
	it := user.bucketHandle.Objects(*user.context, query)
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate through objects in bucket: %v", err)
		}

		imageNames = append(imageNames, attrs.Name)
	}

	c.IndentedJSON(http.StatusOK, imageNames)
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
