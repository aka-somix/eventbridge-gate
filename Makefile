dev-build:
	go build -o egate

release:
	goreleaser --clean
