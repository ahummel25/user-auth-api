BUILD_PREFIX          = env GOOS=linux go build
COMMON_LDFLAGS        = -s -w

.PHONY: build clean deploy gomodgen

build: lambda/*
	for func in lambda/*; do \
		mkdir -p build/packages/$$(basename $${func}) ; \
		$(BUILD_PREFIX) -ldflags="$(COMMON_LDFLAGS)" -o build/packages/$$(basename $${func})/bootstrap lambda/$$(basename $${func})/*.go ; \
		zip -j build/packages/$$(basename $${func}).zip build/packages/$$(basename $${func})/bootstrap ; \
	done

clean:
	rm -rf ./build ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh
