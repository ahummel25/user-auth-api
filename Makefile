BUILD_PREFIX          = env GOARCH=amd64 GOOS=linux CGO_ENABLED=1 go build
COMMON_LDFLAGS        = "-linkmode external -extldflags -static -s -w"

.PHONY: build clean deploy gomodgen

build: lambda/*
	for func in lambda/*; do \
		mkdir -p build/packages/$$(basename $${func}) ; \
		$(BUILD_PREFIX) -ldflags=$(COMMON_LDFLAGS) -o build/packages/$$(basename $${func})/bootstrap lambda/$$(basename $${func})/*.go ; \
		zip -j build/packages/$$(basename $${func}).zip build/packages/$$(basename $${func})/bootstrap ; \
	done

clean:
	rm -rf ./build ./vendor Gopkg.lock

deploy: clean build
	AWS_SDK_LOAD_CONFIG=1 sls deploy --verbose

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh
