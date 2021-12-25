fmt:
	go fmt ./...

test:
	go test ./...

run_worker:
	go run cmd/worker/main.go

run_cli:
	go run cmd/cli/main.go
