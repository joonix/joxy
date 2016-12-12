REGISTRY?=eu.gcr.io/joonix-cloud
IMAGE?=joxy
GITHASH?=$(shell git describe --dirty --tags --always)

all: container

container:
	GOOS=linux GOARCH=amd64 go build .
	docker build . -t $(REGISTRY)/$(IMAGE):latest
	docker tag $(REGISTRY)/$(IMAGE):latest $(REGISTRY)/$(IMAGE):$(GITHASH)

release: container
	gcloud docker -- push $(REGISTRY)/$(IMAGE):$(GITHASH)
	gcloud docker -- push $(REGISTRY)/$(IMAGE):latest

deploy: release
	kubectl set image deployment/$(IMAGE) joxy=$(REGISTRY)/$(IMAGE):$(GITHASH)

clean:
	rm $(IMAGE)

.PHONY: all container release deploy