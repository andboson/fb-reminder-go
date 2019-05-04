package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/andboson/fb-reminder-go/facebook"
	"github.com/andboson/fb-reminder-go/processor"
	"github.com/andboson/fb-reminder-go/reminders"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

type Service struct {
	address string
	srv     *http.Server
	rm      reminders.Reminderer
	fb      facebook.FBManager
	dfp     processor.Processor
}

func NewService(address string, rm reminders.Reminderer, fb facebook.FBManager, dfp processor.Processor) *Service {
	var server = &Service{
		address: address,
		fb:      fb,
		rm:      rm,
		dfp:     dfp,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", server.handleWebhook)

	server.srv = &http.Server{
		Addr:    address,
		Handler: mux,
	}

	return server
}

func (s *Service) Serve() error {
	var err error

	fmt.Println("Started listening...")
	if err = s.srv.ListenAndServe(); err != nil {
		log.Printf("Couldn't start srv: %s", err)
	}

	return err
}

func (s *Service) Stop() error {
	return s.srv.Close()
}

func (s *Service) handleWebhook(w http.ResponseWriter, req *http.Request) {
	var err error

	wr := dialogflow.WebhookRequest{}
	if err = jsonpb.Unmarshal(req.Body, &wr); err != nil {
		log.Printf("Couldn't Unmarshal request to jsonpb: %s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := s.dispatch(wr)
	if err != nil {
		log.Printf("err dispatch request: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(resp)
}

func (s *Service) dispatch(wr dialogflow.WebhookRequest) ([]byte, error) {
	var resp proto.Message
	var ctx = context.Background()
	var err error

	fmt.Printf("\n >>> %+v ", wr.GetQueryResult().GetIntent().GetDisplayName())
	fmt.Printf("\n >>> %+v ", wr.GetQueryResult().GetIntent().GetParameters())

	fbClientID := extractFBClientID(wr)
	switch wr.GetQueryResult().GetIntent().GetDisplayName() {
	case "menu":
		resp = s.dfp.HandleDefault(ctx, fbClientID)
//		err = s.fb.ShowMenu(ctx, fbClientID)

	default:
		resp = s.dfp.HandleDefault(ctx, fbClientID)
	}

	if err != nil {
		log.Printf("err dispatch intent: %s", err)
	}

	//m, err := new(jsonpb.Marshaler).MarshalToString(resp)
	//if err != nil {
	//	log.Printf("err marchal proto: %s", err)
	//}
	//br := []byte(m)
	br, err := json.Marshal(resp)

	if err != nil {
		log.Printf("err marshall response: %s", err)
		br = []byte(err.Error())
	}

	return br, nil
}

func extractFBClientID(wr dialogflow.WebhookRequest) string {
	var fbID string
	var odir = wr.GetOriginalDetectIntentRequest()
	if data, ok := odir.GetPayload().GetFields()["data"]; ok {
		if sender, ok := data.GetStructValue().GetFields()["sender"]; ok {
			senderStruct := sender.GetStructValue()
			if id, ok := senderStruct.GetFields()["id"]; ok {
				fbID = id.GetStringValue()
			}
		}
	}

	return fbID

}
