package processor

import (
	"context"
	"log"

	dialogflow "cloud.google.com/go/dialogflow/apiv2"
	"github.com/golang/protobuf/proto"
	"google.golang.org/api/option"
	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

type Processor interface {
	HandleDefault(ctx context.Context, fbClientID string) proto.Message
}

type DFProcessor struct {
	authJSONFilePath string
	sessionClient    *dialogflow.SessionsClient
}

func NewDFProcessor(authFile string) *DFProcessor {
	var dp = DFProcessor{
		authJSONFilePath: authFile,
	}
	ctx := context.Background()
	sessionClient, err := dialogflow.NewSessionsClient(ctx, option.WithCredentialsFile(dp.authJSONFilePath))
	if err != nil {
		log.Fatal("Error in auth with Dialogflow")
	}
	dp.sessionClient = sessionClient
	return &dp
}

func (dp *DFProcessor) HandleDefault(ctx context.Context, fbClientID string) proto.Message {
	resp := &dialogflowpb.WebhookResponse{
		FulfillmentText: "Sorry, i didnt understand you",
	}

	return resp
}
