.PHONY: run_api
run_api:
	go run cmd/api/main.go

.PHONY: build_api
build_api:
	go build -o bin/dvigus_task cmd/api/main.go

.PHONY: run_bin_api
run_bin_api:
	./bin/dvigus_task

.PHONY: background_run_api
background_run_api:
	nohup go run cmd/api/main.go > api_logs.log &

.PHONY: docker
docker:
	docker stop api || true
	docker rm api || true
	docker build . -t api
	docker run --restart always \
	--add-host=host.docker.internal:host-gateway \
	--name api \
	-p 3001:3001 -d \
	api
