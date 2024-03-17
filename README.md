# di
A simple Data Input tool that reads a CSV and inserts it into a database

### Examples

Create cluster

``` sh
cockroach demo \
  --demo-locality region=aws-us-east-1:region=aws-eu-central-1:region=aws-ap-southeast-1 \
  --nodes 3 \
  --no-example-database \
  --insecure
```

CRDB Data Types

``` sh
cockroach sql \
  --url "postgres://root@localhost:26257?sslmode=disable" \
  --file examples/crdb_data_types/create.sql

dg -c examples/crdb_data_types/dg.yaml -o examples/crdb_data_types/csvs

go run di.go \
  --url "postgres://root@localhost:26257?sslmode=disable" \
  --file examples/crdb_data_types/csvs/example.csv \
  --fmt date:2006-01-02 \
  --fmt time:15:04:05
```
2006-01-02T

Simple

``` sh
cockroach sql \
  --url "postgres://root@localhost:26257?sslmode=disable" \
  --file examples/simple/create.sql

dg -c examples/simple/dg.yaml -o examples/simple/csvs

go run di.go \
  --url "postgres://root@localhost:26257/store?sslmode=disable" \
  --file examples/simple/csvs/customer.csv

go run di.go \
  --url "postgres://root@localhost:26257/store?sslmode=disable" \
  --file examples/simple/csvs/stock.csv
```

### Debugging

Query to fetch column types in CockroachDB

``` sql
SELECT ordinal_position, column_name, udt_name, is_nullable
FROM information_schema.columns
WHERE table_name = 'customer'
AND table_schema = 'public'
ORDER BY ordinal_position;
```

### Todo

* Unit tests
* Documentation
* Binaries