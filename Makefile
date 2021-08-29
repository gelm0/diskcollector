APP=diskstat
PACKAGE_FOLDER=./df

all: test vet build

build:
	go build -o ${APP} main.go

build-static:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ${APP} main.go

clean:
	go clean
	rm -f ${APP}
	rm -f docker/${APP}

test:
	go test ${PACKAGE_FOLDER} -v

vet:
	gofmt -w ${PACKAGE_FOLDER}
	gofmt -w  main.go
	go vet -json

docker: test build-static
	cp ${APP} docker/${APP}
	docker build docker/ -t docker.hub.com/gelm0/diskcollector
