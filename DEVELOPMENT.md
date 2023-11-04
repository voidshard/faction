


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

Switch to using a sqlitedb by default with
```
export FACTION_DEFAULT_DB_DRIVER=sqlite3
```
This works fine unless your tests need to use the async queue (there isn't current a  local only mode for this) 
