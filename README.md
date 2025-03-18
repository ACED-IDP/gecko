# gecko

gecko is a configuration server used for fetching inserting user specified configurations that are set during etl jobs and read from the frontend.

## setup

currently DB scripts are hardcoded for local development
Very basic GET and PUT tests should be working

```
./init_postgres.sh
go build -o bin/gecko && ./bin/gecko
go test -v ./...
```
