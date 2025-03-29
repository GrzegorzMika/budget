build_docker_image:
	docker build -t gregmika/budget_app:v0.3 .
	docker login -u gregmika --password-stdin
	docker push gregmika/budget_app:v0.3

compile:
	GOOS=linux GOARCH=arm64 go build -o ./build/budget_app main.go

deploy:
	scp ./build/budget_app  grzegorzmika@berry1:/home/grzegorzmika/budget_app

PHONY: compile deploy