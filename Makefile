run-server:
	go run cmd/server/main.go

run-client:
	go run cmd/cli-client/main.go "$(SECRET)"
