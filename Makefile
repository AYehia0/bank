# postgres container configs
db_location = localhost
postgres_container_name = postgres12
port = 5432
postgres_user = root
postgres_pass = secret# TODO: add the password as env_var
postgres_image = postgres:12-alpine

# postgres db configs
db_name = simple_bank
postgres_createdb = createdb --username=root --owner=root $(db_name)
postgres_dropdb = dropdb $(db_name)
postgres_url = postgresql://$(postgres_user):$(postgres_pass)@$(db_location):$(port)/$(db_name)?sslmode=disable

# black magic
args = $(foreach a,$($(subst -,_,$1)_args),$(if $(value $a),$a="$($a)"))

# create the postgres container with the configs like username and password
postgres:
	docker run --name $(postgres_container_name) -p $(port):$(port) -e POSTGRES_USER=$(postgres_user) -e POSTGRES_PASSWORD=$(postgres_pass) -d $(postgres_image)

run:
	docker start $(postgres_container_name)

# create the database in the docker container
createdb:
	docker exec -it postgres12 $(postgres_createdb)

# drop the database 
dropdb:
	docker exec -it postgres12 $(postgres_dropdb)

migrateup:
	migrate -path db/migrations -database $(postgres_url) -verbose up

migratedown:
	migrate -path db/migrations -database $(postgres_url) -verbose down

migrateup1:
	migrate -path db/migrations -database $(postgres_url) -verbose up 1

migratedown1:
	migrate -path db/migrations -database $(postgres_url) -verbose down 1

# a simple trick to have the name as arg without passing name=$
name := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
migrate_create:
	migrate create -ext sql -dir db/migrations -seq $(name);
%::
	@true

# generate sqlc queries
sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run .

mock:
	mockgen --destination db/mock/transaction_store.go --package storedb github.com/AYehia0/go-bk-mst/db/sqlc Store

proto:
	rm -f pb/*.go
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
    proto/*.proto

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server mock migrateup1 migratedown1 run migrate_create proto
