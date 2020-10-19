.PHONY: run database justDB service

run: justDB service

service:
	docker-compose up -d service

database: justDB
	docker cp database/schema.sql golangtechtask_database_1:/schema.sql
	docker-compose exec -e PGPASSWORD=tester database psql -U postgres -d postgres -f /schema.sql

justDB:
	docker-compose up -d database
