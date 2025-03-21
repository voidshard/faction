# Build the application
FROM golang:1.22

WORKDIR /go/src/github.com/voidshard/faction
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY cmd cmd
COPY pkg pkg
COPY internal internal

RUN CGO_ENABLED=0 GOOS=linux go build -o /faction ./cmd/faction/*.go

# Create a minimal image
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

COPY --from=0 /faction /faction
RUN chown -R $USER:$GROUPNAME /faction

USER $USER
ENTRYPOINT ["/faction"]
