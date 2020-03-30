HOST_USER := $(shell whoami)

build:
	mkdir -p build

	docker run \
		--rm \
		-e HOST_USER=$(HOST_USER) \
		-v $(PWD):/build \
		-w /build \
		golang \
		bash -c 'go mod vendor && go build -mod=vendor -o build/astoria-hc .'


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

clean:
	rm -f main
	rm -f astoria-hc
	rm -rf build

