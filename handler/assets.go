package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/globalsign/mgo/bson"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"github.com/sonereker/api-manage-mor-app/model"
	"log"
	"mime/multipart"
	"net/http"
	"path/filepath"
)

func GetAllAssets(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	var assets []model.Asset
	db.Find(&assets)
	respondJSON(w, http.StatusOK, assets)
}

func CreateAsset(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	asset := model.Asset{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&asset); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(&asset).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondJSON(w, http.StatusCreated, asset)
}

func GetAsset(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]
	asset := getAssetOr404(db, id, w, r)
	if asset == nil {
		return
	}
	respondJSON(w, http.StatusOK, asset)
}

func UpdateAsset(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	uuid := vars["uuid"]
	asset := getAssetOr404(db, uuid, w, r)
	if asset == nil {
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&asset); err != nil {
		respondError(w, http.StatusBadRequest, err.Error())
		return
	}
	defer r.Body.Close()

	if err := db.Save(&asset).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusOK, asset)
}

func DeleteAsset(db *gorm.DB, w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	uuid := vars["uuid"]
	asset := getAssetOr404(db, uuid, w, r)
	if asset == nil {
		return
	}
	if err := db.Delete(&asset).Error; err != nil {
		respondError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondJSON(w, http.StatusNoContent, nil)
}

// https://medium.com/spankie/upload-images-to-aws-s3-bucket-in-a-golang-web-application-2612bea70dd8
func UploadFile(w http.ResponseWriter, r *http.Request) {
	maxSize := int64(1024000) // allow only 1MB of file size

	err := r.ParseMultipartForm(maxSize)
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "Image too large. Max Size: %v", maxSize)
		return
	}
	file, fileHeader, err := r.FormFile("profile_picture")
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, "Could not get uploaded file")
		return
	}
	defer file.Close()
	// create an AWS session which can be
	// reused if we're uploading many files
	s, err := session.NewSession(&aws.Config{
		Region: aws.String("us-east-2"),
		Credentials: credentials.NewStaticCredentials(
			"secret-id",  // id
			"secret-key", // secret
			""),          // token can be left blank for now
	})
	if err != nil {
		fmt.Fprintf(w, "Could not upload file")
	}
	fileName, err := UploadFileToS3(s, file, fileHeader)
	if err != nil {
		fmt.Fprintf(w, "Could not upload file")
	}
	fmt.Fprintf(w, "Image uploaded successfully: %v", fileName)
}

func UploadFileToS3(s *session.Session, file multipart.File, fileHeader *multipart.FileHeader) (string, error) {
	// get the file size and read
	// the file content into a buffer
	size := fileHeader.Size
	buffer := make([]byte, size)
	file.Read(buffer)

	// create a unique file name for the file
	tempFileName := "assets/" + bson.NewObjectId().Hex() + filepath.Ext(fileHeader.Filename)

	// config settings: this is where you choose the bucket,
	// filename, content-type and storage class of the file
	// you're uploading
	_, err := s3.New(s).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String("test-bucket"),
		Key:                  aws.String(tempFileName),
		ACL:                  aws.String("public-read"), // could be private if you want it to be access by only authorized users
		Body:                 bytes.NewReader(buffer),
		ContentLength:        aws.Int64(int64(size)),
		ContentType:          aws.String(http.DetectContentType(buffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
		StorageClass:         aws.String("INTELLIGENT_TIERING"),
	})
	if err != nil {
		return "", err
	}

	return tempFileName, err
}

func getAssetOr404(db *gorm.DB, uuid string, w http.ResponseWriter, r *http.Request) *model.Asset {
	asset := model.Asset{}
	if err := db.First(&asset, uuid).Error; err != nil {
		respondError(w, http.StatusNotFound, err.Error())
		return nil
	}
	return &asset
}
