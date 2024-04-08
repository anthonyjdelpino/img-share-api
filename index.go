// package imgShareAPI
package main

/*TODO:
implement functions
structured error messages
authentication?
*/

import (
	"context"
	"io"
	"log"
	"net/http"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
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
	//func main() {
	route := gin.Default()
	route.GET("/", imgShareAPIFunc)
	route.GET("/images", getAllImages)
	route.GET("/images/:name", getSpecificImage)
	route.POST("/images", uploadImage)
	route.DELETE("/images/:id", deleteImage)

	functions.HTTP("imgShareAPIFunc", route.Handler().ServeHTTP)

	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create a Google Cloud Storage Client: %v", err)
	}
	defer client.Close()

	bucket := client.Bucket(bucketName)
	user = &StorageUser{
		userClient:   client,
		context:      &ctx,
		bucketHandle: bucket,
		filePath:     "images/",
	}

	//route.Run("localhost:8080")
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

func getSpecificImage(c *gin.Context) { //functions more like a search
	query := &storage.Query{Prefix: user.filePath + c.Param("name")}

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

func uploadImage(c *gin.Context) {
	file, err := c.FormFile("file_input")
	if err != nil {
		log.Fatalf("Failed to receive file from request.")
		return
	}
	blob, _ := file.Open()

	//add timeout?
	wc := *user.userClient.Bucket(bucketName).Object(user.filePath + file.Filename).NewWriter(*user.context)
	io.Copy(&wc, blob)
	wc.Close()
	c.String(http.StatusOK, "upload image")
}

func deleteImage(c *gin.Context) {
	c.String(http.StatusOK, "delete image")
}
