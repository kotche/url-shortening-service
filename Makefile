gen:
	mockgen -source=internal/app/service/service.go -destination=internal/app/storage/mock/mock_repository.go