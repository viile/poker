NAME=poker

.PHONY: build
build:
	go build -v -o $(NAME) -tags=jsoniter main.go
