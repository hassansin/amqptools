

.PHONY: setup test release rabbit

setup:
	go get -u github.com/govend/govend
	go get -u github.com/goreleaser/goreleaser
	govend
rabbit:
	docker run --rm -d --name rabbit -p 5672:5672 -p 15672:15672 rabbitmq:3.6.5-management
# make release tag=v0.0.1
release:
	git tag -a $(tag) -m $(tag) && git push origin $(tag) && goreleaser
