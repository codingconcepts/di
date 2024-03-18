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

#### All CockroachDB data types

Create cluster

``` sh
cockroach demo --insecure --no-example-database
```

Create database and table

``` sh
cockroach sql \
  --url "postgres://root@localhost:26257?sslmode=disable" \
  --file examples/crdb_data_types/create.sql
```

Generate data (using [dg](http://github.com/codingconcepts/dg))

``` sh
dg -c examples/crdb_data_types/dg.yaml -o examples/crdb_data_types/csvs
```

Import data using di

``` sh
di \
  --url "postgres://root@localhost:26257?sslmode=disable" \
  --file examples/crdb_data_types/csvs/example.csv
```

##### Generated columns

Create cluster

``` sh
cockroach demo \
  --demo-locality region=aws-us-east-1:region=aws-eu-central-1:region=aws-ap-southeast-1 \
  --nodes 3 \
  --no-example-database \
  --insecure
```

Create database and table

``` sh
cockroach sql \
  --url "postgres://root@localhost:26257?sslmode=disable" \
  --file examples/generated_columns/create.sql
```

Generate data (using [dg](http://github.com/codingconcepts/dg))

``` sh
dg -c examples/generated_columns/dg.yaml -o examples/generated_columns/csvs
```

Import data using di

``` sh
di \
  --url "postgres://root@localhost:26257/example?sslmode=disable" \
  --file examples/generated_columns/csvs/example.csv
```

### Helpful statements

CockroachDB column type fetch

``` sql
SELECT ordinal_position, column_name, udt_name, is_nullable, is_generated
FROM information_schema.columns
WHERE table_name = 'example'
AND table_schema = 'public'
ORDER BY ordinal_position;
```

### Todo

* Unit tests