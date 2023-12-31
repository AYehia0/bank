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

## ORMS
while using something like GORM is really helpful, it has some drawbacks like performance on high load and learning to use the functions GORM have is a must, SQLC, on the other side, is very fast and easy to use and it has automatic code generation for queries written in SQL unfortunately it only supports few DBs (MySQL, PostgreSQL, SQLite).

### Getting started

- [Download SQLC](https://sqlc.dev/)

Running `sqlc init` to create a config file `sqlc.yaml`

## Database Transactions
A database transaction is a logical unit of work that consists of one or more database operations, such as inserts, updates, or deletions. It represents a series of actions that must be executed together as a single, indivisible unit. The concept of transactions ensures that if one part of the transaction fails, the entire transaction is rolled back, and the database is left in its original state.

Transactions provide the following key properties, often referred to as **ACID** properties:

3. **Atomicity**: Transactions are atomic, meaning they are treated as a single unit of work. Either all the operations within a transaction are completed successfully, or none of them are. If any part of the transaction fails, all changes made by the transaction are rolled back, and the database is left unchanged.

2. **Consistency**: Transactions ensure that the database moves from one consistent state to another. The database must satisfy all defined rules, constraints, and relationships during and after the transaction. If a transaction violates any of these rules, it is rolled back, and the original state is restored.

4. **Isolation**: Transactions are isolated from each other, meaning that the intermediate states of concurrent transactions are not visible to each other. Each transaction operates as if it is the only transaction being executed, which prevents interference or data corruption caused by concurrent access.

5. **Durability**: Once a transaction is committed, its changes are permanent and survive system failures, such as power outages or crashes. The committed data is stored in a way that it can be recovered and restored even in the event of a system failure.

Transactions ensure data integrity and help maintain the reliability and consistency of a database system. They are widely used in applications where maintaining data integrity and handling concurrent access is crucial, such as banking systems, e-commerce platforms, and enterprise applications.

In this project, Ensuring ACID is very important since we're dealing with money transfer, we have to implement transaction store.

## Database Deadlock

```
=== RUN   TestTransferTransaction
    store_test.go:38:
                Error Trace:    /home/none/Things/github/go-backend/db/sqlc/store_test.go:38
                Error:          Received unexpected error:
                                pq: deadlock detected
                Test:           TestTransferTransaction
--- FAIL: TestTransferTransaction (1.01s)
```

A database deadlock occurs when two or more transactions permanently hold resources (such as database locks) and wait for each other to release the resources, resulting in a situation where none of the transactions can proceed. In other words, a deadlock is a state in which multiple transactions are blocked indefinitely, unable to complete their execution.

Deadlocks typically occur in multi-user or concurrent database systems where multiple transactions are executed simultaneously. They can happen due to various factors, such as improper transaction scheduling, concurrent access to shared resources, or incorrect implementation of locking mechanisms.

Here's a simplified example of a database deadlock scenario:

- Transaction A acquires a lock on Resource X and requests a lock on Resource Y.
- Transaction B acquires a lock on Resource Y and requests a lock on Resource X.

Both transactions are now waiting for the other transaction to release the resource they need. As a result, neither transaction can proceed, leading to a deadlock.

Database management systems (DBMS) typically have mechanisms in place to detect and resolve deadlocks. When a deadlock is detected, the DBMS can choose to automatically resolve it by aborting one of the transactions involved, allowing the other transaction(s) to continue. This process is known as deadlock detection and resolution.

Preventing deadlocks usually involves proper `transaction design and concurrency control mechanisms`. Techniques like transaction scheduling, `locking protocols`, and deadlock detection algorithms are used to minimize the occurrence of deadlocks and ensure efficient database operations.

### Deadlock detection and debugging
![deadlock](/.assets/deadlock.png)
1. Check which transaction calling which query 
2. View the blocked and blocking activiy
```sql
  SELECT blocked_locks.pid     AS blocked_pid,
         blocked_activity.usename  AS blocked_user,
         blocking_locks.pid     AS blocking_pid,
         blocking_activity.usename AS blocking_user,
         blocked_activity.query    AS blocked_statement,
         blocking_activity.query   AS current_statement_in_blocking_process
   FROM  pg_catalog.pg_locks         blocked_locks
    JOIN pg_catalog.pg_stat_activity blocked_activity  ON blocked_activity.pid = blocked_locks.pid
    JOIN pg_catalog.pg_locks         blocking_locks 
        ON blocking_locks.locktype = blocked_locks.locktype
        AND blocking_locks.database IS NOT DISTINCT FROM blocked_locks.database
        AND blocking_locks.relation IS NOT DISTINCT FROM blocked_locks.relation
        AND blocking_locks.page IS NOT DISTINCT FROM blocked_locks.page
        AND blocking_locks.tuple IS NOT DISTINCT FROM blocked_locks.tuple
        AND blocking_locks.virtualxid IS NOT DISTINCT FROM blocked_locks.virtualxid
        AND blocking_locks.transactionid IS NOT DISTINCT FROM blocked_locks.transactionid
        AND blocking_locks.classid IS NOT DISTINCT FROM blocked_locks.classid
        AND blocking_locks.objid IS NOT DISTINCT FROM blocked_locks.objid
        AND blocking_locks.objsubid IS NOT DISTINCT FROM blocked_locks.objsubid
        AND blocking_locks.pid != blocked_locks.pid

    JOIN pg_catalog.pg_stat_activity blocking_activity ON blocking_activity.pid = blocking_locks.pid
   WHERE NOT blocked_locks.granted;
```
Or even better :
```sql
SELECT a.datname,
         l.relation::regclass,
         l.transactionid,
         l.mode,
         l.GRANTED,
         a.usename,
         a.query,
         a.query_start,
         age(now(), a.query_start) AS "age",
         a.pid
FROM pg_stat_activity a
JOIN pg_locks l ON l.pid = a.pid
ORDER BY a.query_start;
```
3. Implement a deadlock prevention technique.
