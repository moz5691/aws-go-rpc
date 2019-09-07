GOCMD=go 
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GORUN=$(GOCMD) run
BINARY_NAME=main
BINARY_UNIX=main
CLIENT_BINARY_NAME=client
ZIP_NAME=$(BINARY_UNIX).zip
ROLE_NAME=twirp-lambda-role
ROLE_ARN=arn:aws:iam::031173249615:role/twirp-lambda-role
FUNC_NAME=movies

gen:
	retool do protoc --proto_path=$(GOPATH)/src:. --twirp_out=. --go_out=. rpc/movies/service.proto

dev-build:
	$(GOBUILD) -o build/$(BINARY_NAME) -v cmd/moviesserver/main.go

dev-run:
	$(GOBUILD) -o build/$(BINARY_NAME) -v cmd/moviesserver/main.go
	build/$(BINARY_NAME)

dev-crun:
	$(GOBUILD) -o build/$(CLIENT_BINARY_NAME) -v cmd/moviesclient/main.go
	build/$(CLIENT_BINARY_NAME)

test:
	$(GOTEST) -v ./cmd/movieservices

prod-build:
	env GOOS=linux $(GOBUILD) -o build/$(BINARY_NAME) -v cmd/moviesserver/main.go
	chmod +x build/$(BINARY_UNIX)
	zip -j build/$(ZIP_NAME) build/$(BINARY_UNIX)

clean:
	$(GOCLEAN)
	rm -f build/$(BINARY_NAME)
	rm -f build/$(BINARY_UNIX)
	rm -f build/$(ZIP_NAME)


#AWS configuration, AWS CLI only
role-config:
	aws iam create-role --role-name $(ROLE_NAME) \
	--assume-role-policy-document file://tmp/trust-policy.json

create-func:
	aws lambda create-function --function-name $(FUNC_NAME) --runtime go1.x \
	--role $(ROLE_ARN) --handler main --zip-file fileb://build/$(ZIP_NAME)

update-func:
	aws lambda update-function-code --function-name $(FUNC_NAME) \
	--zip-file fileb://build/$(ZIP_NAME)

invoke-func:
	aws lambda invoke --function-name $(FUNC_NAME) tmp/output.json

# test-api:
# 	curl -H 'Content-Type:application/json' -d '{"term":"test here"}'  <apigateway_url>

# local dynamodb
run/dynamo-up:
	docker-compose up -d dynamo

run/dynamo-init: #run/dynamo-up
	./dynamodb/create-schema-locally.sh 
	./dynamodb/load-data-locally.sh MovieSingle

run/dynamo-down:
	cd docker && docker-compose stop dynamo

run/docker-down:
	cd docker && docker-compose down

# Do this in lieu of run/docker-down @see https://docs.docker.com/compose/reference/down/
run/docker-remove-volume:
	cd docker && docker-compose down --volumes