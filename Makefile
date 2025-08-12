APP_EXECUTABLE="bin/rumi-go"
APP_BOOTSTRAP="./cmd/server"

lint:
	golangci-lint run --timeout=10m

compile:
	mkdir -p bin/
	go build -o $(APP_EXECUTABLE) $(APP_BOOTSTRAP)


copy-config:
	cp ./config.yaml.tmpl ./config.yaml

run:
	env LOG_LEVEL=debug ./bin/rumi-go start

migrate-create:
	go run cmd/server/main.go migrate:create --filename $(FILENAME)

migrate-mysql-up:
	go run cmd/server/main.go migrate:mysql:up

migrate-mysql-down:
	go run cmd/server/main.go migrate:mysql:down

migrate-mysql-status:
	go run cmd/server/main.go migrate:mysql:status

config-show:
	go run cmd/server/main.go config:show
