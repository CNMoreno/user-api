test_to_file = go test -coverprofile=coverage.out


coverage:
	$(test_to_file) ./internal/adapters ./internal/handlers/ ./internal/repository ./internal/security ./internal/usecase
	go tool cover -html=coverage.out

build: docker-compose up --build 

mockery:
	mockery --dir ./internal/repository --output ./mocks/repository --all 