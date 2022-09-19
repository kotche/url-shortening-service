gen:
	mockgen -source=internal/app/service/service.go -destination=internal/app/storage/mock/mock_repository.go

proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative internal/app/transport/grpc/proto/shortener.proto