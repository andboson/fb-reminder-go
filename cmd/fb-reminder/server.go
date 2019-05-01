package main

import (
	"fmt"
	"net/http"

	"github.com/andboson/fb-reminder-go/facebook"
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

	//    agent.fb.ShowMenu(agent.originalRequest.payload.data.sender.id);

	fmt.Printf("\n >>> %+v ", wr.GetQueryResult().GetIntent().GetName())
	fmt.Printf("\n >>> %+v ", wr.GetQueryResult().GetIntent().GetDisplayName())
	fmt.Printf("\n >>> %+v ", wr.GetQueryResult().GetIntent().GetParameters())

	ir := wr.GetOriginalDetectIntentRequest()
	fmt.Printf("\n >>-_->>>>>> %+v ", ir.GetPayload().GetFields())
	fmt.Printf("\n >>-_->>>>>> %+v ", ir.GetSource())

	fmt.Printf("\n >>-_____-> %+v ", wr)
}
