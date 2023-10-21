


### Local Development

Initial spin up
```sh
# Stand up infra
docker compose up

# Setup DB
./cmd/pg_updater/run.sh
```

Reseting infra (drops everything)
```
docker compose rm
```


