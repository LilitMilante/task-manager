test:
	go test -v -race ./...

up:
	docker run --name tm-db -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=your-password -e POSTGRES_DB=task-manager -d postgres:15.3-alpine
