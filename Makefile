test:
	@go test ./... ;

test.verbose:
	@go test -v ./... ;

test.coverage:
	@go test -cover ./... ;

test.coverage.out:
	@go test ./... -coverprofile="coverage.out"
	# now run 'go tool cover -html="coverage.out"'
