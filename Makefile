BUILD_PREFIX          = env GOOS=linux GOARCH=arm64 go build
COMMON_LDFLAGS        = -s -w -extldflags '-static'
COMMON_TAGS           = lambda.norpc netgo
DEPLOYMENT_STAGE      = ${TF_WORKSPACE}

# Update this line to set environment-specific tags for prod and dev
ENV_TAGS              = $(if $(filter prod,$(DEPLOYMENT_STAGE)),prod,$(if $(filter dev,$(DEPLOYMENT_STAGE)),dev,))

.PHONY: build clean deploy gomodgen local dev-deps logs logs-mongo

build: 
	# Build Lambda functions
	for func in lambda/*/cmd/lambda; do \
		funcname=$$(echo $$func | cut -d'/' -f2); \
		mkdir -p build/packages/$$funcname ; \
		$(BUILD_PREFIX) -ldflags "$(COMMON_LDFLAGS)" -tags "$(COMMON_TAGS) $(ENV_TAGS)" \
			-trimpath \
			-o build/packages/$$funcname/bootstrap $$func/*.go ; \
		zip -j build/packages/$$funcname.zip build/packages/$$funcname/bootstrap ; \
	done

clean:
	rm -rf ./build ./vendor Gopkg.lock

deploy: clean build
	sls deploy --stage $(DEPLOYMENT_STAGE) --verbose

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh

# Local development targets
local:
	go run cmd/local/main.go

dev-deps:
	docker compose up -d

# Logging targets
logs-mongo:
	docker compose logs -f mongodb

logs: logs-mongo