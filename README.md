# godash
Redash Go implementation

# Usage
## worker
```console
# first set up postgresql
make setup_pg

# worker start
make run_worker
```

## client
```console
# send task(datasource connection settings) to worker on different terminal
make run_cli_s

# send task(query) to worker
make run_cli_q
```