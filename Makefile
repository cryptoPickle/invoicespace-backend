generate-schemas:
	cd schema ; go run github.com/99designs/gqlgen -v

run-state:
	docker kill users-postgres | true
	docker rm users-postgres -f | true
	docker run --name users-postgres -e POSTGRES_DB=users  -d -p 5432:5432 postgres | true
	docker kill users-redis | true
	docker rm users-redis -f | true
	docker run --name users-redis -p 6379:6379 -d redis

run-dev: run-state
	bash ./scripts/run-dev.sh