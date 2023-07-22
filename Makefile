
run:
	go run cmd/main.go

swag-init:
	swag init -g api/api.go -o api/docs

migration-up:
	migrate -path ./migration/postgres -database 'postgres://asadbek:3066586@localhost:5432/car_service?sslmode=disable' up

migration-down:
	migrate -path ./migration/postgres -database 'postgres://asadbek:3066586@localhost:5432/car_service?sslmode=disable' down


