


### Local Development

Initial spin up
```sh
# Stand up infra
docker compose up postgres

# Create Faction DB & run migrations
go run cmd/faction/*.go migrate setup
go run cmd/faction/*.go migrate up

# Create Igor DB & run migrations (https://github.com/voidshard/igor)
docker run --rm --network=host uristmcdwarf/igor:0.0.5 migrate setup
docker run --rm --network=host uristmcdwarf/igor:0.0.5 migrate up

# Now we can spin up everything
docker compose up
```

Reseting infra (drops everything)
```
docker compose rm
```

