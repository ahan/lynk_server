.DEFAULT_GOAL := local

local:
	@/Users/dpv/go/bin/air

prod:
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o app ./main.go

deploy:
	@scp app root@47.98.111.217:/root/app/lynk/
	@scp app root@116.62.7.86:/root/app/lynk/
