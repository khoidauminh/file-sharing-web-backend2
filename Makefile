import-db:
	docker exec -i postgres-db psql -U postgres -d "file-sharing" < ./internal/infrastructure/database/init.sql
export-db:
	docker exec -i postgres-db pg_dump -U postgres -d "file-sharing" > ./internal/infrastructure/database/backup.sql
server:
	go run ./cmd/server/main.go
