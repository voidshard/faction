


### Local Development

Initial spin up
```sh
# Stand up infra
docker compose up postgres

# In another terminal, setup DB
./cmd/pg_updater/run.sh

# Now we can spin up everything
docker compose up
```

Reseting infra (drops everything)
```
docker compose rm
```

Switch to using a sqlitedb & local 'queue' by default with
```
export ENABLE_LOCAL_MODE=true
```
