

.PHONY: setup test release

setup:
	go get -u github.com/govend/govend
	go get -u github.com/goreleaser/goreleaser
	govend

# make release tag=v0.0.1
release:
	git tag -a $(tag) -m $(tag) && git push origin $(tag) && goreleaser
