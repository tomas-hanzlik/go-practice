# Just for my personal use during dev.
GOTEST=go test
GORUN=go run

all: test app
test: 
		$(GOTEST) ./src/...
test-v: 
		$(GOTEST) ./src/... -v
play:
		$(GORUN) playground/play.go
play-crypto:
		$(GORUN) playground/cryptomood.go
app:
		./cmd/app/run.sh
