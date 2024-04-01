# Build image
FROM golang:1.21

WORKDIR /go/src/github.com/voidshard/faction
COPY go.mod ./
COPY go.sum ./
COPY cmd cmd
COPY pkg pkg
COPY internal internal
RUN go mod download

# go-sqlite3 requires cgo, but if we're running in a container we're not using sqlite3 so ... disable
# If for some reason you do want to use sqlite3 from a container (?) then you'll need to set CGO_ENABLED=1
RUN CGO_ENABLED=0 GOOS=linux go build -o /app ./cmd/faction && chown -R $USER_UID:$USER_GID /app

# App image
FROM alpine

ARG USER=app
ARG GROUPNAME=$USER
ARG UID=12345
ARG GID=23456

RUN addgroup \
    --gid "$GID" \
    "$GROUPNAME" \
&&  adduser \
    --disabled-password \
    --gecos "" \
    --home "$(pwd)" \
    --ingroup "$GROUPNAME" \
    --no-create-home \
    --uid "$UID" \
    $USER

COPY --from=0 /app /app
COPY migrations migrations
RUN chown -R $USER:$GROUPNAME /app migrations

USER $USERNAME
ENTRYPOINT ["/app"]
