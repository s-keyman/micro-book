.PHONY: docker

docker:
	@rm webook||true
	@go mod tidy
	@GOOS=linux GOARCH=amd64 go build -o webook .
	@docker rmi -f a1095001/webook:v0.0.1 || true
	@docker build -t a1095001/webook:v0.0.1 .