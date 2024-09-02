test_to_file = go test -coverprofile=coverage.out


coverage:
	$(test_to_file) ./internal/handlers/ ./internal/usecase
	go tool cover -html=coverage.out