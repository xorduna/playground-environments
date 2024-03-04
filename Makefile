VERSION := 1.0
CONTAINER_REGISTRY := registry.digitalocean.com/maulabs
API_IMAGE_NAME := attendance-api
ARCH := "amd64"

build-api:
	go mod tidy
	docker build \
		-f deployment/docker/api.dockerfile \
		-t $(API_IMAGE_NAME):$(VERSION) \
		--build-arg BUILD_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		--build-arg GOARCH=$(ARCH) \
		.

push-api:
	docker tag $(API_IMAGE_NAME):$(VERSION) $(CONTAINER_REGISTRY)/$(API_IMAGE_NAME):$(VERSION)
	docker push $(CONTAINER_REGISTRY)/$(API_IMAGE_NAME):$(VERSION)
