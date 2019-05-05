# fb-reminder-go

Facebook reminder bot

Powered by DialogFlow AI

### install

Import `reminder.zip` in dialog flow and setup fulfillment in https://console.dialogflow.com/api-client

Setup `docker` if you want to run tests 

#### setup DB

create database and creds

#### setup config

minimal `config.json` content:

```  
{  
    "snooze_period": "5m",
    "pg_user": "********",
    "pg_passwd": "******",
    "pg_db": "***888***",
    "pg_address": "host:port",
    "x_key": "<some_secret_header>",
    "fb_token":"EAAgcpoks048BAB.................",
    "project_id": "rem.....",
    "private_key": "-----BEGIN PRIVATE KEY-----\nMI..................",
    "client_email": "email............",
}
```

#### test, build

`make all` - for mods, tests, build  

or 

just `make build` 

### run

Run `./reminder` app

App will start on `3000` port by default

All requests must be signed with `X-Key` header (value of `x_key` config field)