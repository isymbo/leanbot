all: deps build

deps:
	@godep restore

clean:
	@rm -rf Godeps/_workspace
build:
	@godep go build .
