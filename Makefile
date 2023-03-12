APP_NAME=url-shortener-api
APP_EXECUTABLE_WINDOWS="./bin/${APP_NAME}.exe"
APP_BINARY_LINUX="./bin/${APP_NAME}"

test:
	@go test ./... ;

test.verbose:
	@go test -v ./... ;

test.coverage:
	@go test -cover ./... ;

test.coverage.out:
	@go test ./... -coverprofile="coverage.out"
	# now run 'go tool cover -html="coverage.out"'

build.windows:
	@env GOOS=windows GOARCH=amd64 go build -o ${APP_EXECUTABLE_WINDOWS} ./cmd/server

build.linux:
	@env GOOS=linux GOARCH=amd64 go build -o ${APP_BINARY_LINUX} ./cmd/server