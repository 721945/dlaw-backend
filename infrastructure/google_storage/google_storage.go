package google_storage

import (
	"cloud.google.com/go/storage"
	"context"
	"errors"
	"fmt"
	"github.com/721945/dlaw-backend/infrastructure/google_vision"
	"github.com/721945/dlaw-backend/libs"
	"github.com/721945/dlaw-backend/models"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
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

func (g GoogleStorage) GiveAccessPublic(name, fileName string) (string, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)

	if err != nil {
		return "", err
	}

	defer func(client *storage.Client) {
		err := client.Close()
		if err != nil {
			g.logger.Error("Error closing Google Cloud Storage client: %v", err)
			panic(err)
		}
	}(client)

	acl := client.Bucket(g.bucket).Object(name).ACL()

	// Update the object ACL to make it publicly accessible.
	if err := acl.Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		log.Fatal(err)
	}

	return g.getPublicFile(name, fileName), nil
}

func (g GoogleStorage) GiveAccessPrivate(name string) error {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)

	if err != nil {
		return err
	}

	defer func(client *storage.Client) {
		err := client.Close()
		if err != nil {
			g.logger.Error("Error closing Google Cloud Storage client: %v", err)
			panic(err)
		}
	}(client)

	acl := client.Bucket(g.bucket).Object(name).ACL()

	rules, err := acl.List(ctx)
	if err != nil {
		return err
	}

	// Loop through the existing ACL entries and remove any that grant public access.
	for _, rule := range rules {
		if rule.Entity == storage.AllUsers && rule.Role == storage.RoleReader {
			if err := acl.Delete(ctx, rule.Entity); err != nil {
				return err
			}
		}
	}

	return nil
}

func (g GoogleStorage) GetFileSize(name, version string) (int64, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)

	if err != nil {
		return 0, err
	}

	defer func(client *storage.Client) {
		err := client.Close()
		if err != nil {
			g.logger.Error("Error closing Google Cloud Storage client: %v", err)
			panic(err)
		}
	}(client)

	bucket := client.Bucket(g.bucket)

	objectName := name

	if version != "" {
		objectName = fmt.Sprintf("%s/%s", name, version)
	}

	obj := bucket.Object(objectName)

	attrs, err := obj.Attrs(ctx)

	if err != nil {
		return 0, err
	}

	return attrs.Size, nil
}

func (g GoogleStorage) GetSignedFileUrls(files []models.File) (fileUrls []models.FileUrl, err error) {

	cloudNames := make([]string, len(files))
	previewCloudNames := make([]string, len(files))
	downloadNames := make([]string, len(files))
	isShared := make([]bool, len(files))

	for i, file := range files {
		cloudNames[i] = file.CloudName
		downloadNames[i] = file.Name
		previewCloudNames[i] = file.PreviewCloudName
		isShared[i] = file.IsShared
	}

	urlsCh := make(chan []string)
	previewUrlsCh := make(chan []string)

	// Run the two calls to GetSignedUrls in parallel using goroutines
	go func() {
		urls, err := g.getSignedUrls(cloudNames, []string{}, downloadNames, isShared)
		if err != nil {
			g.logger.Info(err)
			urlsCh <- nil
		} else {
			urlsCh <- urls
		}
	}()

	go func() {
		previewUrls, err := g.getSignedUrls(previewCloudNames, []string{}, downloadNames, isShared)
		if err != nil {
			g.logger.Info(err)
			previewUrlsCh <- nil
		} else {
			previewUrlsCh <- previewUrls
		}
	}()

	// Wait for both goroutines to complete and merge the results
	urls := <-urlsCh
	previewUrls := <-previewUrlsCh

	if urls == nil || previewUrls == nil {
		return nil, errors.New("error getting signed urls")
	}

	fileUrls = make([]models.FileUrl, len(files))

	for i, _ := range files {
		fileUrls[i] = models.FileUrl{
			Url:        urls[i],
			PreviewUrl: previewUrls[i],
		}
	}

	return fileUrls, nil
}

func (g GoogleStorage) getSignedUrls(names, versions, fileNames []string, isShared []bool) ([]string, error) {
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

	expirationTime := time.Now().Add(time.Hour * 2)

	bucket := client.Bucket(g.bucket)

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
		// If share file, so we just only need
		if isShared[i] {
			urls[i] = g.getPublicFile(name, fileNames[i])
			continue
		}

		wg.Add(1)
		go func(i int, name string, version string, downloadName string) {
			defer wg.Done()

			objectName := name

			if version != "" {
				objectName = fmt.Sprintf("%s/%s", name, version)
			}

			signed, err := bucket.SignedURL(objectName, opts)

			if err != nil {
				log.Printf("Error signing URL for object %s: %v", name, err)
			} else {
				urls[i] = signed
				urls[i] += "&response-content-disposition=attachment;filename=" + downloadName
			}
		}(i, name, versions[i], fileNames[i])
	}

	wg.Wait()

	return urls, nil
}

func (g GoogleStorage) getPublicFile(name, downloadName string) string {
	//encodedBucketName := base64.URLEncoding.EncodeToString([]byte(g.bucket))
	//encodedObjectName := base64.URLEncoding.EncodeToString([]byte(name))
	urlStr := fmt.Sprintf("https://storage.googleapis.com/%s/%s?response-content-disposition=attachment;filename=%s", g.bucket, name, url.PathEscape(downloadName))
	//urlStr := fmt.Sprintf("https://storage.googleapis.com/%s/%s?response-content-disposition=attachment;filename=%s", encodedBucketName, encodedObjectName, url.PathEscape(downloadName))
	downloadURL, _ := url.Parse(urlStr)
	return downloadURL.String()
}
