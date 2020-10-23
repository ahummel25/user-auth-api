BUILD_PREFIX          = env GOOS=linux GOARCH=amd64 go build
COMMON_LDFLAGS        = "-s -w"
COMMON_TAGS           = "lambda.norpc"
DEPLOYMENT_STAGE      = ${TF_WORKSPACE}

.PHONY: build clean deploy gomodgen

build: lambda/*
	for func in lambda/*; do \
		mkdir -p build/packages/$$(basename $${func}) ; \
		$(BUILD_PREFIX) -ldflags $(COMMON_LDFLAGS) -tags $(COMMON_TAGS) -o build/packages/$$(basename $${func})/bootstrap lambda/$$(basename $${func})/*.go ; \
		zip -j build/packages/$$(basename $${func}).zip build/packages/$$(basename $${func})/bootstrap ; \
	done

clean:
	rm -rf ./build ./vendor Gopkg.lock

deploy: clean build
	sls deploy --stage $(DEPLOYMENT_STAGE) --verbose

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh
