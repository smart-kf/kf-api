.PHONY: copy-config build run mv-tpl build-image reload

build:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/app cmd/server/main.go

build-image:build
	@docker build -t kf-api .

reload:
	@docker compose stop && docker compose rm -f && docker compose up -d
