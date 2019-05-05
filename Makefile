export GO111MODULE=on

all: mod test build

mod:
	go mod tidy

generate:
	go generate -x ./...

test: mod generate
	go test ./... -race -covermode atomic -coverprofile coverage.profile && go tool cover -func coverage.profile

build: mod generate
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o reminder ./cmd/fb-reminder

buildall: mod generate
	CGO_ENABLED=0 GOOS=linux  go build ${LDFLAGS} -o reminder ./cmd/fb-reminder
