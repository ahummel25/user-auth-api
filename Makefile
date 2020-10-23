BUILD_PREFIX          = env GOOS=linux GOARCH=arm64 go build
COMMON_LDFLAGS        = -s -w -extldflags '-static'
COMMON_TAGS           = lambda.norpc netgo
DEPLOYMENT_STAGE      = ${TF_WORKSPACE}

# Update this line to set environment-specific tags for prod and dev
ENV_TAGS              = $(if $(filter prod,$(DEPLOYMENT_STAGE)),prod,$(if $(filter dev,$(DEPLOYMENT_STAGE)),dev,))

.PHONY: build clean deploy gomodgen

build: lambda/*
	for func in lambda/*; do \
		mkdir -p build/packages/$$(basename $${func}) ; \
		$(BUILD_PREFIX) -ldflags "$(COMMON_LDFLAGS)" -tags "$(COMMON_TAGS) $(ENV_TAGS)" \
			-trimpath \
			-o build/packages/$$(basename $${func})/bootstrap lambda/$$(basename $${func})/*.go ; \
		zip -j build/packages/$$(basename $${func}).zip build/packages/$$(basename $${func})/bootstrap ; \
	done

clean:
	rm -rf ./build ./vendor Gopkg.lock

deploy: clean build
	sls deploy --stage $(DEPLOYMENT_STAGE) --verbose

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh