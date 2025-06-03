package main

import (
	"context"
	awsS3 "go-zrbc/pkg/oss"
	"os"
	"testing"

	"go-zrbc/pkg/utils"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

func TestS3_upload(t *testing.T) {
	// Create an S3 service client
	client := awsS3.S3Client()

	// Open the file for reading
	file, err := os.Open("./FB3Y23869XV5.jpg")
	if err != nil {
		t.Fatal("failed to open file, err: ", err)
	}
	defer file.Close()

	input := &s3.PutObjectInput{
		Bucket: aws.String("yingshi-manager-image"),
		Key:    aws.String("dev/FB3Y23869XV6.jpg"),
		Body:   file,
		ACL:    types.ObjectCannedACLPublicRead,
		//ContentType: aws.String("audio/mpeg"),
		ContentType: aws.String("image/jpeg"),
	}

	output, err := awsS3.PutFile(context.TODO(), client, input)
	if err != nil {
		t.Fatal("error put file, err:", err)
	}

	t.Log("File uploaded successfully, output:", output)
}

func TestS3_GetNewKeyOld(t *testing.T) {
	filekey := utils.GetFileKeyNewNew("./FB3Y23869XV5.jpeGG")
	t.Log("File uploaded successfully, output:", filekey)
}

func TestS3_GetNewKeyNew(t *testing.T) {
	content, err := os.ReadFile("/Users/Admin1/Downloads/FB3Y23869XV5.jpg")
	if err != nil {
		t.Fatal("error read file, err:", err)
	}

	filekey := utils.GetFileKeyNew(content)
	t.Log("File uploaded successfully, output:", filekey)
}

func TestS3_Mobileconfig(t *testing.T) {
	awsURL := "leconfig"

	if len(awsURL) > 13 && awsURL[len(awsURL)-13:] == ".mobileconfig" {
		t.Log("true")
	} else {
		t.Log("false")
	}
}
