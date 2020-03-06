GOTEST=go test
GORUN=go run

all: test play
test: 
		$(GOTEST) ./src/...
test-v: 
		$(GOTEST) ./src/... -v
play:
		$(GORUN) playground/play.go
app:
		$(GORUN) cmd/app/main.go
