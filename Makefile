HOST_USER := $(shell whoami)
.PHONY: build

build:
	mkdir -p build

	docker run \
		--rm \
		-e HOST_USER=$(HOST_USER) \
		-v $(PWD):/build \
		-w /build \
		golang \
		bash -c 'go mod vendor && go build -mod=vendor -o espressod cmd/main.go'
	docker run \
		--rm \
		-e HOST_USER=$(HOST_USER) \
		-e GOOS=linux \
		-e GOARCH=arm \
		-e GOARM=5 \
                -v $(PWD):/build \
                -w /build \
                golang \
                bash -c 'go mod vendor && go build -mod=vendor -o espressod_armv5 cmd/main.go'

check: go_fmt go_lint go_vet

go_fmt:
	docker run \
		--rm \
		-v $(PWD):/installer \
		-w /installer \
		golang \
		bash -c "find . -path ./vendor -prune -o -name '*.go' -exec gofmt -l {} \; | tee fmt.out && if [ -s fmt.out ] ; then exit 1; fi "


go_vet:
	docker run\
		--rm \
		-v $(PWD):/installer \
		-w /installer \
		golang \
		bash -c "go vet ./..."

go_lint:
	docker run \
		--rm \
		-v $(PWD):/installer \
		-w /installer \
		golang \
		bash -c 'go get golang.org/x/lint/golint && go list ./... | xargs -L1 golint -set_exit_status'

run:
	docker run \
		--rm \
		-v $(PWD):/build \
		-w /build \
		golang \
		bash -c 'go run .'
		

docker:
	docker build -t astoria-hc:latest . < Dockerfile

deploy:
	scp espressod_armv5 pi@astoria:espressod_armv5


clean:
	rm -f main
	rm -f espressod
	rm -f espressod_armv5
	rm -rf build

