FROM golang:1.22

ARG USERNAME=app
ARG USER_UID=1000
ARG USER_GID=$USER_UID

RUN groupadd --gid $USER_GID $USERNAME \
    && useradd --uid $USER_UID --gid $USER_GID -m $USERNAME

WORKDIR /go/src/github.com/voidshard/faction
COPY go.mod ./
COPY go.sum ./
COPY cmd cmd
COPY pkg pkg
COPY internal internal
RUN go mod download

# go-sqlite3 requires cgo, but if we're running in a container we're not using sqlite3 so ... disable
# If for some reason you do want to use sqlite3 from a container (?) then you'll need to set CGO_ENABLED=1
RUN CGO_ENABLED=0 GOOS=linux go build -o /faction ./cmd/faction && chown -R $USER_UID:$USER_GID /faction

USER $USERNAME
ENTRYPOINT ["/faction"]
