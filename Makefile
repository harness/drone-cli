
deps:
	cd drone
	go get -u ./...

build:
	cd drone
	go install ./...

test:
	cd drone
	go test ./...

clean:
	rm -rf bin dist

dist:
	mkdir -p bin dist
	echo dist