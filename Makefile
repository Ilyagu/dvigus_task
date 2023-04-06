.PHONY: run_api
run_api:
	go run cmd/api/main.go

.PHONY: docker_build
docker_build:
	docker stop api || true
	docker rm api || true
	docker build . -t api

.PHONY: docker_run
docker_run:
	docker run --restart always \
	--add-host=host.docker.internal:host-gateway \
	--name api \
	-p 3001:3001 \
	-d api
