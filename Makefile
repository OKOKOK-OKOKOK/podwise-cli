default: build

build:
	goreleaser release --snapshot --clean && mv podwise-skills.tar.gz ./dist/