package google_vision

import (
	"cloud.google.com/go/storage"
	vision "cloud.google.com/go/vision/apiv1"
	"cloud.google.com/go/vision/v2/apiv1/visionpb"
	"context"
	"github.com/721945/dlaw-backend/libs"
	"io/ioutil"
	"log"
)

type GoogleVision struct {
	logger *libs.Logger
}

func NewGoogleVision(logger *libs.Logger) GoogleVision {
	return GoogleVision{logger: logger}
}

//func (g *GoogleVision) GetTextFromImageUrl(url string) (string, error) {
//	ctx := context.Background()
//	client, err := vision.NewImageAnnotatorClient(ctx)
//
//	if err != nil {
//		return "", err
//	}
//
//	defer func(client *vision.ImageAnnotatorClient) {
//		err := client.Close()
//		if err != nil {
//			panic(err)
//		}
//	}(client)
//
//	image := vision.NewImageFromURI(url)
//
//	annotations, err := client.DetectTexts(ctx, image, nil, 10)
//
//	if err != nil {
//		return "", err
//	}
//
//	g.logger.Info(annotations)
//
//	if len(annotations) == 0 {
//		return "", nil
//	}
//
//	return annotations[0].Description, nil
//}

func (g *GoogleVision) GetTextFromPdfUrl(obj storage.ObjectHandle) (string, error) {
	ctx := context.Background()
	//client, err := storage.NewClient(ctx)
	//if err != nil {
	//	log.Fatalf("Failed to create client: %v", err)
	//}

	// Download file from GCS bucket

	reader, err := obj.NewReader(ctx)
	if err != nil {
		log.Fatalf("Failed to create reader: %v", err)
	}

	defer reader.Close()

	// Extract text from PDF using Google Vision API
	imageBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}
	image := visionpb.Image{
		Content: imageBytes,
	}
	client, err := vision.NewImageAnnotatorClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create Vision client: %v", err)
	}
	defer client.Close()
	response, err := client.DetectDocumentText(ctx, &image, nil)
	if err != nil {
		log.Fatalf("Failed to detect text: %v", err)
	}

	text := response.Text
	g.logger.Info(text)
	return text, nil
}

//func (g *GoogleVision) GetTextFromPdfUrl(url string, target string) (string, error) {
//ctx := context.Background()
//client, err := vision.NewImageAnnotatorClient(ctx)
//
//if err != nil {
//	return "", err
//}
//
//defer func(client *vision.ImageAnnotatorClient) {
//	err := client.Close()
//	if err != nil {
//		panic(err)
//	}
//}(client)
//
////request := &visionpb.AsyncBatchAnnotateFilesRequest{
////	Requests: []*visionpb.AsyncAnnotateFileRequest{
////		{
////			Features: []*visionpb.Feature{
////				{
////					Type: visionpb.Feature_DOCUMENT_TEXT_DETECTION,
////				},
////			},
////			InputConfig: &visionpb.InputConfig{
////				GcsSource: &visionpb.GcsSource{Uri: url},
////				// Supported MimeTypes are: "application/pdf" and "image/tiff".
////				MimeType: "application/pdf",
////			},
////			OutputConfig: &visionpb.OutputConfig{
////				GcsDestination: &visionpb.GcsDestination{Uri: target},
////				// How many pages should be grouped into each json output file.
////				BatchSize: 20,
////			},
////		},
////	},
////}
////g.logger.Info("------------------")
////operation, err := client.AsyncBatchAnnotateFiles(ctx, request)
//image := vision.NewImageFromURI(url)
//response, err := client.DetectDocumentText(ctx, image, nil)
//
//g.logger.Info("------------------")
//if err != nil {
//	g.logger.Error(err)
//	return "", err
//}
//
//g.logger.Info("Waiting for the operation to finish.")
//
//g.logger.Info(response.Text)

//resp, err := operation.Wait(ctx)
//if err != nil {
//	return "", err
//}
//
//g.logger.Info(resp.String())

//image := vision.NewImageFromURI(url)
//
//annotations, err := client.DetectTexts(ctx, image, nil, 10)
//
//if err != nil {
//	return "", err
//}
//
//g.logger.Info(annotations)
//
//if len(annotations) == 0 {
//	return "", nil
//}

//	return "", nil
//}
