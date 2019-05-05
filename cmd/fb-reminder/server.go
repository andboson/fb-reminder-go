package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/golang/protobuf/jsonpb"
	"google.golang.org/genproto/googleapis/cloud/dialogflow/v2"
)

type Service struct {
	address string
	srv     *http.Server
	disp    Dispatcherer
}

func NewService(address string, disp Dispatcherer) *Service {
	var server = &Service{
		address: address,
		disp:    disp,
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

	resp, err := s.disp.Dispatch(wr)
	if err != nil {
		log.Printf("err dispatch request: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	br, err := json.Marshal(resp)
	if err != nil {
		log.Printf("err marshall response: %s", err)
		br = []byte(err.Error())
	}

	w.Write(br)
}
