version=latest
image=mwaaas/sqs_go_demo

cache_dir=.cache
cache_from := $(shell [ -f $(cache_dir)/index.json ] && echo "--cache-from=type=local,src=$(cache_dir)" || echo )

build:
	docker buildx build $(cache_from) --cache-to=type=local,dest=$(cache_dir) --output=type=docker,name=aws_helper_app .

push_image:
	docker tag sqs_demo_sqs_demo $(image):$(version)
	docker push $(image):$(version)

deploy: build push_image