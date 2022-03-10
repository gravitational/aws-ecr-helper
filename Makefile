VERSION = $(shell cat VERSION)

IMAGE_REGISTRY := quay.io
IMAGE_NAMESPACE := gravitational
IMAGE_NAME := aws-ecr-helper

IMAGE = $(IMAGE_REGISTRY)/$(IMAGE_NAMESPACE)/$(IMAGE_NAME):$(VERSION)

.PHONY: docker-build
docker-build:
	docker build -t $(IMAGE) .

.PHONY: docker-push
docker-push:
	docker push $(IMAGE)

.PHONY: test
test: docker-build
	@IMAGE=$(IMAGE) ./test.sh