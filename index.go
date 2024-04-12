package imgShareAPI

//package main

import (
	"context"
	"encoding/hex"
	"io"
	"log"
	"net/http"
	"slices"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/iterator"
)

type StorageUser struct {
	client       *storage.Client
	context      *context.Context
	bucketHandle *storage.BucketHandle
	filePath     string
}

type ImageListing struct {
	Name string `json:"name"`
	Link string `json:"link"`
	ID   string `json:"id"`
}

const (
	projectName = "img-share-api-project"
	bucketName  = "img-share-api-func-bucket"
)

// All acceptable filetypes for upload
var validFileTypes = []string{".jpg", ".jpeg", ".jpe", ".jif", ".jfif", ".jfi", ".png",
	".gif", ".webp", ".tiff", ".tif", ".psd", ".raw", ".arw", ".cr2", ".nrw", ".k25", ".bmp", ".dib",
	".heif", ".heic", ".ind", ".indd", ".indt", ".jp2", ".j2k", ".jpf", ".jpx", ".jpm", ".mj2", ".svg",
	".svgz", ".ai", ".eps", ".pdf"}

var user *StorageUser

func init() {
	//func main() {

	//start Gin engine and set routes from API points to handler functions
	route := gin.Default()
	route.GET("/", imgShareAPIFunc)
	route.GET("/images", getAllImages)
	route.GET("/images/:id", getSpecificImage)
	route.POST("/images", uploadImage)
	route.DELETE("/images/:id", deleteImage)

	functions.HTTP("imgShareAPIFunc", route.Handler().ServeHTTP)

	//Create Google Cloud Platform API client
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create a Google Cloud Storage Client: %v", err)
	}
	defer client.Close()

	//Populate user struct with information for later use and set bucket
	bucket := client.Bucket(bucketName)
	user = &StorageUser{
		client:       client,
		context:      &ctx,
		bucketHandle: bucket,
		filePath:     "images/",
	}

	//route.Run("localhost:8080")
}

// Here for the sake of Google Cloud Functions endpoint
func imgShareAPIFunc(c *gin.Context) {
	c.String(http.StatusOK, "img-share-api\nFile upload POST requests should be made with \"file\" as form name")
}

// Return a list of all images in the folder with a media link
func getAllImages(c *gin.Context) {
	query := &storage.Query{Prefix: user.filePath}

	var imageListings []ImageListing
	it := user.bucketHandle.Objects(*user.context, query)

	//skip the first item as it will just be the folder name
	_, err := it.Next()
	if err == iterator.Done {
		c.IndentedJSON(http.StatusOK, gin.H{
			user.filePath: imageListings,
		})
	}

	//Iterate through each item in the name path provided
	//	-note that Cloud Storage uses a flat namespace: https://cloud.google.com/storage/docs/objects
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("Failed to iterate through objects in bucket: %v", err)
			return
		}
		imageListings = append(imageListings, ImageListing{attrs.Name, "https://storage.googleapis.com/" + bucketName + "/" + attrs.Name, hex.EncodeToString(attrs.MD5)})
	}

	c.IndentedJSON(http.StatusOK, imageListings)
}

// Return an image's information as identified by it's ID field
//
//	-ID is derived by hex-encoding the image's Cloud Storage object's MD5 hash
func getSpecificImage(c *gin.Context) {
	query := &storage.Query{Prefix: user.filePath}

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
		if hex.EncodeToString(attrs.MD5) == c.Param("id") {
			c.JSON(http.StatusOK, ImageListing{attrs.Name, "https://storage.googleapis.com/" + bucketName + "/" + attrs.Name, hex.EncodeToString(attrs.MD5)})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "No image found with ID: " + c.Param("id"),
	})
}

// Upload images to the bucket in the folder specified in user.filePath
func uploadImage(c *gin.Context) {
	//Get multipart file from the POST request
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	//Make sure the file name uses an acceptable filetype extension
	extensionIndex := strings.LastIndex(file.Filename, ".")
	if len(file.Filename) == 0 || extensionIndex >= len(file.Filename) || extensionIndex == -1 {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Incorrect filetype or file name. Please use an image file of approved type",
		})
		return
	}
	extension := file.Filename[strings.LastIndex(file.Filename, "."):len(file.Filename)]
	if !slices.Contains(validFileTypes, extension) {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Incorrect filetype or file name. Please use an image file of approved type",
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

	_, cancel := context.WithTimeout(c, time.Second*10)
	defer cancel()

	//Create storage writer to write new object to bucket
	wc := *user.client.Bucket(bucketName).Object(user.filePath + file.Filename).NewWriter(*user.context)
	io.Copy(&wc, blob)
	wc.Close()

	c.JSON(http.StatusOK, gin.H{
		"message": file.Filename + " uploaded successfully",
	})
}

// Delete image, identified in request by ID field
func deleteImage(c *gin.Context) {
	query := &storage.Query{Prefix: user.filePath}

	var object *storage.ObjectHandle

	it := user.bucketHandle.Objects(*user.context, query)
	_, err := it.Next()
	if err == iterator.Done {
		c.JSON(http.StatusOK, gin.H{
			"error": "Image not found. Folder is empty",
		})
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
		if hex.EncodeToString(attrs.MD5) == c.Param("id") {
			object = user.client.Bucket(bucketName).Object(attrs.Name)
		}
		if object == nil {
			//error here and return

			log.Fatalf("Object nil: %v", attrs.Name)
			return
		}
		attrs, err = object.Attrs(c)
		if err != nil {
			c.String(http.StatusInternalServerError, "Object does not exist")
			return
		}

		object = object.If(storage.Conditions{GenerationMatch: attrs.Generation})
		name := object.ObjectName()

		_, cancel := context.WithTimeout(c, time.Second*10)
		defer cancel()
		if err := object.Delete(c); err != nil {
			log.Fatalf("Failed to delete object: %v", err)
			return
		} else {
			c.JSON(http.StatusOK, gin.H{
				"message": name + "was deleted",
			})
			return
		}
	}
}
