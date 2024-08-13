build:
	go build -ldflags="-s -w" indexof.go
	$(if $(shell command -v upx || which upx), upx indexof)

mac:
	GOOS=darwin go build -ldflags="-s -w" -o indexof-darwin indexof.go
	$(if $(shell command -v upx || which upx), upx indexof-darwin)

win:
	GOOS=windows go build -ldflags="-s -w" -o indexof.exe indexof.go
	$(if $(shell command -v upx || which upx), upx indexof.exe)

linux:
	GOOS=linux go build -ldflags="-s -w" -o indexof-linux indexof.go
	$(if $(shell command -v upx || which upx), upx indexof-linux)

zip:
	go build -ldflags="-s -w" indexof.go
	cat config.prd.json > config.json
	zip indexof.zip indexof config.json