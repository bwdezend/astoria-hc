pull:
	git pull

build:
	/usr/local/go/bin/go build *.go

run:
	/usr/local/go/bin/go run *.go

install:
	mv gpio astoria-hc
	chmod +x astoria-hc
	-curl -s -XGET localhost:2112/exit
	sudo cp astoria-hc /usr/local/bin/astoria-hc

clean:
	rm main
	rm astoria-hc

stage: pull run

deploy: pull build install

local: build install
