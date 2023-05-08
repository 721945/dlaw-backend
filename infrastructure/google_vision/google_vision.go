package google_vision

import (
	"cloud.google.com/go/storage"
	vision "cloud.google.com/go/vision/apiv1"
	"cloud.google.com/go/vision/v2/apiv1/visionpb"
	"context"
	"encoding/json"
	fmt "fmt"
	"github.com/721945/dlaw-backend/libs"
	"google.golang.org/api/iterator"
	"log"
	"strings"
)

type GoogleVision struct {
	logger *libs.Logger
	bucket string
}

type GcsResponse struct {
	Responses []struct {
		FullTextAnnotation struct {
			Text string `json:"text"`
		} `json:"fullTextAnnotation"`
		Context struct {
			Uri        string `json:"uri"`
			PageNumber int    `json:"pageNumber"`
		} `json:"context"`
	} `json:"responses"`
}

func NewGoogleVision(logger *libs.Logger, env libs.Env) GoogleVision {
	return GoogleVision{logger: logger, bucket: env.Bucket}
}

func (g *GoogleVision) GetTextFromImageName(name string) (string, error) {
	ctx := context.Background()
	client, err := vision.NewImageAnnotatorClient(ctx)

	if err != nil {
		return "", err
	}

	defer func(client *vision.ImageAnnotatorClient) {
		err := client.Close()
		if err != nil {
			panic(err)
		}
	}(client)

	url := "gs://" + g.bucket + "/" + name

	image := vision.NewImageFromURI(url)

	annotations, err := client.DetectTexts(ctx, image, nil, 10)

	if err != nil {
		return "", err
	}

	g.logger.Info(annotations)

	if len(annotations) == 0 {
		return "", nil
	}

	return annotations[0].Description, nil
}

func (g *GoogleVision) GetTextFromPdfUrl(name string) (string, error) {

	ctx := context.Background()

	client, err := storage.NewClient(ctx)

	bucketName := g.bucket
	filename := name

	if err != nil {
		log.Fatalf("Failed to create Google Cloud Storage client: %v", err)
	}
	defer client.Close()

	// Open the PDF file from the Google Cloud Storage bucket
	bucket := client.Bucket(bucketName)
	obj := bucket.Object(filename)
	rc, err := obj.NewReader(ctx)
	if err != nil {
		log.Fatalf("Failed to create object reader: %v", err)
	}
	defer rc.Close()

	g.logger.Info("Start to read PDF file")

	// Create a Vision API client and send a request to detect text in the PDF file
	visionClient, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create Vision API client: %v", err)
	}
	feature := &visionpb.Feature{
		Type: visionpb.Feature_DOCUMENT_TEXT_DETECTION,
	}
	gcsSource := &visionpb.GcsSource{
		Uri: fmt.Sprintf("gs://%s/%s", bucketName, filename),
	}
	inputConfig := &visionpb.InputConfig{
		GcsSource: gcsSource,
		MimeType:  "application/pdf",
	}
	asyncReq := &visionpb.AsyncBatchAnnotateFilesRequest{
		Requests: []*visionpb.AsyncAnnotateFileRequest{
			{
				Features:    []*visionpb.Feature{feature},
				InputConfig: inputConfig,
				OutputConfig: &visionpb.OutputConfig{
					GcsDestination: &visionpb.GcsDestination{
						Uri: fmt.Sprintf("gs://%s/output/%s+", bucketName, filename),
					},
					BatchSize: 10,
				},
			},
		},
	}
	operation, err := visionClient.AsyncBatchAnnotateFiles(ctx, asyncReq)
	if err != nil {
		log.Fatalf("Failed to call AsyncBatchAnnotateFiles: %v", err)
	}

	// Wait for the operation to complete and retrieve the results
	operationResponse, err := operation.Wait(ctx)
	if err != nil {
		log.Fatalf("Failed to wait for operation: %v", err)
	}
	g.logger.Info("Done waiting for operation")
	g.logger.Info("Response: ", operationResponse.GetResponses())
	g.logger.Info("Response [LEN]: ", len(operationResponse.GetResponses()))

	outputConfig := operationResponse.GetResponses()[0].GetOutputConfig()

	g.logger.Info("OutputConfig: ", outputConfig)

	prefix := fmt.Sprintf("output/%s+", filename)

	fileNames := g.findFile(prefix)

	g.logger.Info("FileNames: ", fileNames)

	var texts []string

	for _, fileName := range fileNames {
		text, err := g.getTextFromGcs(fileName, bucket)
		if err != nil {
			//return "", err
		}
		texts = append(texts, text)
	}

	textt := strings.Join(texts, " ")

	return textt, nil
}

func (g *GoogleVision) getTextFromGcs(filename string, bucket *storage.BucketHandle) (string, error) {
	ctx := context.Background()

	obj := bucket.Object(filename)

	rc, err := obj.NewReader(ctx)
	if err != nil {
		log.Fatalf("Failed to create output object reader: %v", err)
	}
	defer func(rc *storage.Reader) {
		err := rc.Close()
		if err != nil {

		}
	}(rc)

	var response GcsResponse
	if err := json.NewDecoder(rc).Decode(&response); err != nil {
		log.Fatalf("Failed to parse output JSON: %v", err)
	}

	texts := make([]string, 0)
	for _, r := range response.Responses {
		g.logger.Info(r.FullTextAnnotation.Text)
		texts = append(texts, r.FullTextAnnotation.Text)
	}
	text := strings.Join(texts, " ")
	return text, nil
}

func (g *GoogleVision) findFile(prefix string) []string {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)

	if err != nil {
		g.logger.Error(err)
	}

	defer func(client *storage.Client) {
		err := client.Close()
		if err != nil {
			panic(err)
		}
	}(client)

	bucket := client.Bucket(g.bucket)

	var files []string

	query := &storage.Query{
		Prefix: prefix,
	}

	it := bucket.Objects(ctx, query)

	for {
		objAttrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			g.logger.Error(err)
		}
		files = append(files, objAttrs.Name)
	}

	return files
}
