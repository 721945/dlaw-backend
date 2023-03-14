package google_storage

import (
	"cloud.google.com/go/storage"
	"context"
	"github.com/721945/dlaw-backend/infrastructure/google_vision"
	"github.com/721945/dlaw-backend/libs"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

type GoogleStorage struct {
	bucket string
	logger *libs.Logger
	vision google_vision.GoogleVision
}

func NewGoogleStorage(env libs.Env, logger *libs.Logger, vision google_vision.GoogleVision) GoogleStorage {
	return GoogleStorage{bucket: env.Bucket, logger: logger, vision: vision}
}

func (g GoogleStorage) UploadFile(file multipart.File, fileName string) (string, error) {
	/// Now we are doing for put all file in the bucket and use database to collect all information
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", err
	}
	defer func(client *storage.Client) {
		err := client.Close()
		if err != nil {
			panic(err)
		}
	}(client)

	bucket := client.Bucket(g.bucket)

	obj := bucket.Object(fileName)

	wc := obj.NewWriter(ctx)

	if _, err = io.Copy(wc, file); err != nil {
		return "", err
	}

	if err := wc.Close(); err != nil {
		return "", err
	}

	url := "https://storage.cloud.google.com/" + g.bucket + "/" + fileName

	//sourceUrl := "gs://" + g.bucket + "/" + fileName
	//targetUrl := "gs://" + g.bucket + "/" + fileName + ".txt"

	visionText, err := g.vision.GetTextFromPdfUrl(*obj)
	//visionText, err := g.vision.GetTextFromPdfUrl(sourceUrl, targetUrl)

	g.logger.Info(visionText)

	return url, nil
}

func (g GoogleStorage) GetSignedUrl(amount int) ([]string, error) {
	var urls []string
	for i := 0; i < amount; i++ {
		url, err := g.getSignedUrl("abc.pdf")
		if err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}
	return urls, nil
}

func (g GoogleStorage) getSignedUrl(objectName string) (string, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		return "", err
	}
	defer func(client *storage.Client) {
		err := client.Close()
		if err != nil {
			panic(err)
		}
	}(client)

	bucket := client.Bucket(g.bucket)

	opts := &storage.SignedURLOptions{
		Method:      http.MethodPut,
		ContentType: "application/octet-stream",
		Expires:     time.Now().Add(10 * time.Minute),
	}

	url, err := bucket.SignedURL(objectName, opts)

	if err != nil {
		return "", err
	}
	return url, nil

}
