REGISTRY?=eu.gcr.io/joonix.se/joonix-cloud
IMAGE?=joxy
GITHASH?=$(shell git describe --dirty --tags)

all: container

container:
	GOOS=linux GOARCH=amd64 go build .
	docker build . -t $(REGISTRY)/$(IMAGE):latest
	docker tag $(REGISTRY)/$(IMAGE):latest $(REGISTRY)/$(IMAGE):$(GITHASH)

release: container
	gcloud docker -- push $(REGISTRY)/$(IMAGE):$(GITHASH)
	gcloud docker -- push $(REGISTRY)/$(IMAGE):latest

deploy: release
	kubectl set image deployment/$(IMAGE) tracker=$(REGISTRY)/$(IMAGE):$(GITHASH)

clean:
	rm $(IMAGE)

.PHONY: all container release deploy