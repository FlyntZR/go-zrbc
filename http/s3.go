package http

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"go-zrbc/pkg/http/response"
	"go-zrbc/pkg/xlog"
	"io/ioutil"
	"path"
	"strconv"

	"go-zrbc/view"

	s3 "go-zrbc/service/s3"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type S3Handler struct {
	srv s3.S3Service
}

func NewOssHandler(srv s3.S3Service) *S3Handler {
	return &S3Handler{
		srv: srv,
	}
}

func (s3 *S3Handler) SetRouter(r *gin.Engine) {
	// s3上传图片或视频到aws
	r.POST("/v1/upload_file", s3.UploadFile)
	// s3上传图片或视频base64到aws
	r.POST("/v1/upload_file_content", s3.UploadFileContent)
	// s3删除文件
	r.POST("/v1/delete_file", s3.DeleteFile)
	// s3下载文件
	r.POST("/v1/get_file", s3.GetFile)
	// s3获取sign url
	r.POST("/v1/get_sign_file", s3.GetSignFile)
}

// swagger:route POST /v1/upload_file S3接口 UploadFile
// s3上传图片或视频到aws
// consumes:
//   - multipart/form-data
//
// responses:
//
//	200: UploadFileResp
//	500: CommonError
func (s3 *S3Handler) UploadFile(c *gin.Context) {
	var req view.UploadFileReq
	req.FileKey, _ = c.GetPostForm("file_key")
	req.ContentType, _ = c.GetPostForm("content_type")
	f, err := c.FormFile("file")
	if err != nil {
		response.BadRequestResp(c, errors.New("no file is received"))
		return
	}
	if req.ContentType != "" && req.ContentType != "image/jpeg" && req.ContentType != "image/png" && req.ContentType != "audio/mpeg" {
		xlog.Errorf("content-type err, req:%+v", req)
		response.ErrResp(c, errors.New("content-type must be image/jpeg or image/png or audio/mpeg or empmty"))
		return
	}
	localFilePath := path.Join("./", f.Filename)
	c.SaveUploadedFile(f, localFilePath)
	req.LocalFilePath = localFilePath

	resp, err := s3.srv.UploadFile(context.TODO(), &req)
	if err != nil {
		fmt.Println("err2:", err)
		response.ErrResp(c, err)
		return
	}
	response.JsonResp(c, resp)
}

// swagger:route POST /v1/upload_file_content S3接口 UploadFileContent
// s3上传图片或视频base64到aws
// consumes:
//   - multipart/form-data
//
// responses:
//
//	200: UploadFileResp
//	500: CommonError
func (s3 *S3Handler) UploadFileContent(c *gin.Context) {
	var req view.UploadFileReq
	req.FileKey, _ = c.GetPostForm("file_key")
	req.ContentType, _ = c.GetPostForm("content_type")
	base64Data, _ := c.GetPostForm("f_base64_data")
	if req.ContentType != "" && req.ContentType != "image/jpeg" && req.ContentType != "image/png" && req.ContentType != "audio/mpeg" {
		xlog.Errorf("content-type err, req:%+v", req)
		response.ErrResp(c, errors.New("content-type must be image/jpeg or image/png or audio/mpeg or empmty"))
		return
	}
	localFilePath := path.Join("./", uuid.New().String())
	data, err := base64.StdEncoding.DecodeString(base64Data)
	if err != nil {
		xlog.Errorf("err:%v", err)
		response.ErrResp(c, err)
		return
	}
	err = ioutil.WriteFile(localFilePath, data, 0644)
	if err != nil {
		xlog.Errorf("err:%v", err)
		response.ErrResp(c, err)
		return
	}
	req.LocalFilePath = localFilePath

	resp, err := s3.srv.UploadFile(context.TODO(), &req)
	if err != nil {
		fmt.Println("err2:", err)
		response.ErrResp(c, err)
		return
	}
	response.JsonResp(c, resp)
}

// swagger:route POST /v1/delete_file S3接口 DeleteFile
// s3删除文件
// consumes:
//   - multipart/form-data
//
// responses:
//
//	200: ok
//	500: CommonError
func (s3 *S3Handler) DeleteFile(c *gin.Context) {
	bucket, _ := c.GetPostForm("bucket")
	fileKey, _ := c.GetPostForm("file_key")
	if bucket == "" || fileKey == "" {
		xlog.Errorf("error delete file, bucket:%s, file key:%s", bucket, fileKey)
		response.ErrResp(c, errors.New("bucket or file key is empty"))
		return
	}

	err := s3.srv.DeleteFile(context.TODO(), bucket, fileKey)
	if err != nil {
		fmt.Println("err2:", err)
		response.ErrResp(c, err)
		return
	}
	response.JsonResp(c, "ok")
}

// swagger:route POST /v1/get_file S3接口 GetFile
// s3获取url
// consumes:
//   - multipart/form-data
//
// responses:
//
//	200: GetFileResp
//	500: CommonError
func (s3 *S3Handler) GetFile(c *gin.Context) {
	bucket, _ := c.GetPostForm("bucket")
	fileKey, _ := c.GetPostForm("file_key")
	if bucket == "" || fileKey == "" {
		xlog.Errorf("error get file, bucket:%s, file key:%s", bucket, fileKey)
		response.ErrResp(c, errors.New("bucket or file key is empty"))
		return
	}

	resp, err := s3.srv.GetFile(context.TODO(), bucket, fileKey)
	if err != nil {
		fmt.Println("err2:", err)
		response.ErrResp(c, err)
		return
	}
	response.JsonResp(c, resp)
}

// swagger:route POST /v1/get_sign_file S3接口 GetSignFile
// s3获取sign url
// consumes:
//   - multipart/form-data
//
// responses:
//
//	200: GetFileResp
//	500: CommonError
func (s3 *S3Handler) GetSignFile(c *gin.Context) {
	bucket, _ := c.GetPostForm("bucket")
	fileKey, _ := c.GetPostForm("file_key")
	expiresStr, _ := c.GetPostForm("expires")
	expires, _ := strconv.ParseInt(expiresStr, 10, 64)
	if bucket == "" || fileKey == "" {
		xlog.Errorf("error get file, bucket:%s, file key:%s", bucket, fileKey)
		response.ErrResp(c, errors.New("bucket or file key is empty"))
		return
	}

	resp, err := s3.srv.GetSignFile(context.TODO(), bucket, fileKey, expires)
	if err != nil {
		fmt.Println("err2:", err)
		response.ErrResp(c, err)
		return
	}
	response.JsonResp(c, resp)
}
