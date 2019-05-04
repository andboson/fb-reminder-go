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
	//ShowMenu(ctx context.Context, fbClientID string)
	HandleDefault(ctx context.Context, fbClientID string) proto.Message
}

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

func (dp *DFProcessor) HandleDefault(ctx context.Context, fbClientID string) proto.Message {
	resp := &dialogflowpb.Intent_Message_SimpleResponse{
		DisplayText: "sorry, i didnt understand you",
	}

	return resp
}