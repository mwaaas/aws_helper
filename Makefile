version=latest
image=mwaaas/aws_helper


deploy:
	docker build -t $(image):$(version) .
	docker push $(image):$(version)