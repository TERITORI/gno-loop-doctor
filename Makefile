DOCKER_REGISTRY=rg.nl-ams.scw.cloud/teritori
DOCKER_IMAGE=$(DOCKER_REGISTRY)/loop-doctor:$(shell git rev-parse --short HEAD)

.PHONY: docker.build
docker.build:
	docker build . -t $(DOCKER_IMAGE)

.PHONY: docker.publish
docker.publish: docker.build
	docker push $(DOCKER_IMAGE)