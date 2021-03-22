package router

import (
	"backup/handler"
	"github.com/labstack/echo/v4"
)

type API struct {
	Echo *echo.Echo
	UploadHandler handler.UploadHandler
}

func (api *API) SetupRouter() {
	upload := api.Echo.Group("/upload")
	upload.POST("/multipart/init", api.UploadHandler.InitUploadHandler)
	upload.POST("/multipart/upload_part", api.UploadHandler.UploadPartHandler)
	upload.POST("/multipart/complete", api.UploadHandler.UploadCompleteHandler)
	upload.POST("/multipart/abort", api.UploadHandler.UploadPartAbort)

	bucket := api.Echo.Group("/bucket")
	bucket.POST("/created", api.UploadHandler.CreateBucketHandler)
}

