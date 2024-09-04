test_to_file = go test -coverprofile=coverage.out


coverage:
	$(test_to_file) ./internal/handlers/ ./internal/usecase
	go tool cover -html=coverage.out

build: docker-compose up --build 

mockery:
	mockery --dir ./internal/repository --output ./mocks/repository --all 