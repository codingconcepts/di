# di
A simple Data Input tool that reads a CSV and inserts it into a database

### Installation

Find the release that matches your architecture on the [releases](https://github.com/codingconcepts/di/releases) page.

Download the tar, extract the executable, and move it into your PATH:

```
$ tar -xvf di_[VERSION]-rc1_macOS.tar.gz
```

### Usage

``` sh
di --help
Usage of di:
  -batch int
        import batch size (default 10000)
  -file string
        absolute or relative path to the CSV file to import
  -fmt value
        format helper (in the form of data_type:format) (default date=2006-01-02, time=15:04:05, timestamp=2006-01-02T15:04:05, timestamptz=2006-01-02T15:04:05)
  -schema string
        schema in which your table resides (default "public")
  -table string
        name of the table to import into
  -url string
        database connection string
  -version
        display version information
```

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
  --file examples/crdb_data_types/csvs/example.csv
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