pull:
	git pull

build:
	mkdir -p build
	go build -o build/astoria-hc *.go

run:
	go run *.go

install:
	chmod +x build/astoria-hc
	-curl -s -XGET localhost:2112/exit
	sudo cp build/astoria-hc /usr/local/bin/astoria-hc

clean:
	rm -f main
	rm -f astoria-hc
	rm -rf build

stage: pull run

deploy: pull build install

local: build install
