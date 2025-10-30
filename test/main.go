package main

import (
	"context"
	"log"

	"github.com/qiniu/go-sdk/v7/storagev2/credentials"
	"github.com/qiniu/go-sdk/v7/storagev2/http_client"
	"github.com/qiniu/go-sdk/v7/storagev2/uploader"
)

func main() {
	accessKey := "l1ck5kcmqbuyEcqb9RIIyIpZW6gYRXg-uTIUi4cB"
	secretKey := "XLm1_etEN2EyY-wt35VLnyM_utDELjVrxZN5boZb"
	//
	mac := credentials.NewCredentials(accessKey, secretKey)
	localFile := "C:\\Users\\Lenovo\\OneDrive\\图片\\博客专用\\从一个坑跳入另一个坑.png"
	bucket := "wdx89"
	key := "image/github-x.png"
	uploadManager := uploader.NewUploadManager(&uploader.UploadManagerOptions{
		Options: http_client.Options{
			Credentials: mac,
		},
	})
	err := uploadManager.UploadFile(context.Background(), localFile, &uploader.ObjectOptions{
		BucketName: bucket,
		ObjectName: &key,
		CustomVars: map[string]string{
			"name": "github logo",
		},
		FileName: key, // 设置 FileName
	}, nil)
	if err != nil {
		log.Fatal(err)
	}
}
