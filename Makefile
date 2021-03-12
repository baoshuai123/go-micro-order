

.PHONY: proto
proto:
	docker run --rm -v d:/GOLANG/src/taobao/order:/d/GOLANG/src/taobao/order -w /d/GOLANG/src/taobao/order  -e ICODE=2606C833CD172F4C cap1573/cap-protoc -I ./   --go_out=./ --micro_out=./ ./proto/order/order.proto

.PHONY: build
build: 

	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o order-service *.go

.PHONY: test
test:
	go test -v ./... -cover

.PHONY: docker
docker:
	docker build . -t order-service:latest
