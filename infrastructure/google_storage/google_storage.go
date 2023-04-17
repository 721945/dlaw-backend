package google_storage

import (
	"cloud.google.com/go/storage"
	"context"
	"fmt"
	"github.com/721945/dlaw-backend/infrastructure/google_vision"
	"github.com/721945/dlaw-backend/libs"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
	"sync"
	"time"
)

type GoogleStorage struct {
	bucket      string
	logger      *libs.Logger
	vision      google_vision.GoogleVision
	clientEmail string
	privateKey  string
}

func NewGoogleStorage(
	env libs.Env,
	logger *libs.Logger,
	vision google_vision.GoogleVision,
) GoogleStorage {
	return GoogleStorage{
		bucket:      env.Bucket,
		logger:      logger,
		vision:      vision,
		clientEmail: env.GoogleCloudStorageClientEmail,
		privateKey:  env.GoogleCloudStoragePrivateKey,
	}
}

func (g GoogleStorage) UploadFile(file multipart.File, fileName string) (string, error) {
	/// Now we are doing for put all file in the bucket and use database to collect all information
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	g.logger.Info("Start of upload file")
	if err != nil {
		return "", err
	}
	defer func(client *storage.Client) {
		g.logger.Info("End of upload file")
		err := client.Close()
		if err != nil {
			panic(err)
		}
	}(client)

	bucket := client.Bucket(g.bucket)

	obj := bucket.Object(fileName)

	wc := obj.NewWriter(ctx)

	g.logger.Info("Start of copy file")
	if _, err = io.Copy(wc, file); err != nil {
		return "", err
	}

	g.logger.Info("End of copy file")
	if err := wc.Close(); err != nil {
		return "", err
	}

	attrs, err := obj.Attrs(ctx)

	if err != nil {
		return "", err
	}

	version := strconv.FormatInt(attrs.Generation, 10)

	return version, nil
}

func (g GoogleStorage) GetSignedUrls(names, versions, fileNames []string) ([]string, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)

	if err != nil {
		return []string{}, err
	}

	defer func(client *storage.Client) {
		err := client.Close()
		if err != nil {
			g.logger.Error("Error closing Google Cloud Storage client: %v", err)
			panic(err)
		}
	}(client)

	// TODO : Change this to read from env
	expirationTime := time.Now().Add(time.Hour * 2)

	bucket := client.Bucket(g.bucket)

	// TODO : Make this more security
	opts := &storage.SignedURLOptions{
		GoogleAccessID: g.clientEmail,
		Method:         http.MethodGet,
		Expires:        expirationTime,
		PrivateKey:     []byte(g.privateKey),
	}

	urls := make([]string, len(names))

	if len(versions) == 0 {
		versions = make([]string, len(names))
	}
	var wg sync.WaitGroup
	for i, name := range names {
		wg.Add(1)
		go func(i int, name string, version string, downloadName string) {
			defer wg.Done()

			objectName := name

			if version != "" {
				objectName = fmt.Sprintf("%s/%s", name, version)
			}

			url, err := bucket.SignedURL(objectName, opts)

			if err != nil {
				log.Printf("Error signing URL for object %s: %v", name, err)
			} else {
				urls[i] = url
				urls[i] += "&response-content-disposition=attachment;filename=" + downloadName
			}
		}(i, name, versions[i], fileNames[i])
	}

	wg.Wait()

	return urls, nil
}

//func signed
