BUILD_PREFIX          = env GOOS=linux go build
COMMON_LDFLAGS        = -s -w

.PHONY: build clean deploy gomodgen

# build: gomodgen
# 	export GO111MODULE=on
# 	env GOOS=linux go build -ldflags="-s -w" -o bin/hello hello/main.go
# 	env GOOS=linux go build -ldflags="-s -w" -o bin/world world/main.go

build: lambda/*
	for func in lambda/*; do \
		echo $$(basename $${func}) ; \
		mkdir -p build/packages/$$(basename $${func}) ; \
		$(BUILD_PREFIX) -ldflags="$(COMMON_LDFLAGS)" -o build/packages/$$(basename $${func})/bootstrap lambda/$$(basename $${func})/*.go ; \
		zip -j build/packages/$$(basename $${func}).zip build/packages/$$(basename $${func})/bootstrap ; \
	done

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh
