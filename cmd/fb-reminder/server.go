package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/andboson/fb-reminder-go/facebook"
	"github.com/andboson/fb-reminder-go/processor"
	"github.com/andboson/fb-reminder-go/reminders"

	"github.com/golang/protobuf/jsonpb"
	log "github.com/sirupsen/logrus"
	"google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

type Service struct {
	address string
	srv     *http.Server
	rm      reminders.Reminderer
	fb      facebook.FBManager
}

func NewService(address string, rm reminders.Reminderer, fb facebook.FBManager) *Service {
	var server = &Service{
		address: address,
		fb:      fb,
		rm:      rm,
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
		log.WithError(err).Errorf("Couldn't start srv")
	}

	return err
}

func (s *Service) handleWebhook(w http.ResponseWriter, req *http.Request) {
	var err error

	wr := dialogflow.WebhookRequest{}
	if err = jsonpb.Unmarshal(req.Body, &wr); err != nil {
		log.WithError(err).Error("Couldn't Unmarshal request to jsonpb")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	resp, err := s.dispatch(wr, nil, s.fb)
	if err != nil {
		log.WithError(err).Printf("err dispatch request")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(resp)
}

func (s *Service) dispatch(wr dialogflow.WebhookRequest, dfp processor.Processor, fb facebook.FBManager) ([]byte, error) {
	var resp interface{}
	var ctx = context.Background()
	//    agent.fb.ShowMenu(agent.originalRequest.payload.data.sender.id);
	fmt.Printf("\n >>> %+v ", wr.GetQueryResult().GetIntent().GetName())
	fmt.Printf("\n >>> %+v ", wr.GetQueryResult().GetIntent().GetDisplayName())
	fmt.Printf("\n >>> %+v ", wr.GetQueryResult().GetIntent().GetParameters())

	fbClientID := extractFBClientID(wr)
	switch wr.GetQueryResult().GetIntent().GetName() {
	case "menu":
		dfp.ShowMenu(ctx, fbClientID)

	}

	br, err := json.Marshal(resp)
	if err != nil {
		log.WithError(err).Printf("err marshall response")
		br = []byte(err.Error())
	}

	return br, nil
}

func extractFBClientID(wr dialogflow.WebhookRequest) string {
	odir := wr.GetOriginalDetectIntentRequest()
	fmt.Printf("\n 00>>-_->>>>>> %+v ", odir.GetPayload().GetFields())
	fmt.Printf("\n 11>>-_->>>>>> %+v ", wr.GetQueryResult().GetParameters())
	fmt.Printf("\n >>-_____-> %+v ", wr)

	var fbID string
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
