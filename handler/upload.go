package handler

import (
	"backup/model"
	"backup/storage"
	"encoding/json"
	"github.com/labstack/echo/v4"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type UploadHandler struct {
	//UploadRepo repository.UploadRepo
}

func (upload *UploadHandler) CreateBucketHandler(c echo.Context) error {
	request := model.ReqUploadS3{}
	if err := c.Bind(&request); err != nil {
		log.Println(err.Error())
		return c.JSON(http.StatusBadRequest, &model.ResponseDataS3{
			StatusCode: http.StatusBadRequest,
			Message:    "syntax error",
		})
	}
	storage.CreateBucketS3(request.Bucket)
	log.Printf("successfully created bucket. bucket: %s", request.Bucket)

	return c.JSON(http.StatusCreated, &model.ResponseDataS3{
		StatusCode: http.StatusCreated,
		Message: "successfully created bucket",
		Bucket: request.Bucket,
	})
}

func (upload *UploadHandler) InitUploadHandler(c echo.Context) error {
	decoder := json.NewDecoder(c.Request().Body)
	request := model.ReqUploadS3{}
	err := decoder.Decode(&request)
	if err != nil {
		log.Println(err)
	}
	uploadID, key, bucket := storage.InitUploadS3(request.Bucket, request.Key)
	log.Printf("successfully uploaded init. upload_id: %s, key: %s, bucket: %s", uploadID, key, bucket)

	return c.JSON(http.StatusOK, &model.ResponseDataS3{
		StatusCode: http.StatusOK,
		Message: "successfully uploaded init",
		Key: key,
		UploadID: uploadID,
		Bucket: bucket,
	})
}

func (upload *UploadHandler) UploadPartAbort(c echo.Context) error {
	request := model.ReqUploadS3{}
	if err := c.Bind(&request); err != nil {
		log.Println(err.Error())
		return c.JSON(http.StatusBadRequest, &model.ResponseDataS3{
			StatusCode: http.StatusBadRequest,
			Message:    "syntax error",
		})
	}
	bucket := request.Bucket
	key := request.Key
	uploadID := request.UploadID
	storage.UploadPartAbortS3(bucket, key, uploadID)
	log.Printf("successfully uploaded abort. upload_id: %s, key: %s, bucket: %s", uploadID, key, bucket)

	return c.JSON(http.StatusOK, &model.ResponseDataS3{
		StatusCode: http.StatusOK,
		Message: "successfully uploaded abort",
		Key: key,
		Bucket: bucket,
		UploadID: uploadID,
	})
}

func (upload *UploadHandler) UploadPartHandler(c echo.Context) error {
	bucket := c.Request().URL.Query().Get("bucket")
	key := c.Request().URL.Query().Get("key")
	uploadID := c.Request().URL.Query().Get("upload_id")
	partNumber, _ := strconv.Atoi(c.Request().URL.Query().Get("part_number"))
	bytePart, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		log.Println(err)
	}
	time.Sleep(500 * time.Millisecond)
	complete, _ := storage.UploadPartS3(bucket, key, uploadID, bytePart, partNumber)
	log.Printf("successfully uploaded part number: %d, byte: %d, upload_id: %s, key: %s, bucket: %s",partNumber, len(bytePart), uploadID, key, bucket)

	return c.JSON(http.StatusOK, &model.DataUploadResp{
		Parts: complete,
	})
}

func (upload *UploadHandler) UploadCompleteHandler(c echo.Context) error {
	decoder := json.NewDecoder(c.Request().Body)
	request := model.ReqUploadS3{}
	err := decoder.Decode(&request)
	if err != nil {
		log.Println(err)
	}
	bucket := request.Bucket
	key := request.Key
	uploadID := request.UploadID
	parts := request.Parts
	etag := storage.UploadCompleteS3(bucket, key, uploadID, parts)
	log.Printf("successfully uploaded complete. upload_id: %s, key: %s, bucket: %s, etag: %s", uploadID, key, bucket, etag)

	return  c.JSON(http.StatusOK, &model.ResponseDataS3{
		StatusCode: http.StatusOK,
		Message: "successfully uploaded complete",
		Key: request.Key,
		Bucket: request.Bucket,
		UploadID: request.UploadID,
		Etag: etag,
	})
}