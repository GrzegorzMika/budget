build_docker_image:
	docker build -t ghcr.io/grzegorzmika/budget_app:v0.1 .
	echo $CR_PAT | docker login ghcr.io -u grzegorzmika --password-stdin
	docker push ghcr.io/grzegorzmika/budget_app:v0.1

compile:
	GOOS=linux GOARCH=arm64 go build -o ./build/budget_app main.go

deploy:
	scp ./build/budget_app  grzegorzmika@berry1:/home/grzegorzmika/budget_app

PHONY: compile deploy