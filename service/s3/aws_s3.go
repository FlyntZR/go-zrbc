package s3

import (
	"context"
	"fmt"
	"go-zrbc/config"
	"go-zrbc/view"
	"os"
	"strings"
	"time"

	awsS3 "go-zrbc/pkg/oss"
	"go-zrbc/pkg/utils"
	"go-zrbc/pkg/xlog"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

type S3Service interface {
	UploadFile(ctx context.Context, req *view.UploadFileReq) (*view.UploadFileResp, error)
	DeleteFile(ctx context.Context, bucket, fileKey string) error
	GetFile(ctx context.Context, bucket, fileKey string) (*view.GetFileResp, error)
	GetSignFile(ctx context.Context, bucket, fileKey string, expires int64) (*view.GetFileResp, error)
}

type s3Service struct {
	s3Client *s3.Client
}

func NewS3Service() S3Service {
	s3Client := awsS3.S3Client()
	srv := &s3Service{
		s3Client: s3Client,
	}

	return srv
}

func (s *s3Service) UploadFile(ctx context.Context, req *view.UploadFileReq) (*view.UploadFileResp, error) {
	defer func() {
		e := recover()
		if e != nil {
			fmt.Println("Recover from Panic Error, Failed to Upload File", e)
		}
		os.RemoveAll(req.LocalFilePath)
	}()

	// Open the file for reading
	file, err := os.Open(req.LocalFilePath)
	if err != nil {
		xlog.Errorf("error open file, err:%+v", err)
		return nil, err
	}
	defer file.Close()

	if req.FileKey == "" {
		keyDir := "dev"
		if config.Global.GinMode == "test" {
			keyDir = "test"
		} else if config.Global.GinMode == "prod" {
			keyDir = "prod"
		}
		// req.FileKey = "dev/" + util.GetFileKeyNewNew(req.LocalFilePath)
		req.FileKey = fmt.Sprintf("%s/%s", keyDir, utils.GetFileKeyNewNew(req.LocalFilePath))
	}
	if req.ContentType == "" {
		req.ContentType = "binary/octet-stream"
	}
	// 视频放yingshi-video桶，其他放yingshi-manager-image桶
	bucket := "yingshi-manager-image"
	if req.ContentType == "audio/mpeg" {
		bucket = "yingshi-video"
	}
	input := &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(req.FileKey),
		Body:        file,
		ACL:         types.ObjectCannedACLPublicRead,
		ContentType: aws.String(req.ContentType),
	}
	awsUrl, err := awsS3.PutFile(context.TODO(), s.s3Client, input)
	if err != nil {
		xlog.Errorf("error aws s3 put file, err:%+v", err)
		return nil, err
	}

	// 如果是上传的app打包文件，则返回cdn的链接
	if len(awsUrl) > 13 && (awsUrl[len(awsUrl)-13:] == ".mobileconfig" || awsUrl[len(awsUrl)-4:] == ".apk") {
		awsUrl = strings.Replace(awsUrl, "yingshi-manager-image.s3.ap-east-1.amazonaws.com", "ysimg.ejdjsn.com", -1)
	}
	return &view.UploadFileResp{
		FileUrl: awsUrl,
		FileKey: req.FileKey,
	}, nil
}

func (s *s3Service) DeleteFile(ctx context.Context, bucket, fileKey string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileKey),
	}
	return awsS3.DeleteFile(context.TODO(), s.s3Client, input)
}

func (s *s3Service) GetFile(ctx context.Context, bucket, fileKey string) (*view.GetFileResp, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileKey),
	}

	awsUrl, err := awsS3.GetFile(context.TODO(), s.s3Client, input)
	if err != nil {
		xlog.Errorf("error aws s3 get file, err:%+v", err)
		return nil, err
	}
	return &view.GetFileResp{
		FileUrl: awsUrl,
		FileKey: fileKey,
	}, nil
}

func (s *s3Service) GetSignFile(ctx context.Context, bucket, fileKey string, expires int64) (*view.GetFileResp, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(fileKey),
	}

	signClient := s3.NewPresignClient(s.s3Client)
	awsUrl, err := awsS3.GetPresignedURL(context.TODO(), signClient, input, time.Duration(expires)*time.Second)
	if err != nil {
		xlog.Errorf("error aws s3 get sign file, err:%+v", err)
		return nil, err
	}
	return &view.GetFileResp{
		FileUrl: awsUrl,
		FileKey: fileKey,
	}, nil
}
