package main

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
)

const configFile = "config.json"

type Config struct {
	ServerAddress           string `json:"server_address"`
	SnoozePeriod            string `json:"snooze_period"`
	PgUser                  string `json:"pg_user"`
	PgPasswd                string `json:"pg_passwd"`
	PgDb                    string `json:"pg_db"`
	PgAddress               string `json:"pg_address"`
	XKey                    string `json:"x_key"`
	FbToken                 string `json:"fb_token"`
	Type                    string `json:"type"`
	ProjectID               string `json:"project_id"`
	PrivateKeyID            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientID                string `json:"client_id"`
	AuthURI                 string `json:"auth_uri"`
	TokenURI                string `json:"token_uri"`
	AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
	ClientX509CertURL       string `json:"client_x509_cert_url"`
}

func config() (*Config, error) {
	var cfg *Config

	b, err := ioutil.ReadFile(configFile)
	if err != nil {
		log.Errorf("unable to load " + configFile)
		return nil, err
	}

	err = json.Unmarshal(b, cfg)

	return cfg, err
}
