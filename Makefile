setup:
	-brew install goreleaser/tap/goreleaser
	-brew install jondot/tap/goweight
	# hack to workaround go 1.11 modules, see 
	# https://github.com/vektra/mockery/issues/213
	# https://github.com/golang/go/issues/24250
	go get -u github.com/onsi/ginkgo/ginkgo
	go get github.com/vektra/mockery/.../


deps:
	go mod tidy && go mod vendor

build-linux:
	GOOS=linux GOARCH=amd64 go build

build:
	go build

check-admission:
	./lightscreen  --config examples/spectral-notary/admission.yaml  --check admission-sample.json

kube-start:
	cd deployment/kind-setup && ./setup.sh
	make kube-load-to-kind
	@echo --- DONE --- 
	kubectl proxy

kube-load-to-kind:
	GOOS=linux GOARCH=amd64 go build
	docker build . -t jondot/lightscreen:v0.6
	kind load docker-image jondot/lightscreen:v0.6

kube-deploy:
	kubectl create -f deployment/local/deployment.yaml
	kubectl create -f deployment/local/service.yaml
	kubectl create -f deployment/local/webhook-ca-bundle.yaml

watch:
	ginkgo watch ./...


release:
	goreleaser --rm-dist


mocks:
	rm -rf mocks && mockery -all -dir pkg


test:
	go test ./pkg/... -cover

test-update:
 	UPDATE_SNAPSHOTS=true make test


coverage:
	go test ./pkg/... -coverprofile=cover.out
	go tool cover -html=cover.out


bench:
	go test ./... -test.bench -test.benchmem


weight:
	goweight


.PHONY: deps setup release mocks bench coverage weight test 
