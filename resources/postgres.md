# Learning about Postgres

- `bigserial` 64bit auto-increment int
- `decimal/numeric` up to 131072 digits before the decimal point; up to 16383 digits after the decimal point
- Adding `indexes` to speed up searching for specific column in the table.

## Database Migration

Database migration refers to the process of making structural changes to a database schema while preserving the existing data. It involves applying a series of changes or updates to the database schema to accommodate modifications, such as adding new tables, altering columns, creating indexes, or updating relationships between tables.

Database migration is commonly used in software development projects, particularly when the application's data model evolves over time or when deploying updates to a production environment.

The migration process typically involves creating migration scripts or files that define the necessary changes to be applied to the database schema. These scripts are often written in a database-specific language (such as SQL) or using migration frameworks/tools that provide a higher-level abstraction for managing database changes.

When applying a database migration, the migration tool or framework will compare the current state of the database schema to the desired state defined in the migration scripts. It will then execute the necessary SQL statements or perform automated operations to bring the database schema in line with the desired state.

Database migration allows for version control and reproducibility of changes made to the database schema, making it easier to manage and track database changes over time. It helps ensure data integrity and consistency during the evolution of an application's data model.

## Migration Tools

- [golang-migrate](https://github.com/golang-migrate/migrate)

A great tutorial from [golang-migrate](https://github.com/golang-migrate/migrate/blob/master/database/postgres/TUTORIAL.md) about postgres migration with examples.

### Installing golang-migrate with Go toolchain

#### Versioned

```bash
$ go get -u -d github.com/golang-migrate/migrate/cmd/migrate
$ cd $GOPATH/src/github.com/golang-migrate/migrate/cmd/migrate
$ git checkout $TAG  # e.g. v4.1.0
$ # Go 1.15 and below
$ go build -tags 'postgres' -ldflags="-X main.Version=$(git describe --tags)" -o $GOPATH/bin/migrate $GOPATH/src/github.com/golang-migrate/migrate/cmd/migrate
$ # Go 1.16+
$ go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@$TAG
```

#### Unversioned

```bash
$ # Go 1.15 and below
$ go get -tags 'postgres' -u github.com/golang-migrate/migrate/cmd/migrate
$ # Go 1.16+
$ go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
```

#### Notes

1. Requires a version of Go that [supports modules](https://golang.org/cmd/go/#hdr-Preliminary_module_support). e.g. Go 1.11+
1. These examples build the cli which will only work with postgres.  In order
to build the cli for use with other databases, replace the `postgres` build tag
with the appropriate database tag(s) for the databases desired.  The tags
correspond to the names of the sub-packages underneath the
[`database`](../../database) package.
1. Similarly to the database build tags, if you need to support other sources, use the appropriate build tag(s).
1. Support for build constraints will be removed in the future: https://github.com/golang-migrate/migrate/issues/60
1. For versions of Go 1.15 and lower, [make sure](https://github.com/golang-migrate/migrate/pull/257#issuecomment-705249902) you're not installing the `migrate` CLI from a module. e.g. there should not be any `go.mod` files in your current directory or any directory from your current directory to the root
