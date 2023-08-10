.DEFAULT_GOAL := all

IP = "178.158.194.3"

clean:
	@rm webserver 2> /dev/null || true

lint-front:
	@npm install && npm run lint

build-front: clean
	npm run build
	rm -rf ./internal/webserver/handlers/home/css && cp -rf ./frontend/styles ./internal/webserver/handlers/home/css
	gzip -c ./frontend/index.html > ./internal/webserver/handlers/home/index.html.gz
	gzip -f ./internal/webserver/handlers/home/index.js
	gzip -f ./internal/webserver/handlers/home/index.js.map

lint-back:
	@gofmt -s -w .
	@golangci-lint run

test-back:
	@go test -v -tags=exec ./internal/...

build-back:
	@GOOS=linux GOARCH=mipsle go build -tags=sysboard -ldflags "-s -w" -o webserver cmd/main.go

assemble: build-front build-back

run: lint-front build-front lint-back test-back
	@go run cmd/main.go

deploy: assemble
	@echo "Deploying to remote..."
	@scp webserver root@${IP}:/root/webserver.new
	@echo "Restarting webserver..."
	@ssh root@${IP} '/etc/init.d/webserver stop; chmod +x webserver.new; mv webserver.new webserver; /etc/init.d/webserver start'
	@echo "Done"

deploy-config:
	@scp config-prod.json root@${IP}:/root/config.json
	@ssh root@${IP} '/etc/init.d/webserver restart'

all: lint-front build-front lint-back test-back build-back deploy clean
