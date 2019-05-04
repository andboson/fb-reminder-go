export GO111MODULE=on

all: mod test build

mod:
	go mod tidy

test: generate
	go test ./... -race -covermode atomic -coverprofile coverage.profile && go tool cover -func coverage.profile


build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o reminder ./cmd/fb-reminder

buildall:
	CGO_ENABLED=0 GOOS=linux  go build ${LDFLAGS} -o reminder ./cmd/fb-reminder
