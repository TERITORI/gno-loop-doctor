DOCKER_REGISTRY=rg.nl-ams.scw.cloud/teritori
DOCKER_IMAGE=$(DOCKER_REGISTRY)/loop-doctor:$(shell git rev-parse --short HEAD)
INDEXER_DOCKER_IMAGE=$(DOCKER_REGISTRY)/loop-indexer-runner:$(shell git rev-parse --short HEAD)

.PHONY: docker.build
docker.build:
	docker build . --platform linux/amd64 -t $(DOCKER_IMAGE)

.PHONY: docker.publish
docker.publish: docker.build
	docker push $(DOCKER_IMAGE)


.PHONY: docker.indexer.build
docker.indexer.build:
	docker build . -f indexer-runner/Dockerfile --platform linux/amd64 -t $(INDEXER_DOCKER_IMAGE)

.PHONY: docker.indexer.publish
docker.indexer.publish: docker.indexer.build
	docker push $(INDEXER_DOCKER_IMAGE)