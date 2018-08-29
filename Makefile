# serverless create -t aws-go -n lambda-qr
all:
	@cat Makefile | grep : | grep -v PHONY | grep -v @ | sed 's/:/ /' | awk '{print $$1}' | sort

#-------------------------------------------------------------------------------

.PHONY: build
build:
	go build -ldflags="-s -w" -o bin/qr main.go

.PHONY: package
package:
	env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o bin/qr main.go

.PHONY: lint
lint:
	gometalinter.v2 ./main.go

.PHONY: deploy
deploy: package
	aws-vault exec personal --no-session -- sls deploy --verbose
