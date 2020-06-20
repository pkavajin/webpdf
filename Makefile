TAG := v0.1.0
IMAGE := kavatech/webpdf

docker-build:
	docker build -t ${IMAGE}:${TAG} .

docker-push:
	docker push ${IMAGE}:${TAG}