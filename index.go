package imgShareAPI

//package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"slices"
	"strings"

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

type ImageListing struct {
	Name string `json:"name"`
	Link string `json:"link"`
}

const (
	projectName = "img-share-api-project"
	bucketName  = "img-share-api-func-bucket"
)

var validFileTypes = []string{".jpg", ".jpeg", ".jpe", ".jif", ".jfif", ".jfi", ".png",
	".gif", ".webp", ".tiff", ".tif", ".psd", ".raw", ".arw", ".cr2", ".nrw", ".k25", ".bmp", ".dib",
	".heif", ".heic", ".ind", ".indd", ".indt", ".jp2", ".j2k", ".jpf", ".jpx", ".jpm", ".mj2", ".svg",
	".svgz", ".ai", ".eps", ".pdf"}

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
	c.String(http.StatusOK, "img-share-api\nFile upload POST requests should be made with \"file\" as form name")
}

func getAllImages(c *gin.Context) {
	query := &storage.Query{Prefix: user.filePath}

	var imageListings []ImageListing
	//var imageNames []string
	it := user.bucketHandle.Objects(*user.context, query)

	//skip the foldername
	_, err := it.Next()
	if err == iterator.Done {
		//c.IndentedJSON(http.StatusOK, imageNames)
		c.IndentedJSON(http.StatusOK, imageListings)
	}
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate through objects in bucket: %v", err)
			return
		}
		imageListings = append(imageListings, ImageListing{attrs.Name, "https://storage.googleapis.com/" + bucketName + "/" + attrs.Name})
		// imageNames = append(imageNames, attrs.Name)
		// imageNames = append(imageNames, "https://storage.googleapis.com/"+bucketName+"/"+attrs.Name)
	}

	//c.IndentedJSON(http.StatusOK, imageNames)
	c.IndentedJSON(http.StatusOK, imageListings)
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
			return
		}

		imageNames = append(imageNames, attrs.Name)
		imageNames = append(imageNames, attrs.MediaLink)
	}

	c.IndentedJSON(http.StatusOK, imageNames)
}

func uploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	extensionIndex := strings.LastIndex(file.Filename, ".")
	if len(file.Filename) == 0 || extensionIndex >= len(file.Filename) || extensionIndex == -1 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Incorrect filetype or file name. Please use an image file of approved type",
			"part2": "length check",
		})
		return
	}
	extension := file.Filename[strings.LastIndex(file.Filename, "."):len(file.Filename)]
	if !slices.Contains(validFileTypes, extension) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":     "Incorrect filetype or file name. Please use an image file of approved type",
			"part2":     "filename check",
			"extension": extension,
		})
		return
	}

	blob, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	//ADD TIMEOUT FEATURE HERE
	//error checking + messages for each of these
	//assign ID + other appropriate metadata
	wc := *user.userClient.Bucket(bucketName).Object(user.filePath + file.Filename).NewWriter(*user.context)
	io.Copy(&wc, blob)
	wc.Close()
	c.JSON(http.StatusOK, gin.H{
		//"message": "upload successful",
		"message": file.Filename + " uploaded successfully",
	})
}

func deleteImage(c *gin.Context) {
	c.String(http.StatusOK, "delete image")
}
