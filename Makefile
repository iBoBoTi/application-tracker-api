include .env
export

docker-up:
	docker compose up --build


run-server:
	go run main.go