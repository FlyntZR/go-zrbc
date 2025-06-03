package awsS3

import (
	"context"
	"fmt"
	"go-zrbc/pkg/xlog"
	"time"

	"net/http"
	"net/url"

	gconfig "go-zrbc/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func S3ClientOld() *s3.Client {
	// 代理服务器的URL
	proxyURL, _ := url.Parse("http://18.167.194.10:443")
	// 或者你也可以直接设置Transport
	transport := &http.Transport{
		Proxy: http.ProxyURL(proxyURL),
	}
	// 创建使用自定义Transport的Client
	client := &http.Client{
		Transport: transport,
	}

	// 创建配置选项
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion("ap-east-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(gconfig.Global.AwsKey, gconfig.Global.AwsSecret, "")),
		config.WithHTTPClient(client),
	)
	if err != nil {
		xlog.Errorf("S3 config err:%+v", err)
	}
	return s3.NewFromConfig(cfg)
}

func S3Client() *s3.Client {
	// 创建配置选项
	var cfg aws.Config
	var err error
	// 如果配置了代理，就试用代理
	if gconfig.Global.Agent != "" {
		// 代理服务器的URL
		proxyURL, _ := url.Parse(gconfig.Global.Agent)
		// 或者你也可以直接设置Transport
		transport := &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		}

		// 创建使用自定义Transport的Client
		client := &http.Client{
			Transport: transport,
		}
		cfg, err = config.LoadDefaultConfig(
			context.TODO(),
			config.WithRegion("ap-east-1"),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(gconfig.Global.AwsKey, gconfig.Global.AwsSecret, "")),
			config.WithHTTPClient(client),
		)
	} else {
		cfg, err = config.LoadDefaultConfig(
			context.TODO(),
			config.WithRegion("ap-east-1"),
			config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(gconfig.Global.AwsKey, gconfig.Global.AwsSecret, "")),
		)
	}
	if err != nil {
		xlog.Errorf("S3 config err:%+v", err)
	}
	return s3.NewFromConfig(cfg)
}

type S3API interface {
	PutObject(ctx context.Context, params *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	DeleteObject(ctx context.Context, params *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error)
	GetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error)
}

func PutFile(c context.Context, api S3API, input *s3.PutObjectInput) (string, error) {
	if _, err := api.PutObject(c, input); err != nil {
		return "", err
	}
	awsURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", *input.Bucket, "ap-east-1", *input.Key)
	return awsURL, nil
}

func DeleteFile(c context.Context, api S3API, input *s3.DeleteObjectInput) error {
	if _, err := api.DeleteObject(c, input); err != nil {
		return err
	}
	return nil
}

func GetFile(c context.Context, api S3API, input *s3.GetObjectInput) (string, error) {
	if _, err := api.GetObject(c, input); err != nil {
		return "", err
	}
	awsURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", *input.Bucket, "ap-east-1", *input.Key)
	return awsURL, nil
}

type S3PresignAPI interface {
	PresignGetObject(ctx context.Context, params *s3.GetObjectInput, optFns ...func(*s3.PresignOptions)) (*v4.PresignedHTTPRequest, error)
}

func GetPresignedURL(c context.Context, api S3PresignAPI, input *s3.GetObjectInput, dur time.Duration) (string, error) {
	resp, err := api.PresignGetObject(c, input, s3.WithPresignExpires(dur))
	if err != nil {
		return "", err
	}
	return resp.URL, nil
}
