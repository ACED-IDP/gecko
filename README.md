# gecko

gecko is a configuration server used for fetching inserting user specified configurations that are set during etl jobs and read from the frontend.

## helm cluster setup

Helm setup should be much simpler since env vars are defined.

```
./init_cluster_pg.sh
go build -o bin/gecko && ./bin/gecko
go test -v ./...
```

## local setup

Make sure the below command matches whatever was specified in the init db script. For local, going to need to disable sslmode. Ex:

```
./init_postgres.sh
go build -o bin/gecko
./bin/gecko -db "postgresql://postgres:your_strong_password@localhost:5432/testdb?sslmode=disable" -port 8080
go test -v ./...
```
