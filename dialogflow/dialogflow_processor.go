package dialogflow

import (
	"context"
	"log"

	dialogflow "cloud.google.com/go/dialogflow/apiv2"
	"google.golang.org/api/option"
	dialogflowpb "google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

type DFProcessor struct {
	authJSONFilePath string
	sessionClient    *dialogflow.SessionsClient
}

func (dp *DFProcessor) NewDFProcessor() {
	ctx := context.Background()
	sessionClient, err := dialogflow.NewSessionsClient(ctx, option.WithCredentialsFile(dp.authJSONFilePath))
	if err != nil {
		log.Fatal("Error in auth with Dialogflow")
	}
	dp.sessionClient = sessionClient
}

func (dp *DFProcessor) ProcessIntent(ctx context.Context) {
	req := &dialogflowpb.DetectIntentRequest{
		// TODO: Fill request struct fields.
	}
	resp, err := dp.sessionClient.DetectIntent(ctx, req)
	if err != nil {
		// TODO: Handle error.
	}
	// TODO: Use resp.
	_ = resp
}
