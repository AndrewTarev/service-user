create-migrations:
	migrate create -ext sql -dir ./internal/app/repository/migrations -seq init

migrateup:
	migrate -path ./internal/app/repository/migrations -database 'postgres://postgres:postgres@localhost:5433/service-user?sslmode=disable' up

migratedown:
	migrate -path ./internal/app/repository/migrations -database 'postgres://postgres:postgres@localhost:5433/service-user?sslmode=disable' down

test-mock:
	mockgen -source=internal/app/service/service.go -destination=internal/app/service/mocks/mock_service.go -package=mocks

gen-docs:
	swag init -g ./cmd/main.go -o ./docs