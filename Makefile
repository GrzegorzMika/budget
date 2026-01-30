build_docker_image:
	docker build -t registry.gregdev.dev/library/budget_app:v0.4 .
	docker push registry.gregdev.dev/library/budget_app:v0.4

compile:
	GOOS=linux GOARCH=arm64 go build -o ./build/budget_app main.go

deploy:
	scp ./build/budget_app  grzegorzmika@berry1:/home/grzegorzmika/budget_app

PHONY: compile deploy